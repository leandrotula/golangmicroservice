package app

import (
	"github.com/golangmicroservice/controllers"
	"net/http"
)

func StartApp() {

	http.HandleFunc("/users", controllers.GetUser)

	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}

}
