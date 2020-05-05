package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/leandrotula/golangmicroservice/src/api/errorApi"
	"github.com/leandrotula/golangmicroservice/src/api/repository"
	"github.com/leandrotula/golangmicroservice/src/api/service"
	"net/http"
)

func CreateRepo(c *gin.Context) {

	var request repository.ApiRequest
	if bindError := c.ShouldBindBodyWith(&request, binding.JSON); bindError != nil {

		errors := errorApi.NewBadRequestError("invalid json body")
		c.JSON(errors.Status(), errors)

		return

	}

	response, err := service.CreateRepoOperation.CreateRepo(&request)

	if err != nil {

		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func CreateRepos(c *gin.Context) {

	var request []repository.ApiRequest
	if bindError := c.ShouldBindBodyWith(&request, binding.JSON); bindError != nil {

		errors := errorApi.NewBadRequestError("invalid json body")
		c.JSON(errors.Status(), errors)

		return

	}

	response, err := service.CreateRepoOperation.CreateRepos(request)

	if err != nil {

		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusCreated, response)
}