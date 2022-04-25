package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hiro942/elden-client/controller"
	"github.com/hiro942/elden-client/middleware"
)

func Routers() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())

	authGroup := r.Group("auth")
	{
		authGroup.POST("prehandover/location", controller.PreHandoverForLocation)
		authGroup.POST("prehandover/new_satellite", controller.PreHandoverForNewSatellite)
	}

	return r
}
