package services

import (
	"fmt"
	"log"
	"time"
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/core"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/models"
)

var mysqlEngine *xorm.Engine

func InitDbEngine() (err error) {
	connectStr := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?collation=utf8mb4_unicode_ci&allowNativePasswords=true",
		models.Config().Mysql.User, models.Config().Mysql.Password, "tcp", models.Config().Mysql.Server, models.Config().Mysql.Port, models.Config().Mysql.DataBase)
	mysqlEngine,err = xorm.NewEngine("mysql", connectStr)
	if err != nil {
		log.Printf("init mysql fail with connect: %s error: %v \n", connectStr, err)
	}else{
		mysqlEngine.SetMaxIdleConns(models.Config().Mysql.MaxIdle)
		mysqlEngine.SetMaxOpenConns(models.Config().Mysql.MaxOpen)
		mysqlEngine.SetConnMaxLifetime(time.Duration(models.Config().Mysql.Timeout)*time.Second)
		mysqlEngine.Charset("utf8")
		// 使用驼峰式映射
		mysqlEngine.SetMapper(core.SnakeMapper{})
		log.Println("init mysql success ")
	}
	return err
}

type Action struct {
	Sql  string
	Param  []interface{}
}

func Transaction(actions []*Action) error {
	if len(actions) == 0 {
		return fmt.Errorf("transaction actions is null")
	}
	session := mysqlEngine.NewSession()
	err := session.Begin()
	for _,action := range actions {
		params := make([]interface{}, 0)
		params = append(params, action.Sql)
		for _,v := range action.Param {
			params = append(params, v)
		}
		_,err = session.Exec(params...)
		if err != nil {
			session.Rollback()
			break
		}
	}
	if err==nil {
		err = session.Commit()
	}
	session.Close()
	return err
}