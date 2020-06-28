package api

import (
	"net/http"
	"log"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/api/v1/capacity"
)

func InitHttpServer(port int) {
	http.Handle("/capacity/monitor/search", http.HandlerFunc(capacity.MonitorSearchHandler))
	http.Handle("/capacity/monitor/chart", http.HandlerFunc(capacity.MonitorDataHandler))
	http.Handle("/capacity/r/analyze", http.HandlerFunc(capacity.RAnalyzeHandler))
	http.Handle("/capacity/r/image", http.HandlerFunc(capacity.RPlotImageHandle))
	http.Handle("/capacity/r/chart", http.HandlerFunc(capacity.RDataChartHandle))
	listenPort := fmt.Sprintf(":%d", port)
	log.Printf("listening %s ...\n", listenPort)
	http.ListenAndServe(listenPort, nil)
}