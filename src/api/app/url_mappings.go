package app

import "github.com/leandrotula/golangmicroservice/src/api/controller"

func mapUrls() {

	ginHttp.GET("/health", controller.Up)
	ginHttp.POST("/repositories", controller.CreateRepo)
}
