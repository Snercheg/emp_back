package response

import (
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusSuccess = "Success"
	StatusError   = "Error"
)

func NewErrorResponse(c *gin.Context, statusCode int, msg string) Response {
	slog.Error(msg)
	c.AbortWithStatusJSON(statusCode, gin.H{"status": StatusError, "error": msg})
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func SuccessResponse() Response {
	return Response{
		Status: StatusSuccess,
	}
}

func ErrorResponse(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}
