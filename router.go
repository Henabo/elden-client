package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hiro942/elden-client/model/response"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	quit chan os.Signal
)

func (auth *AuthenticationService) RunRouter() {
	router := auth.GetRouter()

	server := http.Server{
		Addr:    ":" + config.HttpServerPort,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server listen err: %s", err)
		}
	}()

	quit = make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	auth.Session.Disconnect(false)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown err: %s", err)
	}

	log.Println("server exit.")
}

func (auth *AuthenticationService) GetRouter() *gin.Engine {
	r := gin.Default()
	r.Use(auth.cors())

	authGroup := r.Group("auth")
	{
		authGroup.POST("prehandover/location", auth.ReturnLocation)
		authGroup.POST("prehandover/new_satellite", auth.ReceiveNewSatellite)
	}

	return r
}

func (auth *AuthenticationService) cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		origin := c.Request.Header.Get("Origin")

		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-ID")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

// ReturnLocation 向卫星提供位置区间
// [POST] /auth/prehandover/location
func (auth *AuthenticationService) ReturnLocation(c *gin.Context) {
	// todo 估算自己的位置区间返回给卫星
	// todo 这里直接先返回一个成功响应代表成功
	response.OK(c)
}

// @Summary authentication for normal(fast) access phrase

// ReceiveNewSatellite 接收新卫星id
// [POST] /auth/prehandover/new_satellite?id=？
func (auth *AuthenticationService) ReceiveNewSatellite(c *gin.Context) {
	newSid := c.Query("id")

	var (
		DefaultSuccessMessage = fmt.Sprintf("收到交接卫星ID `%s`，并成功切换", newSid)
		DefaultErrorMessage   = "收到交接卫星ID失败"
	)

	if strings.TrimSpace(newSid) == "" {
		response.FailWithDescription(DefaultErrorMessage, "缺少参数：新卫星ID", c)
		return
	}

	if err := auth.HandoverAccess(newSid); err != nil {
		DefaultErrorMessage = DefaultSuccessMessage
		response.FailWithDescription(DefaultErrorMessage, err.Error(), c)
		return
	}

	response.OKWithMessage(DefaultSuccessMessage, c)
}
