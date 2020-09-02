package models

type PluginRequest struct {
	RequestId  string  `json:"requestId"`
	Inputs  []*PluginRequestInput  `json:"inputs"`
}

type PluginRequestInput struct {
	CallbackParameter  string  `json:"callbackParameter"`
	Guid  string  `json:"guid"`
	TemplateGuid  string  `json:"template_guid"`
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
	FuncOld  string  `json:"func_old"`
	FuncNew  string  `json:"func_new"`
}
