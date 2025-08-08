package user

import (
	"github.com/codepnw/simple-bank/internal/utils"
	"github.com/codepnw/simple-bank/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	uc UserUsecase
}

func NewUserHandler(uc UserUsecase) *userHandler {
	return &userHandler{uc: uc}
}

func (h *userHandler) CreateUser(ctx *gin.Context) {
	req := new(UserRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ErrBadRequest(ctx, err)
		return
	}

	result, err := h.uc.Create(ctx, req)
	if err != nil {
		response.ErrInternalServer(ctx, err)
		return
	}

	response.Created(ctx, result)
}

func (h *userHandler) GetUser(ctx *gin.Context) {
	id, err := utils.GetParamID(ctx, "id")
	if err != nil {
		response.ErrBadRequest(ctx, err)
	}

	user, err := h.uc.GetUserByID(ctx, id)
	if err != nil {
		response.ErrInternalServer(ctx, err)
		return
	}

	response.Success(ctx, user)
}

func (h *userHandler) GetUsers(ctx *gin.Context) {
	users, err := h.uc.GetUsers(ctx)
	if err != nil {
		response.ErrInternalServer(ctx, err)
		return
	}

	response.Success(ctx, users)
}

func (h *userHandler) UpdateUser(ctx *gin.Context) {
	id, err := utils.GetParamID(ctx, "id")
	if err != nil {
		response.ErrBadRequest(ctx, err)
	}

	req := new(UserUpdateRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ErrBadRequest(ctx, err)
		return
	}

	result, err := h.uc.Update(ctx, id, req)
	if err != nil {
		response.ErrInternalServer(ctx, err)
		return
	}

	response.Success(ctx, result)
}

func (h *userHandler) DeleteUser(ctx *gin.Context) {
	id, err := utils.GetParamID(ctx, "id")
	if err != nil {
		response.ErrBadRequest(ctx, err)
	}

	if err := h.uc.Delete(ctx, id); err != nil {
		response.ErrInternalServer(ctx, err)
		return
	}

	response.Success(ctx, "user deleted")
}
