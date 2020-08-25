package api

import (
	"net/http"
	"fmt"
	"time"
	"strings"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/api/v1/capacity"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/util/log"
)

func InitHttpServer(port int) {
	http.Handle("/capacity/api/v1/monitor/search", handleWithLog(capacity.MonitorSearchHandler))
	http.Handle("/capacity/api/v1/monitor/chart", handleWithLog(capacity.MonitorDataHandler))
	http.Handle("/capacity/api/v1/r/data", handleWithLog(capacity.RJustifyDataHandler))
	http.Handle("/capacity/api/v1/r/analyze", handleWithLog(capacity.RAnalyzeHandler))
	http.Handle("/capacity/api/v1/r/calc", handleWithLog(capacity.RCalcDataHandle))
	http.Handle("/capacity/api/v1/r/save", handleWithLog(capacity.SaveAnalyzeConfig))
	http.Handle("/capacity/api/v1/r/list", handleWithLog(capacity.ListAnalyzeConfig))
	http.Handle("/capacity/api/v1/r/get", handleWithLog(capacity.GetAnalyzeConfig))
	http.Handle("/capacity/api/v1/r/delete", handleWithLog(capacity.DeleteAnalyzeConfig))
	http.Handle("/capacity/api/v1/r/excel", handleWithLog(capacity.ExcelUploadHandler))
	http.Handle("/capacity/api/v1/plugin/compare", handleWithLog(capacity.ComparePluginHandler))
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/capacity/", http.StripPrefix("/capacity/", fs))
	listenPort := fmt.Sprintf(":%d", port)
	log.Logger.Info(fmt.Sprintf("listening %s ...", listenPort))
	http.ListenAndServe(listenPort, nil)
}

func handleWithLog(h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w,r)
		log.Logger.Info("Request",log.String("url", r.RequestURI), log.String("method",r.Method), log.String("ip",strings.Split(r.RemoteAddr,":")[0]), log.Float64("cost_second",time.Now().Sub(start).Seconds()))
	})
}