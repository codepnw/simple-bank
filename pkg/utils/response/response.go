package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, message string, data any) {
	if message == "" {
		message = "successfully"
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": message,
		"data":    data,
	})
}

func Created(c *gin.Context, message string, data any) {
	if message == "" {
		message = "created"
	}
	c.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": message,
		"data":    data,
	})
}

func NoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code":    http.StatusBadRequest,
		"type":    "BAD_REQUEST",
		"message": message,
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"code":    http.StatusUnauthorized,
		"type":    "UNAUTHORIZED",
		"message": message,
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"type":    "NOT_FOUND",
		"message": message,
	})
}

func InternalServerError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":  http.StatusInternalServerError,
		"type":  "INTERNAL_SERVER_ERROR",
		"error": err.Error(),
	})
}
