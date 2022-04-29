package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hiro942/elden-client/global"
	"github.com/hiro942/elden-client/model/response"
	"strings"
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
	var (
		DefaultSuccessMessage = "pre-handover for receiving new satellite id success"
		DefaultErrorMessage   = "pre-handover for receiving new satellite id error"
	)

	newSatelliteId := c.Query("id")
	if strings.TrimSpace(newSatelliteId) == "" {
		response.FailWithDescription(DefaultErrorMessage, "blank new satellite id", c)
	}

	// 更新交接卫星
	global.HandoverSatellite = newSatelliteId

	response.OKWithMessage(DefaultSuccessMessage, c)
}
