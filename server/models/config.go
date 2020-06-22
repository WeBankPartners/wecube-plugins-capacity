package models

import (
	"sync"
	"log"
	"encoding/json"
	"os"
	"io/ioutil"
	"strings"
)

type HttpConfig struct {
	Port  string  `json:"port"`
	Token  string  `json:"token"`
}

type CacheConfig struct {
	ImagesDir  string  `json:"images_dir"`
}

type MonitorConfig struct {
	Enable  bool  `json:"enable"`
	BaseUrl  string  `json:"base_url"`
	Token  string  `json:"token"`
}

type DataSourceConfig struct {
	Monitor  MonitorConfig  `json:"monitor"`
}

type GlobalConfig struct {
	Http  *HttpConfig  `json:"http"`
	Cache  CacheConfig  `json:"cache"`
	DataSource  DataSourceConfig  `json:"data_source"`
}

var (
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func InitConfig(cfg string) error {
	if cfg == "" {
		log.Println("use -c to specify configuration file")
	}
	_, err := os.Stat(cfg)
	if os.IsExist(err) {
		log.Println("config file not found")
		return err
	}
	b,err := ioutil.ReadFile(cfg)
	if err != nil {
		log.Printf("read file %s error %v \n", cfg, err)
		return err
	}
	configContent := strings.TrimSpace(string(b))
	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Println("parse config file:", cfg, "fail:", err)
		return err
	}
	lock.Lock()
	defer lock.Unlock()
	config = &c
	log.Println("read config file:", cfg, "successfully")
	return nil
}