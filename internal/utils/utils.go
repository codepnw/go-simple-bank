package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const ContextKeyUser = "user"

func GetParamID(ctx *gin.Context, key string) (int64, error) {
	id := ctx.Param(key)
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, err
	}
	return idInt, nil
}
