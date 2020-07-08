package api

import (
	"net/http"
	"log"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/api/v1/capacity"
)

func InitHttpServer(port int) {
	http.Handle("/capacity/api/v1/monitor/search", http.HandlerFunc(capacity.MonitorSearchHandler))
	http.Handle("/capacity/api/v1/monitor/chart", http.HandlerFunc(capacity.MonitorDataHandler))
	http.Handle("/capacity/api/v1/r/data", http.HandlerFunc(capacity.RJustifyDataHandler))
	http.Handle("/capacity/api/v1/r/analyze", http.HandlerFunc(capacity.RAnalyzeHandler))
	http.Handle("/capacity/api/v1/r/calc", http.HandlerFunc(capacity.RCalcDataHandle))
	http.Handle("/capacity/api/v1/r/save", http.HandlerFunc(capacity.SaveAnalyzeConfig))
	http.Handle("/capacity/api/v1/r/list", http.HandlerFunc(capacity.ListAnalyzeConfig))
	http.Handle("/capacity/api/v1/r/get", http.HandlerFunc(capacity.GetAnalyzeConfig))
	http.Handle("/capacity/api/v1/r/delete", http.HandlerFunc(capacity.DeleteAnalyzeConfig))
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/capacity/", http.StripPrefix("/capacity/", fs))
	listenPort := fmt.Sprintf(":%d", port)
	log.Printf("listening %s ...\n", listenPort)
	http.ListenAndServe(listenPort, nil)
}