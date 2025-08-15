package user

import (
	"errors"

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

	result, err := h.uc.GetUserByID(ctx, id)
	if err != nil {
		response.ErrInternalServer(ctx, err)
		return
	}

	response.Success(ctx, result)
}

func (h *userHandler) GetProfile(ctx *gin.Context) {
	u, err := CurrentUser(ctx)
	if err != nil {
		response.Unauthorized(ctx, err.Error())
		return
	}

	response.Success(ctx, u)
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

func (h *userHandler) UpdateProfile(ctx *gin.Context) {
	u, err := CurrentUser(ctx)
	if err != nil {
		response.Unauthorized(ctx, err.Error())
		return
	}

	req := new(UserUpdateRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ErrBadRequest(ctx, err)
		return
	}

	result, err := h.uc.Update(ctx, u.ID, req)
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

func CurrentUser(ctx *gin.Context) (*User, error) {
	val, ok := ctx.Get("user")
	if !ok {
		return nil, errors.New("user not found in context")
	}

	u, ok := val.(*User)
	if !ok {
		return nil, errors.New("invalid user type in context")
	}

	return u, nil
}
