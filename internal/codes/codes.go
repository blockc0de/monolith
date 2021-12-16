package codes

type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type CodeErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func NewCodeError(status int, msg string) *CodeError {
	return &CodeError{Code: status, Msg: msg}
}

func (e *CodeError) Error() string {
	return e.Msg
}

func (e *CodeError) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		Code:  e.Code,
		Error: e.Msg,
	}
}
