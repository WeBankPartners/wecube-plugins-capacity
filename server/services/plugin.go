package services

import (
	"github.com/WeBankPartners/wecube-plugins-capacity/server/models"
	"fmt"
	"time"
	"strings"
)

func CompareModels(input *models.PluginRequestInput) (err error,output models.PluginResponseOutput) {
	defer func() {
		if err != nil {
			output.ErrorCode = "1"
			output.ErrorMessage = err.Error()
		}else{
			output.ErrorCode = "0"
			output.ErrorMessage = ""
		}
		output.Guid = input.Guid
		output.CallbackParameter = input.CallbackParameter
	}()
	// Check param
	if input.TemplateName == "" {
		err = fmt.Errorf("Parameter template_name can not empty ")
		return
	}
	if input.Start == "" || input.End == "" {
		err = fmt.Errorf("Parameter start and end can not empty ")
		return
	}
	startT,cErr := time.Parse("2006-01-02 15:04:05", input.Start)
	if cErr != nil {
		err = fmt.Errorf("Parameter start validate fail,ensure format like '2006-01-02 15:04:05' ")
		return
	}
	endT,cErr := time.Parse("2006-01-02 15:04:05", input.End)
	if cErr != nil {
		err = fmt.Errorf("Parameter end validate fail,ensure format like '2006-01-02 15:04:05' ")
		return
	}

	// Load template work
	cErr,rWorks := getRWorkByName(input.TemplateName)
	if len(rWorks) == 0 {
		if cErr != nil {
			err = fmt.Errorf("Get r_work table data fail,%s ", cErr.Error())
		}else{
			err = fmt.Errorf("Get empty data from r_work when name=%s ", input.TemplateName)
		}
		return
	}
	if rWorks[0].LegendY == "" || rWorks[0].LegendX == "" {
		err = fmt.Errorf("Please check the legend_x and legned_y of this template,can not empty ")
		return
	}
	output.FuncOld = rWorks[0].Expr
	output.LevelOld = fmt.Sprintf("%d", rWorks[0].Level)

	// Load template monitor
	cErr,rMonitor := getRMonitorTable(rWorks[0].Guid)
	if cErr != nil {
		err = fmt.Errorf("Get r_monitor table data fail,%s ", cErr.Error())
		return
	}
	if len(rMonitor) == 0 {
		err = fmt.Errorf("Please make sure the template model build with monitor data ")
		return
	}

	// Compute new r func
	var monitorParam models.RRequestMonitor
	for _,v := range rMonitor {
		monitorParam.Config = append(monitorParam.Config, models.ChartConfigObj{Endpoint:v.Endpoint, Metric:v.Metric, Aggregate:v.Agg, Start:fmt.Sprintf("%d",startT.Unix()), End:fmt.Sprintf("%d",endT.Unix())})
	}
	monitorParam.LegendY = rWorks[0].LegendY
	monitorParam.LegendX = strings.Split(rWorks[0].LegendX, "^")
	cErr,analyzeResult := RAnalyzeData(models.RRequestParam{Monitor:monitorParam})
	if cErr != nil {
		err = fmt.Errorf("Analyze new model fail,%s ", cErr.Error())
		return
	}
	output.FuncNew = analyzeResult.FuncExpr
	output.LevelNew = fmt.Sprintf("%d", analyzeResult.Level)

	// Save
	isSave := strings.ToLower(input.Save)
	if isSave == "y" || isSave == "yes" || isSave == "true" {
		var saveParam models.SaveWorkParam
		saveParam.Guid = analyzeResult.Guid
		saveParam.FuncExpr = analyzeResult.FuncExpr
		saveParam.Output = analyzeResult.Output
		saveParam.Level = analyzeResult.Level
		saveParam.Workspace = analyzeResult.Workspace
		saveParam.FuncX = analyzeResult.FuncX
		saveParam.FuncB = analyzeResult.FuncB
		saveParam.Images = analyzeResult.Images
		if analyzeResult.Chart.DataSeries[0].Name == "real" {
			saveParam.YReal = analyzeResult.Chart.DataSeries[0].Data
			saveParam.YFunc = analyzeResult.Chart.DataSeries[1].Data
		}else{
			saveParam.YReal = analyzeResult.Chart.DataSeries[1].Data
			saveParam.YFunc = analyzeResult.Chart.DataSeries[0].Data
		}
		saveParam.Monitor = monitorParam
		tmpName := rWorks[0].Name
		if strings.Contains(rWorks[0].Name, "-") {
			tmpName = tmpName[:strings.LastIndex(tmpName, "-")]
		}
		tmpName += fmt.Sprintf("-%s", time.Now().Format("20060102150405"))
		saveParam.Name = tmpName
		output.TemplateName = tmpName
		err = SaveRWork(saveParam)
	}
	return
}