package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hiro942/elden-client/model/response"
	"github.com/hiro942/elden-client/service"
	"log"
	"time"
)

// @Summary authentication for normal(fast) access phrase
// @Router /auth/prehandover/location [post]

func PreHandoverForLocation(c *gin.Context) {
	var (
	//DefaultSuccessMessage = "pre-handover for location success"
	//DefaultErrorMessage   = "pre-handover for location error"
	)

	// todo 估算自己的位置区间返回给卫星
	// todo 这里暂且先返回一个成功响应
	response.OK(c)
}

// @Summary authentication for normal(fast) access phrase
// @Router /auth/prehandover/new_satellite?id= [post]

func PreHandoverForNewSatellite(c *gin.Context) {
	newSatelliteId := c.Query("id")

	response.OK(c)

	// 2秒后接入新卫星
	time.Sleep(time.Second * 2)

	// 接入新卫星
	err := service.HandoverAccess(newSatelliteId)
	if err != nil {
		log.Panicln(err)
	}
}
