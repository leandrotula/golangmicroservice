package errorApi

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ApiError interface {
	Status() int
	Message() string
	Error() string
}

type apiError struct {
	ApiStatus           int    `json:"api_status"`
	ApiMessage          string `json:"api_message"`
	ApiErrorDescription string `json:"api_error_description,omitempty"`
}

func (a *apiError) Status() int {
	return a.ApiStatus
}

func (a *apiError) Message() string {
	return a.ApiMessage
}

func (a *apiError) Error() string {
	return a.ApiErrorDescription
}

func NewApiErrorNotFound(message string) ApiError {
	
	return &apiError{
		ApiStatus:  http.StatusNotFound,
		ApiMessage: message,
	}
}

func NewInternalErrorFound(message string) ApiError {

	return &apiError{
		ApiStatus:  http.StatusInternalServerError,
		ApiMessage: message,
	}
}

func NewBadRequestError(message string) ApiError {

	return &apiError{
		ApiStatus:  http.StatusBadRequest,
		ApiMessage: message,
	}
}

func NewApiError(message string, code int) ApiError {

	return &apiError{
		ApiStatus:  code,
		ApiMessage: message,
	}

}

func DeserializeByteResponse(data []byte)(*apiError, error) {

	var apiError apiError

	if err := json.Unmarshal(data, &apiError); err != nil {
		return nil, errors.New("problem during conversion")
	}

	return &apiError, nil
}





