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
)

func RAnalyzeData(param models.RRequestParam) (err error,result models.RunScriptResult) {
	err = checkRParam(param)
	if err != nil {
		return err,result
	}
	var x,y []float64
	var eChartData models.EChartOption
	if len(param.Monitor.Config) > 0 {
		err,eChartData = MonitorChart(param.Monitor.Config)
		if err != nil {
			return err,result
		}
		var xData,yData [][]float64
		for i,v := range eChartData.Legend {
			if v == param.Monitor.LegendX {
				xData = eChartData.Series[i].Data
			}
			if v == param.Monitor.LegendY {
				yData = eChartData.Series[i].Data
			}
		}
		for i,v := range yData {
			if param.Monitor.XTime {
				x = append(x, float64(i)+1)
			}
			y = append(y, v[1])
		}
		if !param.Monitor.XTime {
			for _,v := range xData {
				x = append(x, v[1])
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
	param.FuncA = result.FuncA
	param.FuncB = result.FuncB
	err,chart := RChartData(param, eChartData)
	result.Chart = chart
	return err,result
}

func RChartData(param models.RRequestParam,chart models.EChartOption) (err error,result models.EChartOption) {
	err = checkRParam(param)
	if err != nil {
		return err,result
	}
	result.IsDataSeries = true
	result.Legend = []string{"real", "fun(y)"}
	var series,newSeries models.TimeSerialModel
	var xAxis models.AxisModel
	var yAxis,newAxis models.DataSerialModel
	if len(param.Monitor.Config) > 0 {
		var eChartData models.EChartOption
		if len(chart.Legend) > 0 {
			eChartData = chart
		}else {
			err, eChartData = MonitorChart(param.Monitor.Config)
			if err != nil {
				return err, result
			}
		}
		var xData,yData [][]float64
		for i,v := range eChartData.Legend {
			if v == param.Monitor.LegendX {
				xData = eChartData.Series[i].Data
			}
			if v == param.Monitor.LegendY {
				yData = eChartData.Series[i].Data
			}
		}
		if param.Monitor.XTime {
			for i,v := range yData {
				series.Data = append(series.Data, v)
				newSeries.Data = append(newSeries.Data, []float64{v[0], param.FuncA*(float64(i)+1) + param.FuncB})
			}
			if param.AddDate > 0 {
				tmpStep := yData[1][0] - yData[0][0]
				if param.AddDate < tmpStep {
					err = fmt.Errorf("param add date less then y data step")
					return err,result
				}
				tmpLastDate := yData[len(yData)-1][0]
				tmpMaxDate := tmpLastDate + param.AddDate
				countIndex := len(yData)
				for {
					tmpLastDate = tmpLastDate + tmpStep
					if tmpLastDate > tmpMaxDate {
						break
					}
					countIndex += 1
					newSeries.Data = append(newSeries.Data, []float64{tmpLastDate, param.FuncA*float64(countIndex) + param.FuncB})
				}
				result.IsDataSeries = false
			}
			series.Name = "real"
			series.Type = "line"
			newSeries.Name = "fun(y)"
			newSeries.Type = "line"
			result.Series = []*models.TimeSerialModel{&series, &newSeries}
		}else{
			for _,v := range xData {
				xAxis.Data = append(xAxis.Data, v[1])
				newAxis.Data = append(newAxis.Data, param.FuncA*v[1] + param.FuncB)
			}
			for _,v := range yData {
				yAxis.Data = append(yAxis.Data, v[1])
			}
			for _,v := range param.AddData {
				xAxis.Data = append(xAxis.Data, v)
				newAxis.Data = append(newAxis.Data, param.FuncA*v + param.FuncB)
			}
			result.Xaxis = xAxis
			yAxis.Name = "real"
			yAxis.Type = "line"
			newAxis.Name = "fun(y)"
			newAxis.Type = "line"
			result.DataSeries = []*models.DataSerialModel{&yAxis, &newAxis}
		}
	}else{
		for _,v := range param.XData {
			xAxis.Data = append(xAxis.Data, v)
			newAxis.Data = append(newAxis.Data, param.FuncA*v + param.FuncB)
		}
		for _,v := range param.YData {
			yAxis.Data = append(yAxis.Data, v)
		}
		for _,v := range param.AddData {
			xAxis.Data = append(xAxis.Data, v)
			newAxis.Data = append(newAxis.Data, param.FuncA*v + param.FuncB)
		}
		result.Xaxis = xAxis
		yAxis.Name = "real"
		yAxis.Type = "line"
		newAxis.Name = "fun(y)"
		newAxis.Type = "line"
		result.DataSeries = []*models.DataSerialModel{&yAxis, &newAxis}
	}
	return err,result
}

func runRscript(x,y []float64,guid string) (err error,result models.RunScriptResult)  {
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
	if len(x) != len(y) {
		if len(x) > len(y) {
			x = x[:len(y)]
		}else{
			y = y[:len(x)]
		}
	}
	tData := strings.Replace(string(b), "{x_data}", turnFloatListToString(x), -1)
	tData = strings.Replace(tData, "{y_data}", turnFloatListToString(y), -1)
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
	result.FuncA = output.FuncA
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
	for _,v := range strings.Split(output, "\n") {
		if strings.HasPrefix(v, "(Intercept)") {
			result.FuncB,_ = getEstimate(v)
			continue
		}
		if strings.HasPrefix(v, "x") {
			result.FuncA,result.Level = getEstimate(v)
		}
	}
	result.Output = strings.Replace(output, "\n", "<br/>", -1)
	result.FuncExpr = fmt.Sprintf("y=%.4fx+(%.4f)", result.FuncA, result.FuncB)
	return result
}

func getEstimate(s string) (estimate float64, level int) {
	level = strings.Count(s, "*")
	var eStr string
	var err error
	for i,v := range strings.Split(s, " ") {
		if i > 1 && v != "" {
			eStr = v
			break
		}
	}
	if strings.Contains(eStr, "e") {
		decimalNum,err := decimal.NewFromString(eStr)
		if err != nil {
			log.Printf("decimal error: %v \n", err)
		}else{
			eStr = decimalNum.String()
		}
	}
	estimate,err = strconv.ParseFloat(eStr, 64)
	if err != nil {
		log.Printf("parse float error: %v \n", err)
	}
	return estimate,level
}

func checkRParam(param models.RRequestParam) error {
	var err error
	if len(param.Monitor.Config) == 0 && (len(param.YData) == 0 || len(param.XData) == 0 ) {
		err = fmt.Errorf("param validate fail,monitor config and data is empty")
		return err
	}
	return nil
}