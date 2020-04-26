package errorApi

import "net/http"

type ApiError interface {
	Status() int
	Message() string
	Error() string
}

type apiError struct {
	apiStatus int `json:"api_status"`
	apiMessage string `json:"api_message"`
	apiErrorDescription string `json:"api_error_description,omitempty"`
}

func (a *apiError) Status() int {
	return a.apiStatus
}

func (a *apiError) Message() string {
	return a.apiMessage
}

func (a *apiError) Error() string {
	return a.apiErrorDescription
}

func NewApiErrorNotFound(message string) ApiError {
	
	return &apiError{
		apiStatus:           http.StatusNotFound,
		apiMessage:          message,
	}
}

func NewInternalErrorFound(message string) ApiError {

	return &apiError{
		apiStatus:           http.StatusInternalServerError,
		apiMessage:          message,
	}
}

func NewBadRequestError(message string) ApiError {

	return &apiError{
		apiStatus:           http.StatusBadRequest,
		apiMessage:          message,
	}
}

func NewApiError(message string, code int) ApiError {

	return &apiError{
		apiStatus:           code,
		apiMessage:          message,
	}

}





