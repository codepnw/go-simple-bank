package transaction

import (
	"github.com/codepnw/simple-bank/internal/utils"
	"github.com/codepnw/simple-bank/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type transactionHandler struct {
	uc       TransactionUsecase
	validate *validator.Validate
}

func NewTransactionHandler(uc TransactionUsecase) *transactionHandler {
	return &transactionHandler{
		uc:       uc,
		validate: validator.New(),
	}
}

func (h *transactionHandler) Deposit(ctx *gin.Context) {
	req := new(DepositReq)

	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ErrBadRequest(ctx, err)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ErrBadRequest(ctx, err)
		return
	}

	// Deposit Usecase
	result, err := h.uc.Deposit(ctx.Request.Context(), req)
	if err != nil {
		response.ErrInternalServer(ctx, err)
		return
	}

	response.Success(ctx, result)
}

func (h *transactionHandler) Withdraw(ctx *gin.Context) {
	req := new(WithdrawReq)

	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ErrBadRequest(ctx, err)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ErrBadRequest(ctx, err)
		return
	}

	// Withdraw Usecase
	result, err := h.uc.Withdraw(ctx.Request.Context(), req)
	if err != nil {
		response.ErrInternalServer(ctx, err)
		return
	}

	response.Success(ctx, result)
}

func (h *transactionHandler) Transfer(ctx *gin.Context) {
	req := new(TransferReq)

	if err := ctx.ShouldBindJSON(req); err != nil {
		response.ErrBadRequest(ctx, err)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ErrBadRequest(ctx, err)
		return
	}

	// Transfer Usecase
	result, err := h.uc.Transfer(ctx.Request.Context(), req)
	if err != nil {
		response.ErrInternalServer(ctx, err)
		return
	}

	response.Success(ctx, result)
}

func (h *transactionHandler) Transactions(ctx *gin.Context) {
	// TODO: get userID from context later
	userID, err := utils.GetParamID(ctx, "userID")
	if err != nil {
		response.ErrBadRequest(ctx, err)
		return
	}

	// Transactions Usecase
	result, err := h.uc.Transactions(ctx.Request.Context(), userID)
	if err != nil {
		response.ErrInternalServer(ctx, err)
		return
	}

	response.Success(ctx, result)
}
