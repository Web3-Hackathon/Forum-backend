package models

// BaseResponseModel is the base of every response
// it helps the client know if an operation was
// successful or if it failed
type BaseResponseModel struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// ErrorModel is used to send the occurred errors
// back to the client
type ErrorModel struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorResponseModel is returned when an error
// occurs on the webserver
type ErrorResponseModel struct {
	BaseResponseModel
	Error ErrorModel `json:"error"`
}
