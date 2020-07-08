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
	Port  int  `json:"port"`
	Token  string  `json:"token"`
}

type CacheConfig struct {
	WorkspaceDir  string  `json:"workspace_dir"`
	CleanInterval  int  `json:"clean_interval"`
}

type MonitorConfig struct {
	Enable  bool  `json:"enable"`
	BaseUrl  string  `json:"base_url"`
	Token  string  `json:"token"`
}

type DataSourceConfig struct {
	Monitor  MonitorConfig  `json:"monitor"`
}

type MysqlConfig struct {
	Type  string  `json:"type"`
	Server  string  `json:"server"`
	Port  int     `json:"port"`
	User  string  `json:"user"`
	Password   string  `json:"password"`
	DataBase  string  `json:"database"`
	MaxOpen  int  `json:"maxOpen"`
	MaxIdle  int  `json:"maxIdle"`
	Timeout  int  `json:"timeout"`
}

type GlobalConfig struct {
	Http  *HttpConfig  `json:"http"`
	Mysql  MysqlConfig  `json:"mysql"`
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
	config = &c
	log.Println("read config file:", cfg, "successfully")
	lock.Unlock()
	initWorkspaceDir()
	return nil
}