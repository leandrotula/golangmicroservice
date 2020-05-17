package service

import (
	"github.com/leandrotula/golangmicroservice/src/api/domain/github"
	"github.com/leandrotula/golangmicroservice/src/api/errorApi"
	"github.com/leandrotula/golangmicroservice/src/api/provider/environment"
	"github.com/leandrotula/golangmicroservice/src/api/provider/github_provider"
	"github.com/leandrotula/golangmicroservice/src/api/repository"
	"net/http"
	"strings"
	"sync"
)

type createRepoInterface interface {

	CreateRepo(request *repository.ApiRequest) (*repository.ApiResponse, errorApi.ApiError)
	CreateRepos(request []repository.ApiRequest) (repository.CreateReposResponse, errorApi.ApiError)
}

type createRepoImpl struct {}

var (
	CreateRepoOperation createRepoInterface
)

func init() {
	CreateRepoOperation = &createRepoImpl{}
}

func (op *createRepoImpl) CreateRepo(request *repository.ApiRequest) (*repository.ApiResponse, errorApi.ApiError) {

	inputName, apiResponse, apiError, done := validate(request)
	if done {
		return apiResponse, apiError
	}

	req := github.CreateRepositoryRequestGithub{Name: inputName, Description: request.Description}

	authorizationHeader := environment.RetrieveAuthorizationHeader()
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

func validate(request *repository.ApiRequest) (string, *repository.ApiResponse, errorApi.ApiError, bool) {
	inputName := strings.TrimSpace(request.Name)
	if inputName == "" {
		return "", nil, errorApi.NewBadRequestError("invalid input name"), true
	}
	return inputName, nil, nil, false
}

func (op *createRepoImpl) CreateRepos(requests []repository.ApiRequest) (repository.CreateReposResponse, errorApi.ApiError) {

	input := make(chan repository.CreateRepositoriesResponse)
	output := make(chan repository.CreateReposResponse)
	var wg sync.WaitGroup

	defer close(output)

	go op.handle(&wg, input, output)
	for _, r := range requests {

		wg.Add(1)
		go op.createSingleRepo(r, input)
	}

	wg.Wait()
	close(input)

	finalResult := <- output

	success := 0
	for _, tmpResult := range finalResult.Results {

		if tmpResult.Response != nil {
			success ++
		}

	}

	if success == 0 {
		finalResult.StatusCode = finalResult.Results[0].Error.Status()
	} else if success == len(requests) {
		finalResult.StatusCode = http.StatusCreated
	} else {
		finalResult.StatusCode = http.StatusPartialContent
	}
	return finalResult, nil
}

func (op *createRepoImpl) handle(wg *sync.WaitGroup, inChannel chan repository.CreateRepositoriesResponse,
	outputChannel chan repository.CreateReposResponse) {

	var result repository.CreateReposResponse

	for event := range inChannel {
		eventResult := repository.CreateRepositoriesResponse{
			Response: event.Response,
			Error:    event.Error,
		}
		result.Results = append(result.Results, eventResult)
		wg.Done()
	}

	outputChannel <- result
}

func (op *createRepoImpl) createSingleRepo(providedRequest repository.ApiRequest, output chan repository.CreateRepositoriesResponse) {

	_, _, apiError, done := validate(&providedRequest)
	if done {
		output <- repository.CreateRepositoriesResponse{
			Response: nil,
			Error:    apiError,
		}

		return
	}

	req := github.CreateRepositoryRequestGithub{Name: providedRequest.Name,
		Description: providedRequest.Description}

	authorizationHeader := environment.RetrieveAuthorizationHeader()
	response, errorResponse, genericError := github_provider.CreatePostRepository(authorizationHeader, req)

	if errorResponse != nil {
		output <- repository.CreateRepositoriesResponse{
			Response: nil,
			Error:    errorApi.NewApiError(errorResponse.Message, errorResponse.StatusCode),
		}
		return
	}

	if genericError != nil {

		output <- repository.CreateRepositoriesResponse{
			Response: nil,
			Error:    errorApi.NewApiError(genericError.Message, http.StatusBadRequest),
		}
		return

	}

	finalResponse := repository.ApiResponse{
		ID:       response.ID,
		Name:     response.Name,
		FullName: response.FullName,
	}

	output <- repository.CreateRepositoriesResponse{
		Response: &finalResponse,
		Error:    nil,
	}


}
