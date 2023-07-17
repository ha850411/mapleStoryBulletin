package routes

import (
	"go_colly/controllers"

	"github.com/gin-gonic/gin"
)

func SetRouter(r *gin.Engine) {
	r.GET("/bulletin", controllers.Test)
}
