package service

import (
	"github.com/leandrotula/golangmicroservice/src/api/domain/github"
	"github.com/leandrotula/golangmicroservice/src/api/errorApi"
	"github.com/leandrotula/golangmicroservice/src/api/repository"
	"github.com/leandrotula/golangmicroservice/src/provider/github_provider"
	"net/http"
	"strings"
)

type createRepoInterface interface {

	CreateRepo(request *repository.ApiRequest) (*repository.ApiResponse, errorApi.ApiError)
}

type createRepoImpl struct {}

var (
	CreateRepoOperation createRepoInterface
)

func init() {
	CreateRepoOperation = &createRepoImpl{}
}

func (op *createRepoImpl) CreateRepo(request *repository.ApiRequest) (*repository.ApiResponse, errorApi.ApiError) {

	inputName := strings.TrimSpace(request.Name)
	if inputName == "" {
		return nil, errorApi.NewBadRequestError("invalid input name")
	}

	req := github.CreateRepositoryRequestGithub{Name: inputName, Description: request.Description}

	response, errorResponse, genericError := github_provider.CreatePostRepository("", req)

	if errorResponse != nil {
		return nil, errorApi.NewApiError(errorResponse.Message, errorResponse.StatusCode)
	}

	if genericError != nil {

		return nil, errorApi.NewApiError(genericError.Message, http.StatusBadRequest)
	}

	return &repository.ApiResponse{
		ID:       response.ID,
		Name:     response.Name,
		FullName: response.FullName,
	}, nil
}