package user

import (
	"net/http"
	"strconv"

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
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	creasted, err := h.uc.Create(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": creasted})
}

func (h *userHandler) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)

	user, err := h.uc.GetUser(ctx, idInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *userHandler) GetUsers(ctx *gin.Context) {
	users, err := h.uc.GetUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": users})
}

func (h *userHandler) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)

	req := new(UserUpdateRequest)

	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	updated, err := h.uc.Update(ctx, idInt, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updated})
}

func (h *userHandler) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)

	if err := h.uc.Delete(ctx, idInt); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "user deleted"})
}
