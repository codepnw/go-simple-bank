package account

import (
	"github.com/codepnw/simple-bank/internal/utils"
	"github.com/codepnw/simple-bank/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type accountHandler struct {
	uc AccountUsecase
}

func NewAccountHandler(uc AccountUsecase) *accountHandler {
	return &accountHandler{uc: uc}
}

func (h *accountHandler) CreateAccount(ctx *gin.Context) {
	req := new(accountRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ErrBadRequest(ctx, err)
		return
	}

	result, err := h.uc.CreateAccount(ctx, req)
	if err != nil {
		response.ErrInternalServer(ctx, err)
		return
	}

	response.Success(ctx, result)
}

func (h *accountHandler) GetAccountByID(ctx *gin.Context) {
	id, err := utils.GetParamID(ctx, "id")
	if err != nil {
		response.ErrBadRequest(ctx, err)
		return
	}

	result, err := h.uc.GetAccountByID(ctx, id)
	if err != nil {
		response.ErrInternalServer(ctx, err)
		return
	}

	response.Success(ctx, result)
}

func (h *accountHandler) ListAccounts(ctx *gin.Context) {
	// TODO: get user_id from middleware later
	userID, err := utils.GetParamID(ctx, "userID")
	if err != nil {
		response.ErrBadRequest(ctx, err)
		return
	}

	result, err := h.uc.ListAccounts(ctx, userID)
	if err != nil {
		response.ErrInternalServer(ctx, err)
		return
	}

	response.Success(ctx, result)
}

func (h *accountHandler) UpdateStatusPending(ctx *gin.Context) {
	id, err := utils.GetParamID(ctx, "id")
	if err != nil {
		response.ErrBadRequest(ctx, err)
		return
	}

	if err = h.uc.UpdateStatusPending(ctx, id); err != nil {
		response.ErrInternalServer(ctx, err)
		return
	}

	response.Success(ctx, "updated account pending")
}

func (h *accountHandler) UpdateStatusApproved(ctx *gin.Context) {
	id, err := utils.GetParamID(ctx, "id")
	if err != nil {
		response.ErrBadRequest(ctx, err)
		return
	}

	if err = h.uc.UpdateStatusApproved(ctx, id); err != nil {
		response.ErrInternalServer(ctx, err)
		return
	}

	response.Success(ctx, "updated account approved")
}

func (h *accountHandler) UpdateStatusRejected(ctx *gin.Context) {
	id, err := utils.GetParamID(ctx, "id")
	if err != nil {
		response.ErrBadRequest(ctx, err)
		return
	}

	if err = h.uc.UpdateStatusRejected(ctx, id); err != nil {
		response.ErrInternalServer(ctx, err)
		return
	}

	response.Success(ctx, "updated account rejected")
}
