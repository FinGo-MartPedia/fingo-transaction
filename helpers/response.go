package helpers

import "github.com/gin-gonic/gin"

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   *string     `json:"error"`
}

func SendResponse(c *gin.Context, statusCode int, message string, data interface{}, err *string) {
	c.JSON(statusCode, Response{
		Message: message,
		Data:    data,
		Error:   err,
	})
}
