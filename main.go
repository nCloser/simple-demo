package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nCloser/simple-demo/controller"
	"github.com/nCloser/simple-demo/service"
)

func main() {

	controller.ConnectDB()

	go service.RunMessageServer()

	r := gin.Default()

	initRouter(r)

	r.Run(":10110") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080"), change the port from 8080 to 10110
}
