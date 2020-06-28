package capacity

import (
	"net/http"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/services"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/models"
	"encoding/json"
	"log"
	"strings"
	"fmt"
	"io/ioutil"
)

func MonitorSearchHandler(w http.ResponseWriter,r *http.Request)  {
	var err error
	var result []models.OptionModel
	searchText := r.FormValue("search")
	searchType := r.FormValue("search_type")
	endpointType := r.FormValue("type")
	if searchType != "endpoint" && searchType != "metric" {
		returnJson(r,w,fmt.Errorf("validate fail,param search_type must be endpoint or metric "),nil)
		return
	}
	if searchType == "metric" {
		if endpointType == "" {
			returnJson(r,w,fmt.Errorf("validate fail,param type can not empty "),nil)
			return
		}
		err,result = services.MonitorMetricSearch(endpointType)
	}else {
		err, result = services.MonitorEndpointSearch(searchText)
	}
	returnJson(r,w,err,result)
}

func MonitorDataHandler(w http.ResponseWriter,r *http.Request)  {
	var param []models.ChartConfigObj
	b,_ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(b, &param)
	if err != nil {
		returnJson(r,w,err,nil)
		return
	}
	err,result := services.MonitorChart(param)
	returnJson(r,w,err,result)
}

func RAnalyzeHandler(w http.ResponseWriter,r *http.Request)  {
	var param models.RRequestParam
	b,_ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(b, &param)
	if err != nil {
		returnJson(r,w,err,nil)
		return
	}
	err,result := services.RAnalyzeData(param)
	returnJson(r,w,err,result)
}

func RPlotImageHandle(w http.ResponseWriter,r *http.Request)  {

}

func RDataChartHandle(w http.ResponseWriter,r *http.Request)  {
	var param models.RRequestParam
	b,_ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(b, &param)
	if err != nil {
		returnJson(r,w,err,nil)
		return
	}
	err,result := services.RChartData(param)
	returnJson(r,w,err,result)
}

func returnJson(r *http.Request,w http.ResponseWriter,err error,result interface{})  {
	w.Header().Set("Content-Type", "application/json")
	var response models.RespJson
	if err != nil {
		log.Printf(" %s  fail,error:%s \n", r.URL.String(), err.Error())
		response.Code = 1
		response.Msg = err.Error()
		if strings.Contains(response.Msg, "validate") {
			w.WriteHeader(http.StatusBadRequest)
		}else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}else{
		log.Printf(" %s success! \n", r.URL.String())
		response.Code = 0
		response.Msg = "success"
		response.Data = result
		w.WriteHeader(http.StatusOK)
	}
	d,_ := json.Marshal(response)
	w.Write(d)
}