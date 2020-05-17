package service

import (
	"errors"
	"github.com/leandrotula/golangmicroservice/src/api/client"
	"github.com/leandrotula/golangmicroservice/src/api/repository"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
