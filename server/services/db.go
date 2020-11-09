package services

import (
	"fmt"
	"time"
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/core"
	_ "github.com/go-sql-driver/mysql"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/models"
	"github.com/WeBankPartners/wecube-plugins-capacity/server/util/log"
	"strings"
	"strconv"
)

var mysqlEngine *xorm.Engine

func InitDbEngine() (err error) {
	connectStr := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?collation=utf8mb4_unicode_ci&allowNativePasswords=true",
		models.Config().Mysql.User, models.Config().Mysql.Password, "tcp", models.Config().Mysql.Server, models.Config().Mysql.Port, models.Config().Mysql.DataBase)
	mysqlEngine,err = xorm.NewEngine("mysql", connectStr)
	if err != nil {
		log.Logger.Error("Init mysql fail", log.String("connectStr",connectStr), log.Error(err))
	}else{
		mysqlEngine.SetMaxIdleConns(models.Config().Mysql.MaxIdle)
		mysqlEngine.SetMaxOpenConns(models.Config().Mysql.MaxOpen)
		mysqlEngine.SetConnMaxLifetime(time.Duration(models.Config().Mysql.Timeout)*time.Second)
		mysqlEngine.Charset("utf8")
		// 使用驼峰式映射
		mysqlEngine.SetMapper(core.SnakeMapper{})
		log.Logger.Info("Init mysql success")
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

func saveRWorkTable(param models.RWorkTable) error {
	var actions []*Action
	actions = append(actions, &Action{Sql:"delete from r_work where guid=?", Param:[]interface{}{param.Guid}})
	actions = append(actions, &Action{Sql:"insert into r_work(guid,name,workspace,output,expr,func_x,func_x_name,func_b,level,legend_x,legend_y,update_at) value (?,?,?,?,?,?,?,?,?,?,?,NOW())", Param:[]interface{}{param.Guid,param.Name,param.Workspace,param.Output,param.Expr,param.FuncX,param.FuncXName,param.FuncB,param.Level,param.LegendX,param.LegendY}})
	return Transaction(actions)
}

func getRImagesTable(guid string) (err error, result []*models.RImagesTable) {
	err = mysqlEngine.SQL("select * from r_images where guid=?", guid).Find(&result)
	return err,result
}

func saveRImagesTable(params []*models.RImagesTable) error {
	if len(params) == 0 {
		return fmt.Errorf("save images fail,no data input ")
	}
	var actions []*Action
	actions = append(actions, &Action{Sql:"delete from r_images where guid=?", Param:[]interface{}{params[0].Guid}})
	for _,v := range params {
		actions = append(actions, &Action{Sql: "insert into r_images(guid,workspace,data) value (?,?,?)", Param: []interface{}{v.Guid, v.Workspace, v.Data}})
	}
	return Transaction(actions)
}

func getRChartTable(guid string) (err error, result []*models.RChartTableInput) {
	var data []*models.RChartTable
	err = mysqlEngine.SQL("select * from r_chart where guid=?", guid).Find(&data)
	if err != nil {
		return err,result
	}
	for _,v := range data {
		var yReal,yFunc []float64
		for _,vv := range strings.Split(v.YReal, ",") {
			tmpV,_ := strconv.ParseFloat(vv, 64)
			yReal = append(yReal, tmpV)
		}
		for _,vv := range strings.Split(v.YFunc, ",") {
			tmpV,_ := strconv.ParseFloat(vv, 64)
			yFunc = append(yFunc, tmpV)
		}
		result = append(result, &models.RChartTableInput{Guid:v.Guid,YReal:yReal,YFunc:yFunc,UpdateAt:v.UpdateAt})
	}
	return err,result
}

func saveRChartTable(param models.RChartTableInput) error {
	var actions []*Action
	actions = append(actions, &Action{Sql:"delete from r_chart where guid=?", Param:[]interface{}{param.Guid}})
	var yReal,yFunc []string
	for _,v := range param.YReal {
		yReal = append(yReal, fmt.Sprintf("%f", v))
	}
	for _,v := range param.YFunc {
		yFunc = append(yFunc, fmt.Sprintf("%f", v))
	}
	actions = append(actions, &Action{Sql:"insert into r_chart(guid,y_real,y_func) value (?,?,?)", Param: []interface{}{param.Guid, strings.Join(yReal, ","), strings.Join(yFunc, ",")}})
	return Transaction(actions)
}

func ListRConfig(guid string) (err error, result []*models.RWorkTable) {
	if guid == "" {
		err = mysqlEngine.SQL("select * from r_work").Find(&result)
	}else{
		err = mysqlEngine.SQL("select * from r_work where guid=?", guid).Find(&result)
	}
	if len(result) == 0 {
		result = []*models.RWorkTable{}
	}
	return err,result
}

func DeleteRConfig(guid string) error {
	var actions []*Action
	actions = append(actions, &Action{Sql:"delete from r_work where guid=?", Param:[]interface{}{guid}})
	actions = append(actions, &Action{Sql:"delete from r_images where guid=?", Param:[]interface{}{guid}})
	actions = append(actions, &Action{Sql:"delete from r_chart where guid=?", Param:[]interface{}{guid}})
	return Transaction(actions)
}

func saveRMonitorTable(guid string,param models.RRequestMonitor) error {
	if len(param.Config) == 0 {
		return nil
	}
	var actions []*Action
	actions = append(actions, &Action{Sql:"DELETE FROM r_monitor WHERE guid=?", Param:[]interface{}{guid}})
	insertSql := "INSERT INTO r_monitor(guid,endpoint,metric,agg,start,end) values "
	for i,v := range param.Config {
		insertSql += fmt.Sprintf("('%s','%s','%s','%s','%s','%s')", guid, v.Endpoint, v.Metric, v.Aggregate, v.Start, v.End)
		if i != len(param.Config)-1 {
			insertSql += ","
		}
	}
	actions = append(actions, &Action{Sql:insertSql})
	return Transaction(actions)
}

func getRMonitorTable(guid string) (err error,result []*models.RMonitorTable) {
	err = mysqlEngine.SQL("select * from r_monitor where guid=?", guid).Find(&result)
	return err,result
}

func getRWorkByName(name string) (err error,result []*models.RWorkTable) {
	err = mysqlEngine.SQL("select * from r_work where name=?", name).Find(&result)
	return err,result
}