package app

import (
	"github.com/gin-gonic/gin"
)

var ginHttp = gin.Default()

func StartApp() {

	mapUrls()

	if err := ginHttp.Run(":8081"); err != nil {
		panic(err)
	}

}
