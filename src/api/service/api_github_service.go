package service

import (
	"github.com/leandrotula/golangmicroservice/src/api/domain/github"
	"github.com/leandrotula/golangmicroservice/src/api/errorApi"
	"github.com/leandrotula/golangmicroservice/src/api/provider/github_provider"
	"github.com/leandrotula/golangmicroservice/src/api/repository"
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

	//authorizationHeader := environment.RetrieveAuthorizationHeader()
	authorizationHeader := "6a544c7a2ca60e6d35375a5824a121114d7909b4"
	response, errorResponse, genericError := github_provider.CreatePostRepository(authorizationHeader, req)

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