package models

type MonitorOptionResponse struct {
	Code  int  `json:"code"`
	Message  string  `json:"message"`
	Data  []OptionModel  `json:"data"`
}

type OptionModel struct {
	Id  int  `json:"id"`
	OptionValue  string  `json:"option_value"`
	OptionText   string  `json:"option_text"`
	Active  bool  `json:"active"`
	OptionType  string  `json:"type"`
	OptionTypeName  string  `json:"option_type_name"`
}

type MonitorMetricResponse struct {
	Code  int  `json:"code"`
	Message  string  `json:"message"`
	Data  []MetricOptionModel  `json:"data"`
}

type MetricOptionModel struct {
	Id  int  `json:"id"`
	Metric  string  `json:"metric"`
	PromQl  string  `json:"prom_ql"`
}

type RespJson struct {
	Code  int  `json:"code"`
	Msg   string    `json:"message"`
	Data  interface{}  `json:"data"`
}

type ChartConfigObj struct {
	Endpoint   string    `form:"endpoint" json:"endpoint"`
	Metric   string    `form:"metric" json:"metric"`
	Start  string  `form:"start" json:"start"`
	End  string  `form:"end" json:"end"`
	Aggregate  string  `form:"aggregate" json:"agg"`
}

type YaxisModel struct {
	Unit  string  `json:"unit"`
}

type TimeSerialModel struct {
	Type  string  `json:"type"`
	Name  string  `json:"name"`
	Data  [][]float64  `json:"data"`
}

type DataSerialModel struct {
	Type  string  `json:"type"`
	Name  string  `json:"name"`
	Data  []float64  `json:"data"`
}

type AxisModel struct {
	Data  []float64  `json:"data"`
}

type MonitorChartResponse struct {
	Code  int  `json:"code"`
	Message  string  `json:"message"`
	Data  EChartOption  `json:"data"`
}

type EChartOption struct {
	Id     int     `json:"id"`
	Title  string  `json:"title"`
	Legend  []string  `json:"legend"`
	Xaxis  AxisModel  `json:"xaxis"`
	Yaxis  YaxisModel  `json:"yaxis"`
	Series []*TimeSerialModel  `json:"series"`
	DataSeries  []*DataSerialModel  `json:"data_series"`
	IsDataSeries  bool  `json:"is_data_series"`
}