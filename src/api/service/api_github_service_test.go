package service

import (
	"errors"
	"github.com/leandrotula/golangmicroservice/src/api/client"
	"github.com/leandrotula/golangmicroservice/src/api/errorApi"
	"github.com/leandrotula/golangmicroservice/src/api/repository"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestMain(m *testing.M) {

	client.StartMockup()
	os.Exit(m.Run())
}

func TestCreateRepoInvalidInputName(t *testing.T) {

	request := &repository.ApiRequest{
		Name:        "",
		Description: "",
	}

	response, err := CreateRepoOperation.CreateRepo(request)

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "invalid input name", err.Message())
}

func TestErrorCreateRepoDueInvalidResponse(t *testing.T) {

	client.RestoreMockup()
	client.AddMockBehavior(client.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response:   nil,
		Err:        errors.New("invalid response"),
	})

	request := &repository.ApiRequest{
		Name:        "test-repo",
		Description: "this is a test repo creation",
	}

	response, err := CreateRepoOperation.CreateRepo(request)

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "invalid response", err.Message())

}

func TestErrorCreateRepoDueGenericError(t *testing.T) {

	client.RestoreMockup()
	client.AddMockBehavior(client.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response:   &http.Response{
			Body: ioutil.NopCloser(strings.NewReader("{\"message\":\"Repository creation failed.\",\"errors\":[{\"resource\":\"Repository\",\"code\":\"custom\",\"field\":\"name\",\"message\":\"name already exists on this account\"}],\"documentation_url\":\"https://developer.github.com/v3/repos/#create\"}")),
			StatusCode: http.StatusUnprocessableEntity,
		},
		Err:        nil,
	})

	request := &repository.ApiRequest{
		Name:        "test-repo",
		Description: "this is a test repo creation",
	}

	response, err := CreateRepoOperation.CreateRepo(request)

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Repository creation failed.", err.Message())
	assert.EqualValues(t, http.StatusBadRequest, err.Status())

}

func TestErrorCreateRepoOk(t *testing.T) {

	client.RestoreMockup()
	client.AddMockBehavior(client.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response:   &http.Response{
			Body: ioutil.NopCloser(strings.NewReader("{\"id\":1296269,\"node_id\":\"MDEwOlJlcG9zaXRvcnkxMjk2MjY5\",\"name\":\"Hello-World\",\"full_name\":\"octocat/Hello-World\",\"owner\":{\"login\":\"octocat\",\"id\":1,\"node_id\":\"MDQ6VXNlcjE=\",\"avatar_url\":\"https://github.com/images/errorApi/octocat_happy.gif\",\"gravatar_id\":\"\",\"url\":\"https://api.github.com/users/octocat\",\"html_url\":\"https://github.com/octocat\",\"followers_url\":\"https://api.github.com/users/octocat/followers\",\"following_url\":\"https://api.github.com/users/octocat/following{/other_user}\",\"gists_url\":\"https://api.github.com/users/octocat/gists{/gist_id}\",\"starred_url\":\"https://api.github.com/users/octocat/starred{/owner}{/repo}\",\"subscriptions_url\":\"https://api.github.com/users/octocat/subscriptions\",\"organizations_url\":\"https://api.github.com/users/octocat/orgs\",\"repos_url\":\"https://api.github.com/users/octocat/repos\",\"events_url\":\"https://api.github.com/users/octocat/events{/privacy}\",\"received_events_url\":\"https://api.github.com/users/octocat/received_events\",\"type\":\"User\",\"site_admin\":false}}")),
			StatusCode: http.StatusCreated,
		},
		Err:        nil,
	})

	request := &repository.ApiRequest{
		Name:        "test-repo",
		Description: "this is a test repo creation",
	}

	response, err := CreateRepoOperation.CreateRepo(request)

	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, 1296269, response.ID)

}

func TestCreateSingleRepoInvalidRequest(t *testing.T) {

	request := repository.ApiRequest{
		Name:        "",
		Description: "",
	}

	output := make(chan repository.CreateRepositoriesResponse)
	service := createRepoImpl{}

	go service.createSingleRepo(request, output)

	result := <- output
	assert.NotNil(t, result)
	assert.Nil(t, result.Response)
	err := result.Error
	assert.EqualValues(t, http.StatusBadRequest, err.Status())

}

func TestCreateSingleRepoInvalidGithubResponse(t *testing.T) {

	client.RestoreMockup()
	client.AddMockBehavior(client.Mock{
		HttpMethod: http.MethodPost,
		Url: "https://api.github.com/user/repos",
		Response:   &http.Response{
			Body: ioutil.NopCloser(strings.NewReader("{\"message\":\"Requires authentication\",\"documentation_url\":\"https://developer.github.com/v3/repos/#create\"}")),
			StatusCode: http.StatusUnauthorized,
		},
		Err:        nil,
	})

	request := repository.ApiRequest{
		Name:        "test_name",
		Description: "test_repos",
	}

	output := make(chan repository.CreateRepositoriesResponse)
	service := createRepoImpl{}

	go service.createSingleRepo(request, output)

	result := <- output
	assert.NotNil(t, result)
	assert.Nil(t, result.Response)
	err := result.Error
	assert.EqualValues(t, http.StatusUnauthorized, err.Status())

}

func TestCreateSingleRepoNotProcessableEntity(t *testing.T) {

	client.RestoreMockup()
	client.AddMockBehavior(client.Mock{
		HttpMethod: http.MethodPost,
		Url: "https://api.github.com/user/repos",
		Response:   &http.Response{
			Body: ioutil.NopCloser(strings.NewReader("{\"message\":\"Repository creation failed.\",\"errors\":[{\"resource\":\"Repository\",\"code\":\"custom\",\"field\":\"name\",\"message\":\"name already exists on this account\"}],\"documentation_url\":\"https://developer.github.com/v3/repos/#create\"}")),
			StatusCode: http.StatusUnprocessableEntity,
		},
		Err:        nil,
	})

	request := repository.ApiRequest{
		Name:        "test_name",
		Description: "test_repos",
	}

	output := make(chan repository.CreateRepositoriesResponse)
	service := createRepoImpl{}

	go service.createSingleRepo(request, output)

	result := <- output
	assert.NotNil(t, result)
	assert.Nil(t, result.Response)
	err := result.Error
	assert.EqualValues(t, http.StatusBadRequest, err.Status())

}

func TestCreateSingleRepoOk(t *testing.T) {

	client.RestoreMockup()
	client.AddMockBehavior(client.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response:   &http.Response{
			Body: ioutil.NopCloser(strings.NewReader("{\"id\":1296269,\"node_id\":\"MDEwOlJlcG9zaXRvcnkxMjk2MjY5\",\"name\":\"Hello-World\",\"full_name\":\"octocat/Hello-World\",\"owner\":{\"login\":\"octocat\",\"id\":1,\"node_id\":\"MDQ6VXNlcjE=\",\"avatar_url\":\"https://github.com/images/errorApi/octocat_happy.gif\",\"gravatar_id\":\"\",\"url\":\"https://api.github.com/users/octocat\",\"html_url\":\"https://github.com/octocat\",\"followers_url\":\"https://api.github.com/users/octocat/followers\",\"following_url\":\"https://api.github.com/users/octocat/following{/other_user}\",\"gists_url\":\"https://api.github.com/users/octocat/gists{/gist_id}\",\"starred_url\":\"https://api.github.com/users/octocat/starred{/owner}{/repo}\",\"subscriptions_url\":\"https://api.github.com/users/octocat/subscriptions\",\"organizations_url\":\"https://api.github.com/users/octocat/orgs\",\"repos_url\":\"https://api.github.com/users/octocat/repos\",\"events_url\":\"https://api.github.com/users/octocat/events{/privacy}\",\"received_events_url\":\"https://api.github.com/users/octocat/received_events\",\"type\":\"User\",\"site_admin\":false}}")),
			StatusCode: http.StatusCreated,
		},
		Err:        nil,
	})

	request := repository.ApiRequest{
		Name:        "test_name",
		Description: "test_repos",
	}

	output := make(chan repository.CreateRepositoriesResponse)
	service := createRepoImpl{}

	go service.createSingleRepo(request, output)

	result := <- output
	assert.NotNil(t, result)
	assert.Nil(t, result.Error)
	assert.NotNil(t, result.Response)
	assert.EqualValues(t, "Hello-World", result.Response.Name)
	assert.EqualValues(t, 1296269, result.Response.ID)
	assert.EqualValues(t, "octocat/Hello-World", result.Response.FullName)

}

func TestHandleConcurrentResponseWithError(t *testing.T) {

	var wg sync.WaitGroup
	input := make(chan repository.CreateRepositoriesResponse)
	output := make(chan repository.CreateReposResponse)

	service := createRepoImpl{}
	wg.Add(1)

	go func() {
		input <- repository.CreateRepositoriesResponse{
			Response: nil,
			Error: errorApi.NewInternalErrorFound("Internal server error"),
		}
		close(input)
	}()


	go service.handle(&wg, input, output)
	wg.Wait()

	result := <- output
	assert.NotNil(t, result)
	assert.Nil(t, result.Results[0].Response)
	assert.NotNil(t, result.Results[0].Error)

}

func TestHandleConcurrentSuccessfulResponse(t *testing.T) {

	var wg sync.WaitGroup
	input := make(chan repository.CreateRepositoriesResponse)
	output := make(chan repository.CreateReposResponse)

	service := createRepoImpl{}
	wg.Add(1)

	go func() {
		input <- repository.CreateRepositoriesResponse{
			Response: &repository.ApiResponse{
				ID:       123,
				Name:     "repo created",
				FullName: "user repo created",
			},
			Error: nil,
		}
		close(input)
	}()


	go service.handle(&wg, input, output)
	wg.Wait()

	result := <- output
	assert.NotNil(t, result)
	assert.NotNil(t, result.Results[0].Response)
	assert.Nil(t, result.Results[0].Error)

}

func TestCreateReposStatusCreated(t *testing.T) {

	client.RestoreMockup()
	client.AddMockBehavior(client.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response:   &http.Response{
			Body: ioutil.NopCloser(strings.NewReader("{\"id\":1296269,\"node_id\":\"MDEwOlJlcG9zaXRvcnkxMjk2MjY5\",\"name\":\"Hello-World\",\"full_name\":\"octocat/Hello-World\",\"owner\":{\"login\":\"octocat\",\"id\":1,\"node_id\":\"MDQ6VXNlcjE=\",\"avatar_url\":\"https://github.com/images/errorApi/octocat_happy.gif\",\"gravatar_id\":\"\",\"url\":\"https://api.github.com/users/octocat\",\"html_url\":\"https://github.com/octocat\",\"followers_url\":\"https://api.github.com/users/octocat/followers\",\"following_url\":\"https://api.github.com/users/octocat/following{/other_user}\",\"gists_url\":\"https://api.github.com/users/octocat/gists{/gist_id}\",\"starred_url\":\"https://api.github.com/users/octocat/starred{/owner}{/repo}\",\"subscriptions_url\":\"https://api.github.com/users/octocat/subscriptions\",\"organizations_url\":\"https://api.github.com/users/octocat/orgs\",\"repos_url\":\"https://api.github.com/users/octocat/repos\",\"events_url\":\"https://api.github.com/users/octocat/events{/privacy}\",\"received_events_url\":\"https://api.github.com/users/octocat/received_events\",\"type\":\"User\",\"site_admin\":false}}")),
			StatusCode: http.StatusCreated,
		},
		Err:        nil,
	})

	requests := []repository.ApiRequest{
		{
			Name:        "test-repo",
			Description: "test first repo",
		},
	}

	response, err := CreateRepoOperation.CreateRepos(requests)
	assert.Nil(t, err)

	assert.NotNil(t, response)
	assert.EqualValues(t, response.StatusCode, http.StatusCreated)
	assert.NotNil(t, response.Results[0].Response)
	assert.Nil(t, response.Results[0].Error)

}

func TestCreateReposPartialContent(t *testing.T) {

	client.RestoreMockup()
	client.AddMockBehavior(client.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response:   &http.Response{
			Body: ioutil.NopCloser(strings.NewReader("{\"id\":1296269,\"node_id\":\"MDEwOlJlcG9zaXRvcnkxMjk2MjY5\",\"name\":\"Hello-World\",\"full_name\":\"octocat/Hello-World\",\"owner\":{\"login\":\"octocat\",\"id\":1,\"node_id\":\"MDQ6VXNlcjE=\",\"avatar_url\":\"https://github.com/images/errorApi/octocat_happy.gif\",\"gravatar_id\":\"\",\"url\":\"https://api.github.com/users/octocat\",\"html_url\":\"https://github.com/octocat\",\"followers_url\":\"https://api.github.com/users/octocat/followers\",\"following_url\":\"https://api.github.com/users/octocat/following{/other_user}\",\"gists_url\":\"https://api.github.com/users/octocat/gists{/gist_id}\",\"starred_url\":\"https://api.github.com/users/octocat/starred{/owner}{/repo}\",\"subscriptions_url\":\"https://api.github.com/users/octocat/subscriptions\",\"organizations_url\":\"https://api.github.com/users/octocat/orgs\",\"repos_url\":\"https://api.github.com/users/octocat/repos\",\"events_url\":\"https://api.github.com/users/octocat/events{/privacy}\",\"received_events_url\":\"https://api.github.com/users/octocat/received_events\",\"type\":\"User\",\"site_admin\":false}}")),
			StatusCode: http.StatusCreated,
		},
		Err:        nil,
	})

	requests := []repository.ApiRequest{
		{
			Name:        "test-repo",
		},
		{
			Name:        "",
		},
	}

	response, err := CreateRepoOperation.CreateRepos(requests)
	assert.Nil(t, err)

	for _, tmp := range response.Results {

		if tmp.Error != nil {
			assert.EqualValues(t, http.StatusBadRequest ,tmp.Error.Status())
			assert.EqualValues(t, "invalid input name", tmp.Error.Message())
		} else {
			assert.NotNil(t, tmp.Response)
			assert.EqualValues(t, "Hello-World", tmp.Response.Name)
			assert.EqualValues(t, 1296269, tmp.Response.ID)
			assert.EqualValues(t, "octocat/Hello-World", tmp.Response.FullName)

		}

	}

	assert.NotNil(t, response)

}

func TestCreateReposWithInvalidReqBody(t *testing.T) {

	requests := []repository.ApiRequest{
		{
			Name:        "",
			Description: "",
		},
	}

	response, _ := CreateRepoOperation.CreateRepos(requests)
	assert.NotNil(t, response.Results[0].Error)
	assert.Nil(t, response.Results[0].Response)
	assert.EqualValues(t, response.StatusCode, http.StatusBadRequest)

}

func TestCreateReposWithError(t *testing.T) {

	client.RestoreMockup()
	client.AddMockBehavior(client.Mock{
		HttpMethod: http.MethodPost,
		Url: "https://api.github.com/user/repos",
		Response:   &http.Response{
			Body: ioutil.NopCloser(strings.NewReader("{\"message\":\"Requires authentication\",\"documentation_url\":\"https://developer.github.com/v3/repos/#create\"}")),
			StatusCode: http.StatusUnauthorized,
		},
		Err:        nil,
	})

	requests := []repository.ApiRequest{
		{
			Name:        "test-repo",
			Description: "test first repo",
		},
	}

	response, err := CreateRepoOperation.CreateRepos(requests)
	assert.NotNil(t, response)

	assert.Nil(t, err)
	assert.EqualValues(t, response.StatusCode, http.StatusUnauthorized)
	assert.Nil(t, response.Results[0].Response)
	assert.NotNil(t, response.Results[0].Error)

}
