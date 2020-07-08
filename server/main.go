package main

import (
	"flag"
	"log"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/models"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/api"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/services"
)

func main() {
	cfgFile := flag.String("c", "conf/default.json", "config file")
	flag.Parse()
	err := models.InitConfig(*cfgFile)
	if err != nil {
		log.Printf("init config fail : %v \n", err)
		return
	}
	services.InitDbEngine()
	go services.AutoCleanWorkspace()
	api.InitHttpServer(models.Config().Http.Port)
}