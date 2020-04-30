package github_provider

import (
	"github.com/leandrotula/golangmicroservice/src/api/client"
	"github.com/leandrotula/golangmicroservice/src/api/domain/github"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

//Before All implementation
func TestMain(m *testing.M) {
	client.StartMockup()
	os.Exit(m.Run())
}

func TestCreateInvalidResponse(t *testing.T) {

	//implements invalid closer interface
	client.RestoreMockup()
	invalidCloser, _ := os.Open("a2s")
	client.AddMockBehavior(client.Mock{
		HttpMethod: http.MethodPost,
		Url: "https://api.github.com/user/repos",
		Response:   &http.Response{
			Body: invalidCloser,
			StatusCode: http.StatusInternalServerError,
		},
		Err:        nil,
	})

	response, err, invalidResponse := CreatePostRepository("", github.CreateRepositoryRequestGithub{})
	assert.Nil(t, response)
	assert.Nil(t, invalidResponse)
	assert.NotNil(t, err)
	assert.EqualValues(t, "internal server errorApi", err.Message)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)

}

func TestCreateAndProcessValidResponse(t *testing.T) {

	client.RestoreMockup()
	client.AddMockBehavior(client.Mock{
		HttpMethod: http.MethodPost,
		Url: "https://api.github.com/user/repos",
		Response:   &http.Response{
			Body: ioutil.NopCloser(strings.NewReader("{\"id\":1296269,\"node_id\":\"MDEwOlJlcG9zaXRvcnkxMjk2MjY5\",\"name\":\"Hello-World\",\"full_name\":\"octocat/Hello-World\",\"owner\":{\"login\":\"octocat\",\"id\":1,\"node_id\":\"MDQ6VXNlcjE=\",\"avatar_url\":\"https://github.com/images/errorApi/octocat_happy.gif\",\"gravatar_id\":\"\",\"url\":\"https://api.github.com/users/octocat\",\"html_url\":\"https://github.com/octocat\",\"followers_url\":\"https://api.github.com/users/octocat/followers\",\"following_url\":\"https://api.github.com/users/octocat/following{/other_user}\",\"gists_url\":\"https://api.github.com/users/octocat/gists{/gist_id}\",\"starred_url\":\"https://api.github.com/users/octocat/starred{/owner}{/repo}\",\"subscriptions_url\":\"https://api.github.com/users/octocat/subscriptions\",\"organizations_url\":\"https://api.github.com/users/octocat/orgs\",\"repos_url\":\"https://api.github.com/users/octocat/repos\",\"events_url\":\"https://api.github.com/users/octocat/events{/privacy}\",\"received_events_url\":\"https://api.github.com/users/octocat/received_events\",\"type\":\"User\",\"site_admin\":false}}")),
			StatusCode: http.StatusOK,
		},
		Err:        nil,
	})

	response, err, invalidResponse := CreatePostRepository("", github.CreateRepositoryRequestGithub{})
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.Nil(t, invalidResponse)
	assert.EqualValues(t, 1296269, response.ID)
	assert.EqualValues(t, "MDEwOlJlcG9zaXRvcnkxMjk2MjY5", response.NodeID)
	assert.EqualValues(t, "Hello-World", response.Name)
	assert.EqualValues(t, "https://github.com/images/errorApi/octocat_happy.gif", response.Owner.AvatarURL)

}

//This test should fail because we are receiving and id of type string instead of an id of type int
func TestCreateProcessNonValidSuccessResponse(t *testing.T) {

	client.RestoreMockup()
	client.AddMockBehavior(client.Mock{
		HttpMethod: http.MethodPost,
		Url: "https://api.github.com/user/repos",
		Response:   &http.Response{
			Body: ioutil.NopCloser(strings.NewReader("{\"id\":\"1296269\",\"node_id\":\"MDEwOlJlcG9zaXRvcnkxMjk2MjY5\",\"name\":\"Hello-World\",\"full_name\":\"octocat/Hello-World\"}")),
			StatusCode: http.StatusOK,
		},
		Err:        nil,
	})

	response, err, invalidResponse := CreatePostRepository("", github.CreateRepositoryRequestGithub{})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Nil(t, invalidResponse)
	assert.EqualValues(t, "parsing errorMarshalling response", err.Message)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)

}

func TestCreateAndProcessUnauthorizedResponse(t *testing.T) {

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

	response, err, invalidResponse := CreatePostRepository("", github.CreateRepositoryRequestGithub{})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Nil(t, invalidResponse)
	assert.EqualValues(t, "unauthorized access", err.Message)
	assert.EqualValues(t, http.StatusUnauthorized, err.StatusCode)

}

func TestCreateUnprocessableEntityResponse(t *testing.T) {

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

	response, err, invalidResponse := CreatePostRepository("", github.CreateRepositoryRequestGithub{})
	assert.Nil(t, response)
	assert.Nil(t, err)
	assert.NotNil(t, invalidResponse)
	assert.EqualValues(t, "Repository creation failed.", invalidResponse.Message)
	assert.EqualValues(t, "name", invalidResponse.Errors[0].Field)
	assert.EqualValues(t, "custom", invalidResponse.Errors[0].Code)
	assert.EqualValues(t, "Repository", invalidResponse.Errors[0].Resource)

}

func TestCreateNonExpectedStatusCode(t *testing.T) {

	client.RestoreMockup()
	client.AddMockBehavior(client.Mock{
		HttpMethod: http.MethodPost,
		Url: "https://api.github.com/user/repos",
		Response:   &http.Response{
			Body: ioutil.NopCloser(strings.NewReader("{\"message\":\"Repository creation failed.\",\"errors\":[{\"resource\":\"Repository\",\"code\":\"custom\",\"field\":\"name\",\"message\":\"name already exists on this account\"}],\"documentation_url\":\"https://developer.github.com/v3/repos/#create\"}")),
			StatusCode: http.StatusAlreadyReported, //Non expected statusCode
		},
		Err:        nil,
	})

	response, err, invalidResponse := CreatePostRepository("", github.CreateRepositoryRequestGithub{})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Nil(t, invalidResponse)
	assert.EqualValues(t, "Got invalid status code 208", err.Message)

}