package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response[T any] struct {
	Code        int    `json:"code"`
	Data        T      `json:"data"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

const (
	Success = 0
	Error   = 7
)

func Result[T any](code int, data T, message string, description string, c *gin.Context) {
	c.JSON(http.StatusOK, Response[T]{
		code,
		data,
		message,
		description,
	})
}

func OK(c *gin.Context) {
	Result[any](Success, nil, "ok", "", c)
}

func OKWithMessage(message string, c *gin.Context) {
	Result[any](Success, nil, message, "", c)
}

func OKWithData[T any](data T, message string, c *gin.Context) {
	Result[T](Success, data, message, "", c)
}

func FailWithDescription(message string, description string, c *gin.Context) {
	Result[any](Error, nil, message, description, c)
}
