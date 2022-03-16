package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/models"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func MonitorEndpointSearch(search string) (err error, result []models.OptionModel) {
	var response models.MonitorOptionResponse
	err, data := requestMonitor(http.MethodGet, fmt.Sprintf("dashboard/search?page=1&size=1000&search=%s", search), nil)
	if err != nil {
		return err, result
	}
	err = json.Unmarshal(data, &response)
	return err, response.Data
}

func MonitorMetricSearch(endpointType string) (err error, result []models.MetricOptionModel) {
	var response models.MonitorMetricResponse
	err, data := requestMonitor(http.MethodGet, fmt.Sprintf("dashboard/config/metric/list?type=%s", endpointType), nil)
	if err != nil {
		return err, result
	}
	err = json.Unmarshal(data, &response)
	return err, response.Data
}

func MonitorChart(param []models.ChartConfigObj) (err error, result models.EChartOption) {
	var response models.MonitorChartResponse
	err, data := requestMonitor(http.MethodPost, "dashboard/chart", transChartConfig(param))
	if err != nil {
		return err, result
	}
	err = json.Unmarshal(data, &response)
	return err, response.Data
}

func transChartConfig(param []models.ChartConfigObj) (output models.ChartQueryParam) {
	output = models.ChartQueryParam{}
	if len(param) == 0 {
		return
	}
	startInt, _ := strconv.Atoi(param[0].Start)
	endInt, _ := strconv.Atoi(param[0].End)
	output.Start = int64(startInt)
	output.End = int64(endInt)
	for _, v := range param {
		output.Data = append(output.Data, &models.ChartQueryConfigObj{Endpoint: v.Endpoint, Metric: v.Metric})
	}
	return
}

func requestMonitor(method, url string, postData interface{}) (err error, bodyData []byte) {
	var postBytes []byte
	if postData != nil {
		postBytes, err = json.Marshal(postData)
		if err != nil {
			return err, bodyData
		}
	}
	request, _ := http.NewRequest(method, fmt.Sprintf("%s/monitor/api/v1/%s", models.Config().DataSource.Monitor.BaseUrl, url), bytes.NewBuffer(postBytes))
	tokenValue := models.Config().DataSource.Monitor.TokenValue
	if models.Config().DataSource.Monitor.TokenKey == "Authorization" {
		tokenValue = models.GetCoreToken()
	}
	request.Header.Set(models.Config().DataSource.Monitor.TokenKey, tokenValue)
	request.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err, bodyData
	}
	bodyData, err = ioutil.ReadAll(response.Body)
	response.Body.Close()
	return err, bodyData
}

func ExportExprResult(endpoints []string, metric string) (result []*models.ExportResultObj, err error) {
	result = []*models.ExportResultObj{}
	specSql, endpointsParam := createListParams(endpoints, "")
	var rMonitorTables []*models.RMonitorTable
	err = mysqlEngine.SQL("select distinct guid,endpoint from r_monitor where endpoint in ("+specSql+") order by endpoint", endpointsParam...).Find(&rMonitorTables)
	if err != nil {
		err = fmt.Errorf("Try to query rMonitor table fail,%s ", err.Error())
		return
	}
	if len(rMonitorTables) == 0 {
		return
	}
	rMonitorGuidList := []string{}
	for _, v := range rMonitorTables {
		rMonitorGuidList = append(rMonitorGuidList, v.Guid)
	}
	var rWorkTables []*models.RWorkTable
	err = mysqlEngine.SQL("select guid,name,expr,func_x,func_x_name,func_b,legend_x,legend_y,update_at from r_work where guid in ('" + strings.Join(rMonitorGuidList, "','") + "')").Find(&rWorkTables)
	if err != nil {
		err = fmt.Errorf("Try to query rWork table fail,%s ", err.Error())
		return
	}
	rWorkMap := make(map[string]*models.RWorkTable)
	for _, v := range rWorkTables {
		rWorkMap[v.Guid] = v
	}
	for _, rMonitor := range rMonitorTables {
		if _, b := rWorkMap[rMonitor.Guid]; !b {
			continue
		}
		tmpLegendY := getLegendMetric(rMonitor.Endpoint, rWorkMap[rMonitor.Guid].LegendY)
		if metric != "" {
			if metric != tmpLegendY {
				continue
			}
		}
		tmpResultObj := models.ExportResultObj{RWorkName: rWorkMap[rMonitor.Guid].Name, Endpoint: rMonitor.Endpoint, FuncExpr: rWorkMap[rMonitor.Guid].Expr, UpdateTime: rWorkMap[rMonitor.Guid].UpdateAt.Format(models.DatetimeFormat)}
		tmpResultObj.YFunc = models.ExportResultParamObj{Metric: tmpLegendY}
		tmpResultObj.XParams = buildExportLegendX(rWorkMap[rMonitor.Guid].FuncX, rWorkMap[rMonitor.Guid].FuncXName, rWorkMap[rMonitor.Guid].FuncB, rMonitor.Endpoint)
		result = append(result, &tmpResultObj)
		delete(rWorkMap, rMonitor.Guid)
	}
	return
}

func getLegendMetric(endpoint, legend string) string {
	if strings.HasPrefix(legend, endpoint) {
		legend = legend[len(endpoint)+1:]
	}
	return legend
}

func buildExportLegendX(estimate, funcX, funcB, endpoint string) []*models.ExportResultParamObj {
	result := []*models.ExportResultParamObj{}
	funcXList := strings.Split(funcX, "^")
	for i, v := range strings.Split(estimate, ",") {
		tmpFloatV, _ := strconv.ParseFloat(v, 64)
		tmpObj := models.ExportResultParamObj{Name: fmt.Sprintf("x%d", i+1), Estimate: tmpFloatV}
		tmpObj.Metric = getLegendMetric(endpoint, funcXList[i])
		result = append(result, &tmpObj)
	}
	floatB, _ := strconv.ParseFloat(funcB, 64)
	result = append(result, &models.ExportResultParamObj{Name: "b", Estimate: floatB})
	return result
}

func createListParams(inputList []string, prefix string) (specSql string, paramList []interface{}) {
	if len(inputList) > 0 {
		var specList []string
		for _, v := range inputList {
			specList = append(specList, "?")
			paramList = append(paramList, prefix+v)
		}
		specSql = strings.Join(specList, ",")
	}
	return
}
