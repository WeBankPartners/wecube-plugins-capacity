package capacity

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/models"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/services"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/util/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func MonitorSearchHandler(w http.ResponseWriter, r *http.Request) {
	searchText := r.FormValue("search")
	searchType := r.FormValue("search_type")
	endpointType := r.FormValue("type")
	if searchType != "endpoint" && searchType != "metric" {
		returnJson(r, w, fmt.Errorf("validate fail,param search_type must be endpoint or metric "), nil)
		return
	}
	if searchType == "metric" {
		if endpointType == "" {
			returnJson(r, w, fmt.Errorf("validate fail,param type can not empty "), nil)
			return
		}
		err, result := services.MonitorMetricSearch(endpointType)
		returnJson(r, w, err, result)
	} else {
		err, result := services.MonitorEndpointSearch(searchText)
		returnJson(r, w, err, result)
	}
}

func MonitorDataHandler(w http.ResponseWriter, r *http.Request) {
	var param []models.ChartConfigObj
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(b, &param)
	if err != nil {
		returnJson(r, w, err, nil)
		return
	}
	err, result := services.MonitorChart(param)
	returnJson(r, w, err, result)
}

func ExcelUploadHandler(w http.ResponseWriter, r *http.Request) {
	var uploadBytes []byte
	multipartFile, _, err := r.FormFile("file")
	if err != nil {
		log.Logger.Error("accept form file fail", log.Error(err))
		returnJson(r, w, err, nil)
		return
	}
	uploadBytes, err = ioutil.ReadAll(multipartFile)
	if err != nil {
		returnJson(r, w, fmt.Errorf("Read upload file fail,%s ", err.Error()), nil)
		return
	}
	if len(uploadBytes) == 0 {
		returnJson(r, w, fmt.Errorf("Upload file is empty "), nil)
		return
	}
	err, result := services.SaveExcelFile(uploadBytes)
	returnJson(r, w, err, result)
}

func RJustifyDataHandler(w http.ResponseWriter, r *http.Request) {
	var param models.RRequestMonitor
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(b, &param)
	if err != nil {
		returnJson(r, w, err, nil)
		return
	}
	if len(param.Config) == 0 || param.LegendY == "" || len(param.LegendX) == 0 {
		err = fmt.Errorf("Param validate fail,config legendY legendX can not empty ")
		returnJson(r, w, err, nil)
		return
	}
	err, result := services.AutoJustifyData(param)
	if err != nil {
		returnJson(r, w, err, result)
		return
	}
	var tableData models.YXDataTable
	tableData.Title = []string{"index", "time"}
	for i, v := range result.Legend {
		if i == 0 {
			continue
		}
		tableData.Title = append(tableData.Title, v)
	}
	for index, v := range result.Data {
		tmpMap := make(map[string]string)
		tmpMap["index"] = strconv.Itoa(index + 1)
		for i, vv := range v {
			if i == 0 {
				tmpMap["time"] = time.Unix(int64(vv/1000), 0).Format("2006-01-02 15:04:05")
				tmpMap["id"] = fmt.Sprintf("%.0f", vv)
				tmpMap[result.Legend[i]] = fmt.Sprintf("%.0f", vv)
			} else {
				tmpMap[result.Legend[i]] = fmt.Sprintf("%.4f", vv)
			}
		}
		tableData.Data = append(tableData.Data, tmpMap)
	}
	returnJson(r, w, err, tableData)
}

func RAnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	var param models.RRequestParam
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(b, &param)
	if err != nil {
		returnJson(r, w, err, nil)
		return
	}
	err, result := services.RAnalyzeData(param)
	returnJson(r, w, err, result)
}

func RCalcDataHandle(w http.ResponseWriter, r *http.Request) {
	var param models.RCalcParam
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(b, &param)
	if err != nil {
		returnJson(r, w, err, nil)
		return
	}
	if len(param.AddData) == 0 || param.Guid == "" {
		err = fmt.Errorf("Param validate fail,add_data and guid can not empty ")
		returnJson(r, w, err, nil)
		return
	}
	err, result := services.RCalcData(param)
	returnJson(r, w, err, result)
}

func SaveAnalyzeConfig(w http.ResponseWriter, r *http.Request) {
	var param models.SaveWorkParam
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(b, &param)
	if err != nil {
		returnJson(r, w, err, nil)
		return
	}
	if param.Guid == "" || param.Name == "" || param.Workspace == "" {
		err = fmt.Errorf("Param validate fail,guid name workspace can not empty ")
	} else {
		param.Name = fmt.Sprintf("%s-%s", param.Name, time.Now().Format("20060102150405"))
		err = services.SaveRWork(param)
	}
	returnJson(r, w, err, nil)
}

func GetAnalyzeConfig(w http.ResponseWriter, r *http.Request) {
	guid := r.FormValue("guid")
	err, result := services.GetRWork(guid)
	returnJson(r, w, err, result)
}

func ListAnalyzeConfig(w http.ResponseWriter, r *http.Request) {
	guid := r.FormValue("guid")
	err, result := services.ListRConfig(guid)
	returnJson(r, w, err, result)
}

func DeleteAnalyzeConfig(w http.ResponseWriter, r *http.Request) {
	var err error
	guid := r.FormValue("guid")
	if guid == "" {
		err = fmt.Errorf("Param validate fail,guid can not empty ")
	} else {
		err = services.DeleteRConfig(guid)
	}
	returnJson(r, w, err, nil)
}

func returnJson(r *http.Request, w http.ResponseWriter, err error, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	//w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, Authorization, authorization, Token, X-Auth-Token")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS")
	var response models.RespJson
	if err != nil {
		log.Logger.Error("Request fail", log.String("url", r.URL.String()), log.Error(err))
		response.Code = 1
		response.Msg = err.Error()
	} else {
		response.Code = 0
		response.Msg = "success"
	}
	response.Data = result
	w.WriteHeader(http.StatusOK)
	d, _ := json.Marshal(response)
	w.Write(d)
}

func ReturnAuthFail(r *http.Request, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	var response models.RespJson
	response.Code = 1
	response.Msg = "Token validate fail"
	response.Data = "Token validate fail"
	w.WriteHeader(http.StatusUnauthorized)
	d, _ := json.Marshal(response)
	w.Write(d)
}

func ExportExprResult(w http.ResponseWriter, r *http.Request) {
	metric := r.FormValue("metric")
	endpointList := strings.Split(r.FormValue("endpoints"), ",")
	if len(endpointList) == 0 {
		returnJson(r, w, fmt.Errorf("Param validate fail,endpoint can not emtpy "), nil)
		return
	}
	result, err := services.ExportExprResult(endpointList, metric)
	returnJson(r, w, err, result)
}
