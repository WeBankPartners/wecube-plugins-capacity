package models

import (
	"sync"
	"fmt"
	"time"
)

var getNameLock = new(sync.RWMutex)
var WorkspaceDir string

func GetWorkspaceName() (name string) {
	getNameLock.Lock()
	name = fmt.Sprintf("R_%d", time.Now().UnixNano())
	getNameLock.Unlock()
	return name
}

func initWorkspaceDir()  {
	if Config().Cache.WorkspaceDir == "" {
		WorkspaceDir = "."
	}else{
		WorkspaceDir = Config().Cache.WorkspaceDir
		if WorkspaceDir[len(WorkspaceDir)-1:] == "/" {
			WorkspaceDir = WorkspaceDir[:len(WorkspaceDir)-1]
		}
	}
}

type RunScriptResult struct {
	Guid       string  `json:"guid"`
	Workspace  string  `json:"workspace"`
	Level      int     `json:"level"`
	Output     string  `json:"output"`
	Images     []string `json:"images"`
	FuncExpr   string  `json:"func_expr"`
	FuncX      []*FuncXObj  `json:"func_x"`
	FuncB      float64  `json:"func_b"`
	Chart      EChartOption `json:"chart"`
}

type RRequestMonitor struct {
	Config  []ChartConfigObj
	LegendY   string    `json:"legend_y"`
	LegendX   []string    `json:"legend_x"`
	RemoveList  []float64  `json:"remove_list"`
}

type RRequestExcel struct {
	Enable    bool      `json:"enable"`
	LegendY   string    `json:"legend_y"`
	LegendX   []string    `json:"legend_x"`
	RemoveList  []int  `json:"remove_list"`
}

type RRequestParam struct {
	Guid      string    `json:"guid"`
	Monitor   RRequestMonitor `json:"monitor"`
	XData    [][]float64  `json:"x_data"`
	YData    []float64  `json:"y_data"`
	FuncX    []*FuncXObj  `json:"func_x"`
	FuncB      float64  `json:"func_b"`
	Excel     RRequestExcel  `json:"excel"`
	MinLevel  int  `json:"min_level"`
}

type SaveWorkParam struct{
	Guid       string  `json:"guid"`
	Name       string  `json:"name"`
	Workspace  string  `json:"workspace"`
	Level      int     `json:"level"`
	Output     string  `json:"output"`
	Images     []string `json:"images"`
	FuncExpr   string  `json:"func_expr"`
	FuncX      []*FuncXObj  `json:"func_x"`
	FuncB      float64  `json:"func_b"`
	YReal      []float64 `json:"y_real"`
	YFunc      []float64 `json:"y_func"`
	Monitor    RRequestMonitor `json:"monitor"`
}

type RWorkTable struct {
	Guid  string  `json:"guid"`
	Name  string  `json:"name"`
	Workspace  string  `json:"workspace"`
	Output     string  `json:"output"`
	Expr       string  `json:"expr"`
	FuncX      string  `json:"func_x"`
	FuncXName  string  `json:"func_x_name"`
	FuncB      string  `json:"func_b"`
	Level      int     `json:"level"`
	LegendX    string  `json:"legend_x"`
	LegendY    string  `json:"legend_y"`
	UpdateAt   time.Time  `json:"update_at"`
}

type RImagesTable struct {
	Id    int     `json:"id"`
	Guid  string  `json:"guid"`
	Workspace  string  `json:"workspace"`
	Data     []uint8  `json:"data"`
	UpdateAt   time.Time  `json:"update_at"`
}

type RChartTableInput struct {
	Guid  string  `json:"guid"`
	YReal  []float64  `json:"y_real"`
	YFunc  []float64  `json:"y_func"`
	UpdateAt   time.Time  `json:"update_at"`
}

type RChartTable struct {
	Guid  string  `json:"guid"`
	YReal  string  `json:"y_real"`
	YFunc  string  `json:"y_func"`
	UpdateAt   time.Time  `json:"update_at"`
}

type RMonitorTable struct {
	Id  int  `json:"id"`
	Guid  string  `json:"guid"`
	Endpoint  string  `json:"endpoint"`
	Metric  string  `json:"metric"`
	Agg     string  `json:"agg"`
	Start   string  `json:"start"`
	End     string  `json:"end"`
	UpdateAt   time.Time  `json:"update_at"`
}

type FuncXObj struct {
	PValue   float64  `json:"p_value"`
	Estimate float64  `json:"estimate"`
	Level    int      `json:"level"`
	Index    int      `json:"index"`
	FuncName string   `json:"func_name"`
	Legend   string   `json:"legend"`
	Data     []float64 `json:"data"`
}

type FuncXSortList []*FuncXObj

func (s FuncXSortList) Len() int {
	return len(s)
}

func (s FuncXSortList) Swap(i,j int)  {
	s[i], s[j] = s[j], s[i]
}

func (s FuncXSortList) Less(i,j int) bool {
	return s[i].PValue < s[j].PValue
}

type YXDataTable struct {
	Title  []string  `json:"title"`
	Data   []map[string]string `json:"data"`
}

type YXDataObj struct {
	Legend  []string  `json:"legend"`
	Data  [][]float64 `json:"data"`
}

type RCalcParam struct {
	Guid  string  `json:"guid"`
	AddData  [][]float64  `json:"add_data"`
}

type RCalcResult struct {
	Guid  string  `json:"guid"`
	Chart EChartOption `json:"chart"`
	Table YXDataTable `json:"table"`
}