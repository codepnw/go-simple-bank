package auth

import (
	"net/http"

	"github.com/codepnw/simple-bank/internal/modules/user"
	"github.com/gin-gonic/gin"
)

type authHandler struct {
	uc AuthUsecase
}

func NewAuthHandler(uc AuthUsecase) *authHandler {
	return &authHandler{uc: uc}
}

func (h *authHandler) Login(ctx *gin.Context) {
	req := new(authRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	response, err := h.uc.Login(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *authHandler) Register(ctx *gin.Context) {
	req := new(user.UserRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	response, err := h.uc.Register(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": response})
}
