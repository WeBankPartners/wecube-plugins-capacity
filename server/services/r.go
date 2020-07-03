package services

import (
	"github.com/WeBankPartners/wecube-plugins-capacity/server/models"
	"os/exec"
	"fmt"
	"log"
	"io/ioutil"
	"strings"
	"github.com/shopspring/decimal"
	"strconv"
	"sort"
	"math"
)

func RAnalyzeData(param models.RRequestParam) (err error,result models.RunScriptResult) {
	err = checkRParam(param)
	if err != nil {
		return err,result
	}
	var x [][]float64
	var y []float64
	if len(param.Monitor.Config) > 0 {
		err,yXData := AutoJustifyData(param.Monitor)
		if err != nil {
			return err,result
		}
		for _,v := range yXData.Data {
			for vi,vv := range v {
				if vi == 1 {
					y = append(y, vv)
				}
				if vi >= 2 {
					if len(x) < (vi-1) {
						x = append(x, []float64{vv})
					}else{
						x[vi-2] = append(x[vi-2], vv)
					}
				}
			}
		}
	}else{
		x = param.XData
		y = param.YData
	}
	err,result = runRscript(x,y,param.Guid)
	if err != nil {
		return err,result
	}
	param.FuncX = result.FuncX
	param.FuncB = result.FuncB
	err,chart := RChartData(param, x, y)
	result.Chart = chart
	return err,result
}

func RChartData(param models.RRequestParam,x [][]float64,y []float64) (err error,result models.EChartOption) {
	var xAxis models.AxisModel
	var yAxis,newAxis models.DataSerialModel
	var newYData []float64
	result.IsDataSeries = true
	result.Legend = []string{"real", "fun(y)"}
	for i,_ := range y {
		xAxis.Data = append(xAxis.Data, float64(i+1))
	}
	yAxis.Data = y
	yAxis.Name = "real"
	yAxis.Type = "line"
	for i,v := range param.FuncX {
		param.FuncX[i].Data = x[v.Index-1]
	}
	for _,v := range param.FuncX {
		if len(newYData) == 0 {
			for _,vv := range v.Data {
				newYData = append(newYData, v.Estimate*vv)
			}
		}else{
			for i,vv := range v.Data {
				newYData[i] = newYData[i] + v.Estimate*vv
			}
		}
	}
	if param.FuncB > 0 {
		if len(newYData) == 0 {
			for i:=0;i<len(y);i++ {
				newYData = append(newYData, param.FuncB)
			}
		}else {
			for i,_ := range newYData {
				newYData[i] = newYData[i] + param.FuncB
			}
		}
	}
	newAxis.Data = newYData
	newAxis.Name = "fun(y)"
	newAxis.Type = "line"
	result.DataSeries = []*models.DataSerialModel{&yAxis, &newAxis}
	result.Xaxis = xAxis
	return err,result
}

func runRscript(x [][]float64,y []float64,guid string) (err error,result models.RunScriptResult)  {
	var b []byte
	// build workspace
	if guid == "" {
		result.Guid = models.GetWorkspaceName()
	}else{
		result.Guid = guid
	}
	result.Workspace = fmt.Sprintf("%s/%s", models.WorkspaceDir, result.Guid)
	b,err = exec.Command("/bin/sh", "-c", fmt.Sprintf("mkdir -p %s && /bin/cp -f conf/template.r %s/template.r", result.Workspace, result.Workspace)).Output()
	if err != nil {
		log.Printf("build tmp workspace:%s fail,output:%s error:%v \n", result.Workspace, string(b), err)
		return err,result
	}

	// replace data
	b,err = ioutil.ReadFile(result.Workspace+"/template.r")
	if err != nil {
		log.Printf("replace %s/template.r data,read file fail,error:%v \n", result.Workspace, err)
		return err,result
	}
	var xDataString,xExpr string
	for i,v := range x {
		if i == len(x)-1 {
			xDataString = xDataString + fmt.Sprintf("x%d<-c(%s)", i+1, turnFloatListToString(v))
			xExpr = xExpr + fmt.Sprintf("x%d", i+1)
		}else {
			xDataString = xDataString + fmt.Sprintf("x%d<-c(%s)\n", i+1, turnFloatListToString(v))
			xExpr = xExpr + fmt.Sprintf("x%d+", i+1)
		}
	}
	tData := strings.Replace(string(b), "{x_data}", xDataString, -1)
	tData = strings.Replace(tData, "{y_data}", turnFloatListToString(y), -1)
	tData = strings.Replace(tData, "{x_expr}", xExpr, -1)
	tData = strings.Replace(tData, "{workspace}", result.Workspace, -1)
	err = ioutil.WriteFile(result.Workspace+"/template.r", []byte(tData), 0666)
	if err != nil {
		log.Printf("replace %s/template.r data,write file fail,error:%v \n", result.Workspace, err)
		return err,result
	}

	// run script
	b,err = exec.Command("/bin/sh", "-c", fmt.Sprintf("Rscript %s/template.r", result.Workspace)).Output()
	if err != nil {
		log.Printf("run Rscript %s/template.r fail,output:%s error:%v \n", result.Workspace, string(b), err)
		return err,result
	}
	output := dealWithScriptOutput(string(b))
	result.FuncX = output.FuncX
	result.FuncB = output.FuncB
	result.FuncExpr = output.FuncExpr
	result.Output = output.Output
	result.Level = output.Level
	result.Images = []string{result.Workspace+"/rp001.png",result.Workspace+"/rp002.png",result.Workspace+"/rp003.png",result.Workspace+"/rp004.png"}
	return err,result
}

func turnFloatListToString(data []float64) string {
	if len(data) == 0 {
		return ""
	}
	var s []string
	for _,v := range data {
		s = append(s, fmt.Sprintf("%.3f", v))
	}
	return strings.Join(s, ",")
}

func dealWithScriptOutput(output string) models.RunScriptResult {
	var result models.RunScriptResult
	var funcBLevel int
	var sortFuncList models.FuncXSortList
	expr := "y="
	xCount := 0
	for _,v := range strings.Split(output, "\n") {
		if strings.HasPrefix(v, "(Intercept)") {
			result.FuncB,_,funcBLevel = getEstimate(v)
			continue
		}
		if strings.HasPrefix(v, "x") {
			xCount = xCount + 1
			tmpEstimate,tmpP,tmpL := getEstimate(v)
			if tmpL > 0 {
				sortFuncList = append(sortFuncList, &models.FuncXObj{PValue:tmpP, Estimate:tmpEstimate, FuncName:fmt.Sprintf("x%d", xCount), Level:tmpL, Index:xCount})
			}
		}
	}
	sort.Sort(sortFuncList)
	result.Output = strings.Replace(output, "\n", "<br/>", -1)
	for i,v := range sortFuncList {
		if i == 0 {
			result.Level = v.Level
		}
		expr += fmt.Sprintf("%.4f%s+", v.Estimate, v.FuncName)
		result.FuncX = append(result.FuncX, v)
	}
	if funcBLevel > 0 {
		expr += fmt.Sprintf("(%.4f)", result.FuncB)
	}else{
		result.FuncB = 0
		if len(sortFuncList) > 0 {
			expr = expr[:len(expr)-1]
		}else{
			expr += "?"
		}
	}
	result.FuncExpr = expr
	return result
}

func getEstimate(s string) (estimate,pValue float64, level int) {
	level = strings.Count(s, "*")
	var eStr,pStr string
	var err error
	count := 0
	for _,v := range strings.Split(s, " ") {
		if v != "" {
			count = count + 1
		}
		if count == 2 && eStr == "" {
			eStr = v
		}
		if count == 5 {
			if strings.Contains(v, "<") || strings.Contains(v, ">") {
				pStr = v[1:]
			}else {
				pStr = v
			}
		}
	}
	if strings.Contains(eStr, "e") {
		decimalNum,err := decimal.NewFromString(eStr)
		if err != nil {
			log.Printf("decimal estimate error: %v \n", err)
		}else{
			eStr = decimalNum.String()
		}
	}
	estimate,err = strconv.ParseFloat(eStr, 64)
	if err != nil {
		log.Printf("parse estimate float error: %v \n", err)
	}
	if strings.Contains(pStr, "e") {
		decimalNum,err := decimal.NewFromString(pStr)
		if err != nil {
			log.Printf("decimal p value error: %v \n", err)
		}else{
			pStr = decimalNum.String()
		}
	}
	pValue,err = strconv.ParseFloat(pStr, 64)
	if err != nil {
		log.Printf("parse p value float error: %v \n", err)
	}
	return estimate,pValue,level
}

func checkRParam(param models.RRequestParam) error {
	var err error
	if len(param.Monitor.Config) == 0 && (len(param.YData) == 0 || len(param.XData) == 0 ) {
		err = fmt.Errorf("param validate fail,monitor config and data is empty")
		return err
	}
	return nil
}

func AutoJustifyData(param models.RRequestMonitor) (err error, result models.YXDataObj) {
	if param.LegendY == "" || len(param.LegendX) == 0 {
		return fmt.Errorf("param validate fail,legendY and legendX can not empty"), result
	}
	err,eChartData := MonitorChart(param.Config)
	if err != nil {
		return err,result
	}
	result.Legend = []string{"timestamp", param.LegendY}
	var xData [][][]float64
	var yData [][]float64
	for i,v := range eChartData.Legend {
		if v == param.LegendY {
			yData = eChartData.Series[i].Data
			break
		}
		for _,vv := range param.LegendX {
			if v == vv {
				xData = append(xData, eChartData.Series[i].Data)
				result.Legend = append(result.Legend, vv)
				break
			}
		}
	}
	if len(yData) < 2 {
		return fmt.Errorf("data Y length=%d is too short! ", len(yData)),result
	}
	yStep,yData := clearYXData(yData)
	var xMapList []map[float64][]float64
	for i,v := range xData {
		if len(v) < 2 {
			return fmt.Errorf("data X %s length=%d is too short! ", param.LegendX[i], len(v)),result
		}else {
			tmpXStep,tmpXData := clearYXData(v)
			if tmpXStep != yStep {
				return fmt.Errorf("data X %s step=%.1f is diff from Y step=%.1f ", tmpXStep, yStep),result
			}
			xMapList = append(xMapList, offsetYXData(yData, tmpXData, tmpXStep))
		}
	}
	for _,v := range yData {
		removeFlag := false
		for _,vv := range param.RemoveList {
			if v[0] == vv {
				removeFlag = true
				break
			}
		}
		if removeFlag {
			continue
		}
		illegalFlag := false
		tmpXList := []float64{v[0], v[1]}
		for _,vv := range xMapList {
			if _,b:=vv[v[0]];!b{
				illegalFlag = true
				break
			}else{
				tmpXList = append(tmpXList, vv[v[0]][1])
			}
		}
		if !illegalFlag {
			result.Data = append(result.Data, tmpXList)
		}
	}
	return err,result
}

func clearYXData(data [][]float64) (step float64, newData [][]float64) {
	step = 60
	dataLength := len(data)
	for i,v := range data {
		if i < dataLength-1 {
			if step < (data[i+1][0]-v[0]) {
				step = data[i+1][0]-v[0]
			}
		}
	}
	for i,v := range data {
		if i < dataLength-1 {
			if (data[i+1][0]-v[0]) == step {
				newData = append(newData, v)
			}
		}else{
			if (v[0]-data[i-1][0]) == step {
				newData = append(newData, v)
			}
		}
	}
	return step,newData
}

func offsetYXData(yData,xData [][]float64,step float64) map[float64][]float64 {
	newXData := make(map[float64][]float64)
	var offset float64 = 0
	if yData[0][0] != xData[0][0] {
		for _, v := range yData {
			fetchFlag := false
			for i, vv := range xData {
				if math.Abs(v[0]-vv[0]) < step {
					fetchFlag = true
					if i == len(xData)-1 {
						offset = v[0] - vv[0]
						continue
					}
					if math.Abs(v[0]-vv[0]) < math.Abs(v[0]-xData[i+1][0]) {
						offset = v[0] - vv[0]
					} else {
						offset = v[0] - xData[i+1][0]
					}
					break
				}
			}
			if fetchFlag {
				break
			}
		}
	}
	for _,v :=range xData {
		newT := v[0]+offset
		newXData[newT] = []float64{newT, v[1]}
	}
	return newXData
}