package controllers

import (
	"encoding/json"
	"github.com/golangmicroservice/services"
	"github.com/golangmicroservice/util"
	"net/http"
	"strconv"
)

func GetUser(response http.ResponseWriter, request *http.Request) {

	id := request.URL.Query().Get("id")

	userId, parserError := strconv.ParseInt(id, 10, 64)
	response.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if parserError != nil {

		responseError := util.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "Could not convert to desired id type",
		}
		marshalledError, _ := json.Marshal(responseError)
		response.WriteHeader(responseError.Code)
		_, _ = response.Write(marshalledError)

		return
	}

	user, err := services.GetUser(userId)

	if err != nil {

		marshalledError, _ := json.Marshal(err)
		response.WriteHeader(err.Code)
		_, _ = response.Write(marshalledError)

		return
	}

	resp, _ := json.Marshal(user)
	_, _ = response.Write(resp)

}
