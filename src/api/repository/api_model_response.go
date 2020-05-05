package repository

import "github.com/leandrotula/golangmicroservice/src/api/errorApi"

type ApiResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}

type CreateReposResponse struct {
	StatusCode int `json:"status_code"`
	Results []CreateRepositoriesResponse `json:"results"`
}

type CreateRepositoriesResponse struct {

	Response *ApiResponse `json:"response"`
	Error errorApi.ApiError `json:"error"`
}