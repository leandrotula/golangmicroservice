package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/leandrotula/golangmicroservice/src/api/client"
	"github.com/leandrotula/golangmicroservice/src/api/errorApi"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {

	client.StartMockup()
	os.Exit(m.Run())
}

func TestCreateRepoWithInvalidName(t *testing.T) {

	response := httptest.NewRecorder()
	c, _:= gin.CreateTestContext(response)

	request, _:= http.NewRequest(http.MethodPost, "/repositories", strings.NewReader(`{"name":""}`))
	c.Request = request

	CreateRepo(c)

	assert.EqualValues(t, http.StatusBadRequest, response.Code)
	apiError, _ := errorApi.DeserializeByteResponse(response.Body.Bytes())
	assert.NotNil(t, apiError)
	assert.EqualValues(t, "invalid input name", apiError.Message())

}

func TestCreateRepoWithInvalidJsonBody(t *testing.T) {

	response := httptest.NewRecorder()
	c, _:= gin.CreateTestContext(response)

	request, _:= http.NewRequest(http.MethodPost, "/repositories", strings.NewReader(``))
	c.Request = request

	CreateRepo(c)

	assert.EqualValues(t, http.StatusBadRequest, response.Code)
	apiError, _ := errorApi.DeserializeByteResponse(response.Body.Bytes())
	assert.NotNil(t, apiError)
	assert.EqualValues(t, "invalid json body", apiError.Message())

}

func TestCreateRepoSuccess(t *testing.T) {

	response := httptest.NewRecorder()
	c, _:= gin.CreateTestContext(response)
	request, _:= http.NewRequest(http.MethodPost, "/repositories",
		strings.NewReader(`{"name":"repo-from-go-api","description":"test repo creation"}`))
	c.Request = request

	client.RestoreMockup()
	client.AddMockBehavior(client.Mock{
		HttpMethod: http.MethodPost,
		Url: "https://api.github.com/user/repos",
		Response:   &http.Response{
			Body: ioutil.NopCloser(strings.NewReader("{\"id\":1296269,\"node_id\":\"MDEwOlJlcG9zaXRvcnkxMjk2MjY5\",\"name\":\"Hello-World\",\"full_name\":\"octocat/Hello-World\",\"owner\":{\"login\":\"octocat\",\"id\":1,\"node_id\":\"MDQ6VXNlcjE=\",\"avatar_url\":\"https://github.com/images/errorApi/octocat_happy.gif\",\"gravatar_id\":\"\",\"url\":\"https://api.github.com/users/octocat\",\"html_url\":\"https://github.com/octocat\",\"followers_url\":\"https://api.github.com/users/octocat/followers\",\"following_url\":\"https://api.github.com/users/octocat/following{/other_user}\",\"gists_url\":\"https://api.github.com/users/octocat/gists{/gist_id}\",\"starred_url\":\"https://api.github.com/users/octocat/starred{/owner}{/repo}\",\"subscriptions_url\":\"https://api.github.com/users/octocat/subscriptions\",\"organizations_url\":\"https://api.github.com/users/octocat/orgs\",\"repos_url\":\"https://api.github.com/users/octocat/repos\",\"events_url\":\"https://api.github.com/users/octocat/events{/privacy}\",\"received_events_url\":\"https://api.github.com/users/octocat/received_events\",\"type\":\"User\",\"site_admin\":false}}")),
			StatusCode: http.StatusCreated,
		},
		Err:        nil,
	})

	CreateRepo(c)

	assert.EqualValues(t, http.StatusCreated, response.Code)

}
