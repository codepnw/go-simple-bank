package auth

import (
	"github.com/codepnw/simple-bank/internal/modules/user"
	"github.com/codepnw/simple-bank/internal/utils/response"
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
		response.ErrBadRequest(ctx, err)
		return
	}

	result, err := h.uc.Login(ctx, req)
	if err != nil {
		response.ErrInternalServer(ctx, err)
		return
	}

	response.Success(ctx, result)
}

func (h *authHandler) Register(ctx *gin.Context) {
	req := new(user.UserRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ErrBadRequest(ctx, err)
		return
	}

	result, err := h.uc.Register(ctx, req)
	if err != nil {
		response.ErrInternalServer(ctx, err)
		return
	}

	response.Created(ctx, result)
}
