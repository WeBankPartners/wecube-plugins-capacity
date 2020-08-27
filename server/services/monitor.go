package services

import (
	"github.com/WeBankPartners/wecube-plugins-capacity/server/models"
	"net/http"
	"fmt"
	"bytes"
	"encoding/json"
	"io/ioutil"
)

func MonitorEndpointSearch(search string) (err error,result []models.OptionModel) {
	var response models.MonitorOptionResponse
	err,data := requestMonitor(http.MethodGet, fmt.Sprintf("dashboard/search?page=1&size=1000&search=%s", search), nil)
	if err != nil {
		return err, result
	}
	err = json.Unmarshal(data, &response)
	return err,response.Data
}

func MonitorMetricSearch(endpointType string) (err error,result []models.MetricOptionModel)  {
	var response models.MonitorMetricResponse
	err,data := requestMonitor(http.MethodGet, fmt.Sprintf("dashboard/config/metric/list?type=%s", endpointType), nil)
	if err != nil {
		return err, result
	}
	err = json.Unmarshal(data, &response)
	return err,response.Data
}

func MonitorChart(param []models.ChartConfigObj) (err error,result models.EChartOption) {
	var response models.MonitorChartResponse
	err,data := requestMonitor(http.MethodPost, "dashboard/newchart", param)
	if err != nil {
		return err,result
	}
	err = json.Unmarshal(data, &response)
	return err,response.Data
}

func requestMonitor(method,url string,postData interface{}) (err error,bodyData []byte) {
	var postBytes []byte
	if postData != nil {
		postBytes,err = json.Marshal(postData)
		if err != nil {
			return err,bodyData
		}
	}
	request,_ := http.NewRequest(method, fmt.Sprintf("%s/wecube-monitor/api/v1/%s", models.Config().DataSource.Monitor.BaseUrl, url), bytes.NewBuffer(postBytes))
	tokenValue := models.Config().DataSource.Monitor.TokenValue
	if models.Config().DataSource.Monitor.TokenKey == "Authorization" {
		tokenValue = "Bearer eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJTWVNfTU9OSVRPUiIsImlhdCI6MTU5MDExNzM5NSwidHlwZSI6ImFjY2Vzc1Rva2VuIiwiY2xpZW50VHlwZSI6IlNVQl9TWVNURU0iLCJleHAiOjE3NDU2MzczOTUsImF1dGhvcml0eSI6IltTVUJfU1lTVEVNXSJ9.soBixGZNfZKJTsm-augwGsu-fOoPuTYvZNc_VlxS8oAcZsRn4-ccSjEeAvVAso-7y0dxvzz5c_gw9iUE9LVK2w"
	}
	request.Header.Set(models.Config().DataSource.Monitor.TokenKey, tokenValue)
	request.Header.Set("X-Auth-Token", "default-token-used-in-server-side")
	request.Header.Set("Content-Type", "application/json")
	response,err := http.DefaultClient.Do(request)
	if err != nil {
		return err,bodyData
	}
	bodyData,err = ioutil.ReadAll(response.Body)
	response.Body.Close()
	return err,bodyData
}