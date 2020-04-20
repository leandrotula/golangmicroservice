package app

import (
	"github.com/gin-gonic/gin"
	"github.com/leandrotula/golangmicroservice/controllers"
)

var ginHttp = gin.Default()

func StartApp() {

	ginHttp.GET("/user/:id", controllers.GetUser)

	if err := ginHttp.Run(":8081"); err != nil {
		panic(err)
	}

}
