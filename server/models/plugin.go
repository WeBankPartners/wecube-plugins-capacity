package models

type PluginRequest struct {
	RequestId  string  `json:"requestId"`
	Inputs  []*PluginRequestInput  `json:"inputs"`
}

type PluginRequestInput struct {
	CallbackParameter  string  `json:"callbackParameter"`
	Guid  string  `json:"guid"`
	TemplateName  string  `json:"templateName"`
	Start string  `json:"start"`
	End   string  `json:"end"`
	Save  string  `json:"save"`
}

type PluginResponse struct {
	ResultCode  string  `json:"resultCode"`
	ResultMessage string  `json:"resultMessage"`
	Results  PluginResponseOutputs  `json:"results"`
}

type PluginResponseOutputs struct {
	Outputs  []*PluginResponseOutput  `json:"outputs"`
}

type PluginResponseOutput struct {
	CallbackParameter  string  `json:"callbackParameter"`
	Guid  string  `json:"guid"`
	ErrorCode  string  `json:"errorCode"`
	ErrorMessage  string  `json:"errorMessage"`
	FuncOld  string  `json:"funcOld"`
	FuncNew  string  `json:"funcNew"`
	LevelOld string  `json:"levelOld"`
	LevelNew string  `json:"levelNew"`
	TemplateName  string  `json:"templateName"`
}
