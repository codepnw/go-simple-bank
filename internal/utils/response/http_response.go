package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Created(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    data,
	})
}

func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

func Unauthorized(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusUnauthorized, gin.H{"message": msg})
}

func ErrInternalServer(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}

func ErrBadRequest(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}
