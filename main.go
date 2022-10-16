package main

import (
	"github.com/gin-gonic/gin"
	"users-task/configs"
	"users-task/routes"
)

func main() {

	configs.Logger = configs.NewLogger("log.txt")
	router := gin.Default()

	//run database
	configs.ConnectDB()

	//routes
	routes.UserRoute(router)

	router.Run(configs.ServerPort)
}
