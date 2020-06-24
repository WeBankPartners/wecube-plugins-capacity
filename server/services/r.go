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

func RAnalyzeData()  {

}

func RImageData()  {

}

func RChartData()  {

}
func runRscript(x,y []float64) (err error,result models.RunScriptResult)  {
	var b []byte
	// build workspace
	result.Guid = models.GetWorkspaceName()
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
	result = dealWithScriptOutput(string(b))
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
			result.FuncA,result.Level = getEstimate(v)
			continue
		}
		if strings.HasPrefix(v, "x") {
			result.FuncB,_ = getEstimate(v)
		}
	}
	result.Output = strings.Replace(output, "\n", "<br/>", -1)
	result.FuncExpr = fmt.Sprintf("y=%.4fx+(%.4f)", result.FuncA, result.FuncB)
	result.Images = []string{result.Workspace+"/rp001.png",result.Workspace+"/rp002.png",result.Workspace+"/rp003.png",result.Workspace+"/rp004.png"}
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