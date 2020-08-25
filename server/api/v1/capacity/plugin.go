package capacity

import (
	"net/http"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/models"
	"io/ioutil"
	"encoding/json"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/services"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/util/log"
)

func ComparePluginHandler(w http.ResponseWriter,r *http.Request)  {
	var param models.PluginRequest
	b,_ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(b, &param)
	if err != nil {
		returnJson(r,w,err,nil)
		return
	}
	var result models.PluginResponse
	var outputs []*models.PluginResponseOutput
	for _,v := range param.Inputs {
		tmpErr,tmpOutput := services.CompareModels(v)
		if tmpErr != nil {
			log.Logger.Error("Handle compare plugin fail", log.Error(tmpErr))
			err = tmpErr
		}
		outputs = append(outputs, &tmpOutput)
	}
	result.Results = models.PluginResponseOutputs{Outputs:outputs}
	if err != nil {
		result.ResultCode = "1"
		result.ResultMessage = err.Error()
	}else{
		result.ResultCode = "0"
		result.ResultMessage = "success"
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	d,_ := json.Marshal(result)
	w.Write(d)
}
