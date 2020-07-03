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

type RRequestParam struct {
	Guid      string    `json:"guid"`
	Monitor   RRequestMonitor `json:"monitor"`
	XData    [][]float64  `json:"x_data"`
	YData    []float64  `json:"y_data"`
	FuncX    []*FuncXObj  `json:"func_x"`
	FuncB      float64  `json:"func_b"`
	AddData  []float64  `json:"add_data"`
}

type RWorkTable struct {
	Guid  string  `json:"guid"`
	Name  string  `json:"name"`
	Workspace  string  `json:"workspace"`
	EndpointA  string  `json:"endpoint_a"`
	EndpointB  string  `json:"endpoint_b"`
	MetricA    string  `json:"metric_a"`
	MetricB    string  `json:"metric_b"`
	TimeSelect string  `json:"time_select"`
	LegendX    string  `json:"legend_x"`
	LegendY    string  `json:"legend_y"`
	Output     string  `json:"output"`
	Expr       string  `json:"expr"`
	FuncA      string  `json:"func_a"`
	FuncB      string  `json:"func_b"`
	Level      int     `json:"level"`
	UpdateAt   time.Time  `json:"update_at"`
}

type FuncXObj struct {
	PValue  float64
	Estimate float64
	Level    int
	Index    int
	FuncName  string
	Data     []float64
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

type YXDataObj struct {
	Legend  []string  `json:"legend"`
	Data  [][]float64 `json:"data"`
}