package main

import (
	"go_colly/conf"
	"go_colly/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	ginServer := gin.Default()
	routes.SetRouter(ginServer)
	ginServer.Run(conf.Settings.Server.Port)
}
