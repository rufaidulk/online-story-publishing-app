package helper

type SuccessResponse struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

type ErrorResponse struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

func NewSuccessResponse(status int, msg string, data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Status: status,
		Msg:    msg,
		Data:   data,
	}
}

func NewErrorResponse(status int, msg string) *ErrorResponse {
	return &ErrorResponse{
		Status: status,
		Msg:    msg,
	}
}
