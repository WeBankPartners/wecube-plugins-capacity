package services

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/models"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/util/log"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func RAnalyzeData(param models.RRequestParam) (err error, result models.RunScriptResult) {
	err = checkRParam(param)
	if err != nil {
		return err, result
	}
	var x [][]float64
	var y []float64
	if param.Excel.Enable {
		err, yxData := getExcelData(param.Guid, param.Excel)
		if err != nil {
			return err, result
		}
		for _, v := range yxData.Data {
			for vi, vv := range v {
				if vi == 0 {
					y = append(y, vv)
				} else {
					if len(x) < vi {
						x = append(x, []float64{vv})
					} else {
						x[vi-1] = append(x[vi-1], vv)
					}
				}
			}
		}
	}
	if len(param.Monitor.Config) > 0 {
		err, yXData := AutoJustifyData(param.Monitor)
		if err != nil {
			return err, result
		}
		for _, v := range yXData.Data {
			for vi, vv := range v {
				if vi == 1 {
					y = append(y, vv)
				}
				if vi >= 2 {
					if len(x) < (vi - 1) {
						x = append(x, []float64{vv})
					} else {
						x[vi-2] = append(x[vi-2], vv)
					}
				}
			}
		}
	}
	err, result = runRscript(x, y, param.Guid, param.MinLevel)
	if err != nil {
		return err, result
	}
	for i, v := range result.FuncX {
		if param.Excel.Enable {
			result.FuncX[i].Legend = param.Excel.LegendX[v.Index-1]
		} else {
			result.FuncX[i].Legend = param.Monitor.LegendX[v.Index-1]
		}
	}
	param.FuncX = result.FuncX
	param.FuncB = result.FuncB
	err, chart := RChartData(param, x, y)
	result.Chart = chart
	return err, result
}

func RChartData(param models.RRequestParam, x [][]float64, y []float64) (err error, result models.EChartOption) {
	var xAxis models.AxisModel
	var yAxis, newAxis models.DataSerialModel
	var newYData []float64
	result.IsDataSeries = true
	result.Legend = []string{"real", "func(y)"}
	for i, _ := range y {
		xAxis.Data = append(xAxis.Data, float64(i+1))
	}
	yAxis.Data = y
	yAxis.Name = "real"
	yAxis.Type = "line"
	for i, v := range param.FuncX {
		param.FuncX[i].Data = x[v.Index-1]
	}
	for _, v := range param.FuncX {
		if len(newYData) == 0 {
			for _, vv := range v.Data {
				newYData = append(newYData, v.Estimate*vv)
			}
		} else {
			for i, vv := range v.Data {
				newYData[i] = newYData[i] + v.Estimate*vv
			}
		}
	}
	if param.FuncB > 0 {
		if len(newYData) == 0 {
			for i := 0; i < len(y); i++ {
				newYData = append(newYData, param.FuncB)
			}
		} else {
			for i, _ := range newYData {
				newYData[i] = newYData[i] + param.FuncB
			}
		}
	}
	for _, v := range newYData {
		tmpV, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", v), 64)
		newAxis.Data = append(newAxis.Data, tmpV)
	}
	newAxis.Name = "func(y)"
	newAxis.Type = "line"
	result.DataSeries = []*models.DataSerialModel{&yAxis, &newAxis}
	result.Xaxis = xAxis
	return err, result
}

func runRscript(x [][]float64, y []float64, guid string, minLevel int) (err error, result models.RunScriptResult) {
	if len(x) == 0 || len(y) == 0 {
		err = fmt.Errorf("Run r script fail,x data and y data can not empty ")
		return err, result
	}
	var b []byte
	// build workspace
	if guid == "" {
		result.Guid = models.GetWorkspaceName()
	} else {
		result.Guid = guid
	}
	result.Workspace = fmt.Sprintf("%s/%s", models.WorkspaceDir, result.Guid)
	b, err = exec.Command("/bin/sh", "-c", fmt.Sprintf("mkdir -p %s && /bin/cp -f conf/template.r %s/template.r", result.Workspace, result.Workspace)).Output()
	if err != nil {
		log.Logger.Error("Build tmp workspace fail", log.String("workspace", result.Workspace), log.String("output", string(b)), log.Error(err))
		return err, result
	}

	// replace data
	b, err = ioutil.ReadFile(result.Workspace + "/template.r")
	if err != nil {
		log.Logger.Error(fmt.Sprintf("Replace %s/template.r data,read file fail", result.Workspace), log.Error(err))
		return err, result
	}
	var xDataString, xExpr string
	for i, v := range x {
		if i == len(x)-1 {
			xDataString = xDataString + fmt.Sprintf("x%d<-c(%s)", i+1, turnFloatListToString(v))
			xExpr = xExpr + fmt.Sprintf("x%d", i+1)
		} else {
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
		log.Logger.Error(fmt.Sprintf("Replace %s/template.r data,write file fail", result.Workspace), log.Error(err))
		return err, result
	}

	// run script
	b, err = exec.Command("/bin/sh", "-c", fmt.Sprintf("Rscript %s/template.r", result.Workspace)).Output()
	if err != nil {
		log.Logger.Error(fmt.Sprintf("Run Rscript %s/template.r fail", result.Workspace), log.String("output", string(b)), log.Error(err))
		return err, result
	}
	output := dealWithScriptOutput(string(b), minLevel)
	result.FuncX = output.FuncX
	result.FuncB = output.FuncB
	result.FuncExpr = output.FuncExpr
	result.Output = output.Output
	result.Level = output.Level
	pngDir := strings.Replace(result.Workspace, "public/", "", -1)
	result.Images = []string{pngDir + "/rp001.png", pngDir + "/rp002.png", pngDir + "/rp003.png", pngDir + "/rp004.png"}
	return err, result
}

func turnFloatListToString(data []float64) string {
	if len(data) == 0 {
		return ""
	}
	var s []string
	for _, v := range data {
		s = append(s, fmt.Sprintf("%.3f", v))
	}
	return strings.Join(s, ",")
}

func dealWithScriptOutput(output string, minLevel int) models.RunScriptResult {
	var result models.RunScriptResult
	var funcBLevel int
	var sortFuncList models.FuncXSortList
	if minLevel > 3 {
		minLevel = 0
	}
	expr := "y="
	xCount := 0
	for _, v := range strings.Split(output, "\n") {
		if strings.HasPrefix(v, "(Intercept)") {
			result.FuncB, _, funcBLevel = getEstimate(v)
			continue
		}
		if strings.HasPrefix(v, "x") {
			xCount = xCount + 1
			tmpEstimate, tmpP, tmpL := getEstimate(v)
			if tmpL >= minLevel {
				sortFuncList = append(sortFuncList, &models.FuncXObj{PValue: tmpP, Estimate: tmpEstimate, FuncName: fmt.Sprintf("x%d", xCount), Level: tmpL, Index: xCount})
			}
		}
	}
	sort.Sort(sortFuncList)
	result.Output = strings.Replace(output, "\n", "<br/>", -1)
	for i, v := range sortFuncList {
		if i == 0 {
			result.Level = v.Level
		}
		expr += fmt.Sprintf("%.4f*%s+", v.Estimate, v.FuncName)
		result.FuncX = append(result.FuncX, v)
	}
	if funcBLevel > 0 {
		expr += fmt.Sprintf("(%.4f)", result.FuncB)
	} else {
		result.FuncB = 0
		if len(sortFuncList) > 0 {
			expr = expr[:len(expr)-1]
		} else {
			expr += "?"
		}
	}
	result.FuncExpr = expr
	return result
}

func getEstimate(s string) (estimate, pValue float64, level int) {
	level = strings.Count(s, "*")
	var eStr, pStr string
	var err error
	count := 0
	for _, v := range strings.Split(s, " ") {
		if v != "" {
			count = count + 1
		}
		if count == 2 && eStr == "" {
			eStr = v
		}
		if count == 5 {
			if strings.Contains(v, "<") || strings.Contains(v, ">") {
				pStr = v[1:]
			} else {
				pStr = v
			}
		}
	}
	if strings.Contains(eStr, "e") {
		decimalNum, err := decimal.NewFromString(eStr)
		if err != nil {
			log.Logger.Error("Decimal estimate error", log.Error(err))
		} else {
			eStr = decimalNum.String()
		}
	}
	estimate, err = strconv.ParseFloat(eStr, 64)
	if err != nil {
		log.Logger.Error("Parse estimate float error", log.Error(err))
	}
	if strings.Contains(pStr, "e") {
		decimalNum, err := decimal.NewFromString(pStr)
		if err != nil {
			log.Logger.Error("Decimal p value error", log.Error(err))
		} else {
			pStr = decimalNum.String()
		}
	}
	pValue, err = strconv.ParseFloat(pStr, 64)
	if err != nil {
		log.Logger.Error("Parse p value float error", log.Error(err))
	}
	return estimate, pValue, level
}

func checkRParam(param models.RRequestParam) error {
	var err error
	if param.Excel.Enable {
		if param.Excel.LegendY == "" || len(param.Excel.LegendX) == 0 || param.Guid == "" {
			err = fmt.Errorf("param validate fail,excel guid,legend_x and legend_y can not empty")
			return err
		}
	} else {
		if len(param.Monitor.Config) == 0 && (len(param.YData) == 0 || len(param.XData) == 0) {
			err = fmt.Errorf("param validate fail,monitor config and data can not empty")
			return err
		}
	}
	return nil
}

func AutoJustifyData(param models.RRequestMonitor) (err error, result models.YXDataObj) {
	if param.LegendY == "" || len(param.LegendX) == 0 {
		return fmt.Errorf("param validate fail,legendY and legendX can not empty"), result
	}
	err, eChartData := MonitorChart(param.Config)
	if err != nil {
		return err, result
	}
	result.Legend = []string{"timestamp", param.LegendY}
	var xData [][][]float64
	var yData [][]float64
	for i, v := range eChartData.Legend {
		if v == param.LegendY {
			yData = eChartData.Series[i].Data
		}
		for _, vv := range param.LegendX {
			if v == vv {
				xData = append(xData, eChartData.Series[i].Data)
				result.Legend = append(result.Legend, vv)
				break
			}
		}
	}
	if len(yData) < 2 {
		return fmt.Errorf("data Y length=%d is too short! ", len(yData)), result
	}
	yStep, yData := clearYXData(yData)
	var xMapList []map[float64][]float64
	for i, v := range xData {
		if len(v) < 2 {
			return fmt.Errorf("data X %s length=%d is too short! ", param.LegendX[i], len(v)), result
		} else {
			tmpXStep, tmpXData := clearYXData(v)
			if tmpXStep != yStep {
				return fmt.Errorf("data X %s step=%.1f is diff from Y step=%.1f ", tmpXStep, yStep), result
			}
			xMapList = append(xMapList, offsetYXData(yData, tmpXData, tmpXStep))
		}
	}
	for _, v := range yData {
		removeFlag := false
		for _, vv := range param.RemoveList {
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
		for _, vv := range xMapList {
			if _, b := vv[v[0]]; !b {
				illegalFlag = true
				break
			} else {
				tmpXList = append(tmpXList, vv[v[0]][1])
			}
		}
		if !illegalFlag {
			result.Data = append(result.Data, tmpXList)
		}
	}
	return err, result
}

func clearYXData(data [][]float64) (step float64, newData [][]float64) {
	step = 3600000
	dataLength := len(data)
	for i, v := range data {
		if i < dataLength-1 {
			if step > (data[i+1][0] - v[0]) {
				step = data[i+1][0] - v[0]
			}
		}
	}
	for i, v := range data {
		if i < dataLength-1 {
			if (data[i+1][0] - v[0]) == step {
				newData = append(newData, v)
			}
		} else {
			if (v[0] - data[i-1][0]) == step {
				newData = append(newData, v)
			}
		}
	}
	log.Logger.Debug(fmt.Sprintf("Clear data --> len(data)=%d len(newData)=%d step=%.1f \n", dataLength, len(newData), step))
	return step, newData
}

func offsetYXData(yData, xData [][]float64, step float64) map[float64][]float64 {
	newXData := make(map[float64][]float64)
	if len(yData) == 0 || len(xData) == 0 {
		return newXData
	}
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
	for _, v := range xData {
		newT := v[0] + offset
		newXData[newT] = []float64{newT, v[1]}
	}
	return newXData
}

func SaveRWork(param models.SaveWorkParam) error {
	var err error
	// Save r_work
	workTable := models.RWorkTable{Guid: param.Guid, Name: param.Name, Workspace: param.Workspace, Output: param.Output, Expr: param.FuncExpr, FuncB: fmt.Sprintf("%.4f", param.FuncB), Level: param.Level, LegendX: strings.Join(param.Monitor.LegendX, "^"), LegendY: param.Monitor.LegendY}
	var workFuncX, workFuncXName []string
	for _, v := range param.FuncX {
		workFuncX = append(workFuncX, fmt.Sprintf("%.4f", v.Estimate))
		workFuncXName = append(workFuncXName, v.Legend)
	}
	workTable.FuncX = strings.Join(workFuncX, ",")
	workTable.FuncXName = strings.Join(workFuncXName, "^")
	err = saveRWorkTable(workTable)
	if err != nil {
		return err
	}
	// Save r_chart
	chartTable := models.RChartTableInput{Guid: param.Guid, YReal: param.YReal, YFunc: param.YFunc}
	err = saveRChartTable(chartTable)
	if err != nil {
		return err
	}
	// Save r_images
	var imagesTable []*models.RImagesTable
	for i := 1; i <= 4; i++ {
		tmpFilePath := fmt.Sprintf("%s/rp00%d.png", param.Workspace, i)
		_, tmpError := os.Stat(tmpFilePath)
		if os.IsNotExist(tmpError) {
			err = fmt.Errorf("image file %s not exist ", tmpFilePath)
			break
		}
		b, _ := ioutil.ReadFile(tmpFilePath)
		imagesTable = append(imagesTable, &models.RImagesTable{Guid: param.Guid, Workspace: param.Workspace, Data: b})
	}
	if err != nil {
		return err
	}
	err = saveRImagesTable(imagesTable)
	if err != nil {
		return err
	}
	// Save r_monitor
	err = saveRMonitorTable(param.Guid, param.Monitor)
	return err
}

func GetRWork(guid string) (err error, result models.RunScriptResult) {
	err, rWorkTables := ListRConfig(guid)
	if err != nil {
		return err, result
	}
	if len(rWorkTables) == 0 {
		return fmt.Errorf("there is no record in r_work with guid=%s ", guid), result
	}
	result.Guid = guid
	result.Workspace = rWorkTables[0].Workspace
	result.Output = rWorkTables[0].Output
	result.FuncExpr = rWorkTables[0].Expr
	result.Level = rWorkTables[0].Level
	result.LegendX = strings.Split(rWorkTables[0].FuncXName, "^")
	result.LegendY = rWorkTables[0].LegendY
	pngDir := strings.Replace(result.Workspace, "public/", "", -1)
	result.Images = []string{pngDir + "/rp001.png", pngDir + "/rp002.png", pngDir + "/rp003.png", pngDir + "/rp004.png"}
	for _, v := range strings.Split(rWorkTables[0].FuncXName, "^") {
		result.FuncX = append(result.FuncX, &models.FuncXObj{Legend: v})
	}
	isNeedCreateImage := false
	for i := 1; i <= 4; i++ {
		tmpFilePath := fmt.Sprintf("%s/rp00%d.png", result.Workspace, i)
		_, tmpError := os.Stat(tmpFilePath)
		if os.IsNotExist(tmpError) {
			isNeedCreateImage = true
			break
		}
	}
	if isNeedCreateImage {
		err, imagesTable := getRImagesTable(guid)
		if err != nil {
			return err, result
		}
		exec.Command("bash", "-c", fmt.Sprintf("mkdir -p %s && rm -f %s/*png", result.Workspace, result.Workspace)).Run()
		for i, v := range imagesTable {
			tmpErr := ioutil.WriteFile(fmt.Sprintf("%s/rp00%d.png", result.Workspace, i+1), v.Data, 0644)
			if tmpErr != nil {
				log.Logger.Error(fmt.Sprintf("Write images file=%s/rp00%d.png error", result.Workspace, i), log.Error(tmpErr))
			}
		}
	}
	err, chartTables := getRChartTable(guid)
	if err != nil {
		return err, result
	}
	if len(chartTables) == 0 {
		return fmt.Errorf("there is no chart data in r_chart with guid=%s ", guid), result
	}
	var chartOption models.EChartOption
	var yReal, yFunc models.DataSerialModel
	yReal.Name = "real"
	yReal.Type = "line"
	yReal.Data = chartTables[0].YReal
	yFunc.Name = "func(y)"
	yFunc.Type = "line"
	yFunc.Data = chartTables[0].YFunc
	chartOption.DataSeries = []*models.DataSerialModel{&yReal, &yFunc}
	chartOption.IsDataSeries = true
	chartOption.Legend = []string{"real", "func(y)"}
	var xAxis models.AxisModel
	for i, _ := range chartTables[0].YReal {
		xAxis.Data = append(xAxis.Data, float64(i+1))
	}
	chartOption.Xaxis = xAxis
	result.Chart = chartOption
	return err, result
}

func RCalcData(param models.RCalcParam) (err error, result models.RCalcResult) {
	result.Guid = param.Guid
	err, rWorkTables := ListRConfig(param.Guid)
	if err != nil {
		return fmt.Errorf("list r config error -> %v \n", err), result
	}
	if len(rWorkTables) == 0 {
		return fmt.Errorf("there is no record in r_work with guid=%s ", param.Guid), result
	}
	var yXTable models.YXDataTable
	var estimate, yList []float64
	for _, v := range strings.Split(rWorkTables[0].FuncX, ",") {
		tmpEstimate, _ := strconv.ParseFloat(v, 64)
		estimate = append(estimate, tmpEstimate)
	}
	yXTable.Title = strings.Split(rWorkTables[0].FuncXName, "^")
	yXTable.Title = append(yXTable.Title, "func(y)")
	funcB, _ := strconv.ParseFloat(rWorkTables[0].FuncB, 64)
	for i, v := range param.AddData {
		if len(v) != len(estimate) {
			err = fmt.Errorf("add_data row index %d is validate,length != len(estimate) ", i)
			break
		}
		tmpTableMap := make(map[string]string)
		var tmpY float64
		for j, vv := range v {
			tmpY += estimate[j] * vv
			tmpTableMap[yXTable.Title[j]] = fmt.Sprintf("%.4f", vv)
		}
		tmpY += funcB
		tmpTableMap["func(y)"] = fmt.Sprintf("%.4f", tmpY)
		yList = append(yList, tmpY)
		yXTable.Data = append(yXTable.Data, tmpTableMap)
	}
	if err != nil {
		return fmt.Errorf("calc add data to funcY error -> %v \n", err), result
	}
	result.Table = yXTable
	err, chartTables := getRChartTable(param.Guid)
	if err != nil {
		return fmt.Errorf("get r chart table error -> %v \n", err), result
	}
	if len(chartTables) == 0 {
		return fmt.Errorf("there is no chart data in r_chart with guid=%s ", param.Guid), result
	}
	var chartOption models.EChartOption
	var yReal, yFunc models.DataSerialModel
	yReal.Name = "real"
	yReal.Type = "line"
	yReal.Data = chartTables[0].YReal
	yFunc.Name = "func(y)"
	yFunc.Type = "line"
	yFunc.Data = append(chartTables[0].YFunc, yList...)
	chartOption.DataSeries = []*models.DataSerialModel{&yReal, &yFunc}
	chartOption.IsDataSeries = true
	chartOption.Legend = []string{"real", "func(y)"}
	var xAxis models.AxisModel
	for i, _ := range yFunc.Data {
		xAxis.Data = append(xAxis.Data, float64(i+1))
	}
	chartOption.Xaxis = xAxis
	result.Chart = chartOption
	return err, result
}

func AutoCleanWorkspace() {
	log.Logger.Info("Start auto clean useless workspace cron job")
	t := time.NewTicker(time.Duration(60) * time.Minute).C
	for {
		<-t
		workspaceDir := models.Config().Cache.WorkspaceDir
		fs, err := ioutil.ReadDir(workspaceDir)
		if err != nil {
			log.Logger.Error("Clean workspace job fail with read dir error", log.String("dir", workspaceDir), log.Error(err))
			continue
		}
		err, saveList := ListRConfig("")
		if err != nil {
			log.Logger.Error("Clean workspace job fail with get saved list error", log.Error(err))
			continue
		}
		tn := time.Now()
		for _, v := range fs {
			if !strings.HasPrefix(v.Name(), "R_") {
				continue
			}
			isSave := false
			for _, vv := range saveList {
				if vv.Guid == v.Name() {
					isSave = true
					break
				}
			}
			if isSave {
				continue
			}
			if tn.Sub(v.ModTime()).Minutes() > 60 {
				err = exec.Command("/bin/sh", "-c", fmt.Sprintf("rm -rf %s/%s", workspaceDir, v.Name())).Run()
				if err != nil {
					log.Logger.Error("Clean workspace job error", log.String("subDir", v.Name()), log.Error(err))
				} else {
					log.Logger.Info("Clean workspace useless dir success", log.String("subDir", v.Name()))
				}
			}
		}
	}
}

func SaveExcelFile(content []byte) (err error, result models.RCalcResult) {
	result.Guid = models.GetWorkspaceName()
	output, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("mkdir -p %s/%s", models.WorkspaceDir, result.Guid)).Output()
	if err != nil {
		err = fmt.Errorf("Try to create workspce dir fail,output=%s,err=%s ", string(output), err.Error())
		return err, result
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s/data.xlsx", models.WorkspaceDir, result.Guid), content, 0644)
	if err != nil {
		err = fmt.Errorf("Try to save excel file fail,%s ", err.Error())
		return err, result
	}
	err, yxData := getExcelData(result.Guid, models.RRequestExcel{})
	if err != nil {
		return err, result
	}
	result.Table.Title = append([]string{"index"}, yxData.Legend...)
	for i, v := range yxData.Data {
		rowMap := make(map[string]string)
		rowMap["index"] = strconv.Itoa(i + 1)
		rowMap["id"] = strconv.Itoa(i + 1)
		for ii, vv := range v {
			rowMap[yxData.Legend[ii]] = fmt.Sprintf("%.3f", vv)
		}
		result.Table.Data = append(result.Table.Data, rowMap)
	}
	return err, result
}

func getExcelData(guid string, config models.RRequestExcel) (err error, result models.YXDataObj) {
	excelObj, err := excelize.OpenFile(fmt.Sprintf("%s/%s/data.xlsx", models.WorkspaceDir, guid))
	if err != nil {
		err = fmt.Errorf("Try to read excel file fail,%s ", err.Error())
		return err, result
	}
	sheetList := excelObj.GetSheetList()
	if len(sheetList) == 0 {
		err = fmt.Errorf("Excel sheet is empty ")
		return err, result
	}
	rows, err := excelObj.GetRows(sheetList[0])
	if err != nil {
		err = fmt.Errorf("Excel get rows fail,%s ", err.Error())
		return err, result
	}
	var indexList []int
	for i, v := range rows {
		if i == 0 {
			if config.Enable {
				for ii, vv := range v {
					if vv == config.LegendY {
						indexList = []int{ii}
						break
					}
				}
				for _, vx := range config.LegendX {
					for ii, vv := range v {
						if vv == vx {
							indexList = append(indexList, ii)
							break
						}
					}
				}
				result.Legend = append([]string{config.LegendY}, config.LegendX...)
			} else {
				result.Legend = v
			}
		} else {
			var rowData []float64
			if config.Enable {
				removeFlag := false
				for _, removeIndex := range config.RemoveList {
					if removeIndex == i {
						removeFlag = true
						break
					}
				}
				if removeFlag {
					continue
				}
				for _, indexY := range indexList {
					for ii, vv := range v {
						if ii == indexY {
							vFloat, _ := strconv.ParseFloat(vv, 64)
							rowData = append(rowData, vFloat)
							break
						}
					}
				}
			} else {
				if len(v) != len(result.Legend) {
					err = fmt.Errorf("Excel sheet 1 row %d width is not equal the title ", i+1)
					break
				}
				for ii, vv := range v {
					vFloat, parseErr := strconv.ParseFloat(vv, 64)
					if parseErr != nil {
						err = fmt.Errorf("Excel sheet 1 row %d num %d data validate fail ", i+1, ii+1)
						break
					}
					rowData = append(rowData, vFloat)
				}
				if err != nil {
					break
				}
			}
			result.Data = append(result.Data, rowData)
		}
	}
	return err, result
}
