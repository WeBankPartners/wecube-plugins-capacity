package services

import (
	"fmt"
	"log"
	"time"
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/core"
	_ "github.com/go-sql-driver/mysql"
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

func SaveRConfig(param models.RWorkTable) error {
	var actions []*Action
	actions = append(actions, &Action{Sql:"delete from r_work where guid=?", Param:[]interface{}{param.Guid}})
	actions = append(actions, &Action{Sql:"insert into r_work value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,NOW())", Param:[]interface{}{param.Guid,param.Name,param.Workspace,param.EndpointA,param.EndpointB,param.MetricA,param.MetricB,param.TimeSelect,param.LegendX,param.LegendY,param.Output,param.Expr,param.FuncA,param.FuncB,param.Level}})
	return Transaction(actions)
}

func ListRConfig(guid string) (err error, result []*models.RWorkTable) {
	if guid == "" {
		err = mysqlEngine.SQL("SELECT * FROM r_work").Find(&result)
	}else{
		err = mysqlEngine.SQL("SELECT * FROM r_work WHERE guid=?", guid).Find(&result)
	}
	if len(result) == 0 {
		result = []*models.RWorkTable{}
	}
	return err,result
}

func DeleteRConfig(guid string) error {
	_,err := mysqlEngine.Exec("DELETE FROM r_work WHERE guid=?", guid)
	return err
}