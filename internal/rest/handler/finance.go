package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-finance-tracker/internal/models"
	"go-finance-tracker/internal/repository"
	"go-finance-tracker/internal/rest/form"
	"go-finance-tracker/pkg/logger"
	"net/http"
	"strconv"
)

type FinanceHandlers struct {
	financeRepo repository.FinanceRepo
}

func NewFinanceHandlers(financeRepo repository.FinanceRepo) *FinanceHandlers {
	return &FinanceHandlers{financeRepo: financeRepo}
}

func (h *FinanceHandlers) GetAllFinance(ctx *gin.Context) {
	userIdStr, exists := ctx.Get("id")
	if !exists {
		logger.GetLogger().Error("User not authenticated")
		ctx.JSON(http.StatusUnauthorized, &models.CustomResponse{
			Status:  http.StatusUnauthorized,
			Message: "User not authenticated",
		})
		return
	}

	id, _ := strconv.Atoi(userIdStr.(string))

	finances, err := h.financeRepo.GetAll(id)
	if err != nil {
		if errors.Is(err, repository.ErrFinanceHistoryNotFound) {
			ctx.JSON(http.StatusNotFound, &models.CustomResponse{
				Status:  http.StatusNotFound,
				Message: "User finance history not found:",
				Error:   err.Error(),
			})
		}
		ctx.JSON(http.StatusInternalServerError, &models.CustomResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, &models.CustomResponse{
		Status:  http.StatusOK,
		Message: "User Finance History fetched successfully",
		Data:    finances,
	})
}

func (h *FinanceHandlers) AddFinanceRecord(ctx *gin.Context) {
	userIdStr, exists := ctx.Get("id")
	if !exists {
		logger.GetLogger().Error("User not authenticated")
		ctx.JSON(http.StatusUnauthorized, &models.CustomResponse{
			Status:  http.StatusUnauthorized,
			Message: "User not authenticated",
		})
		return
	}
	userID, _ := strconv.Atoi(userIdStr.(string))

	var financeForm form.FinanceRecordInput
	if err := ctx.ShouldBindJSON(&financeForm); err != nil {
		logger.GetLogger().Error("Invalid finance record request:", err)
		ctx.JSON(http.StatusBadRequest, &models.CustomResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		})
		return
	}

	var financeRecord models.FinanceRecord

	financeRecord.UserID = uint(userID)
	financeRecord.Amount = financeForm.Amount
	financeRecord.TransactionTypeID = financeForm.TransactionTypeID
	financeRecord.CategoryID = financeForm.CategoryID
	financeRecord.Note = financeForm.Note

	if err := h.financeRepo.Create(&financeRecord); err != nil {
		logger.GetLogger().Error("Failed to create finance record:", err)
		ctx.JSON(http.StatusInternalServerError, &models.CustomResponse{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "ok"})
}
