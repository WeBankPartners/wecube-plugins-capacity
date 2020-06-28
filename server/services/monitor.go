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
	err,data := requestMonitor(http.MethodGet, fmt.Sprintf("dashboard/search?page=1&size=1000&search=%s", search), nil)
	if err != nil {
		return err, result
	}
	err = json.Unmarshal(data, &result)
	return err,result
}

func MonitorMetricSearch(endpointType string) (err error,result []models.OptionModel)  {
	err,data := requestMonitor(http.MethodGet, fmt.Sprintf("dashboard/config/metric/list?type=%s", endpointType), nil)
	if err != nil {
		return err, result
	}
	err = json.Unmarshal(data, &result)
	return err,result
}

func MonitorChart(param []models.ChartConfigObj) (err error,result models.EChartOption) {
	err,data := requestMonitor(http.MethodPost, "dashboard/newchart", param)
	if err != nil {
		return err,result
	}
	err = json.Unmarshal(data, &result)
	return err,result
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
	request.Header.Set("X-Auth-Token", models.Config().DataSource.Monitor.Token)
	request.Header.Set("Content-Type", "application/json")
	response,err := http.DefaultClient.Do(request)
	if err != nil {
		return err,bodyData
	}
	bodyData,err = ioutil.ReadAll(response.Body)
	response.Body.Close()
	return err,bodyData
}