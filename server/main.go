package main

import (
	"fmt"
	"flag"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/models"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/api"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/services"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/util/log"
)

func main() {
	cfgFile := flag.String("c", "conf/default.json", "config file")
	flag.Parse()
	err := models.InitConfig(*cfgFile)
	if err != nil {
		fmt.Printf("Init config fail : %v \n", err)
		return
	}
	log.InitArchiveZapLogger()
	services.InitDbEngine()
	go services.AutoCleanWorkspace()
	api.InitHttpServer(models.Config().Http.Port)
}