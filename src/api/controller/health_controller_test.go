package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/leandrotula/golangmicroservice/src/api/errorApi"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUp(t *testing.T) {

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	c.Request = request

	Up(c)

	assert.NotNil(t, response)
	assert.EqualValues(t, http.StatusOK, response.Code)
	messages, _ := errorApi.DeserializeByteResponse(response.Body.Bytes())
	assert.EqualValues(t, "ok", messages.ApiMessage)

}
