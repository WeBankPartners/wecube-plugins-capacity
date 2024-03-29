package models

const (
	DatetimeFormat = `2006-01-02 15:04:05`
)

type MonitorOptionResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    []OptionModel `json:"data"`
}

type OptionModel struct {
	Id             int    `json:"id"`
	OptionValue    string `json:"option_value"`
	OptionText     string `json:"option_text"`
	Active         bool   `json:"active"`
	OptionType     string `json:"type"`
	OptionTypeName string `json:"option_type_name"`
}

type MonitorMetricResponse struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Data    []MetricOptionModel `json:"data"`
}

type MetricOptionModel struct {
	Id     int    `json:"id"`
	Metric string `json:"metric"`
	PromQl string `json:"prom_ql"`
}

type RespJson struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

type ChartConfigObj struct {
	Endpoint  string `form:"endpoint" json:"endpoint"`
	Metric    string `form:"metric" json:"metric"`
	Start     string `form:"start" json:"start"`
	End       string `form:"end" json:"end"`
	Aggregate string `form:"aggregate" json:"agg"`
}

type YaxisModel struct {
	Unit string `json:"unit"`
}

type TimeSerialModel struct {
	Type string      `json:"type"`
	Name string      `json:"name"`
	Data [][]float64 `json:"data"`
}

type DataSerialModel struct {
	Type string    `json:"type"`
	Name string    `json:"name"`
	Data []float64 `json:"data"`
}

type AxisModel struct {
	Data []float64 `json:"data"`
}

type MonitorChartResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    EChartOption `json:"data"`
}

type EChartOption struct {
	Id           int                `json:"id"`
	Title        string             `json:"title"`
	Legend       []string           `json:"legend"`
	Xaxis        AxisModel          `json:"xaxis"`
	Yaxis        YaxisModel         `json:"yaxis"`
	Series       []*TimeSerialModel `json:"series"`
	DataSeries   []*DataSerialModel `json:"data_series"`
	IsDataSeries bool               `json:"is_data_series"`
}

type ExportResultParamObj struct {
	Name     string  `json:"name"`
	Metric   string  `json:"metric"`
	Estimate float64 `json:"estimate"`
}

type ExportResultObj struct {
	RWorkName  string                  `json:"r_work_name"`
	Endpoint   string                  `json:"endpoint"`
	XParams    []*ExportResultParamObj `json:"x_params"`
	YFunc      ExportResultParamObj    `json:"y_func"`
	FuncExpr   string                  `json:"func_expr"`
	UpdateTime string                  `json:"update_time"`
}

type ChartQueryParam struct {
	ChartId    int                    `json:"chart_id"`
	Title      string                 `json:"title"`
	Unit       string                 `json:"unit"`
	Start      int64                  `json:"start"`
	End        int64                  `json:"end"`
	TimeSecond int64                  `json:"time_second"`
	Aggregate  string                 `json:"aggregate"`
	AggStep    int64                  `json:"agg_step"`
	Step       int                    `json:"step"`
	Data       []*ChartQueryConfigObj `json:"data"`
}

type ChartQueryConfigObj struct {
	Endpoint     string `json:"endpoint"`
	Metric       string `json:"metric"`
	PromQl       string `json:"prom_ql"`
	AppObject    string `json:"app_object"`
	EndpointType string `json:"endpoint_type"`
}
