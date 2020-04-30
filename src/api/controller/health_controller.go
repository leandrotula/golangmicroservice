package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Up(c *gin.Context) {

	c.JSON(http.StatusOK, "ok")

}
