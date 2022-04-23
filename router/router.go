package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hiro942/elden-client/middleware"
)

func Routers() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())

	{

	}

	return r
}
