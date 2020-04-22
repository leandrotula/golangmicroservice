package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/leandrotula/golangmicroservice/services"
	"github.com/leandrotula/golangmicroservice/util"
	"net/http"
	"strconv"
)

func GetUser(c *gin.Context) {

	id := c.Param("id")

	userId, parserError := strconv.ParseInt(id, 10, 64)

	if parserError != nil {

		responseError := util.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "Could not convert to desired id type",
		}

		c.JSON(responseError.Code, responseError)
		return
	}

	user, err := services.GetUser(userId)

	if err != nil {

		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusOK, user)

}
