package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Soliard/gophermart/internal/dto"
	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/logger"
	"github.com/Soliard/gophermart/internal/services"
)

type withdrawalHandler struct {
	service services.WithdrawalServiceInterface
}

func NewWithdrawalHandler(service services.WithdrawalServiceInterface) *withdrawalHandler {
	return &withdrawalHandler{service: service}
}

func (h *withdrawalHandler) ProcessWithdrawal(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := logger.FromContext(ctx)

	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(res, "Incorrect body format", http.StatusBadRequest)
		return
	}

	userCtx, err := services.GetUserFromContext(ctx)
	if err != nil {
		log.Error("Failed to get user context from ctx after authentication", logger.F.Error(err))
		http.Error(res, "Failed to get user context", http.StatusInternalServerError)
		return
	}

	reqData := &dto.WithdrawalRequest{}
	err = json.NewDecoder(req.Body).Decode(reqData)
	if err != nil {
		log.Error("Failed to decode body", logger.F.Error(err))
		http.Error(res, "Failed to decode body", http.StatusInternalServerError)
		return
	}

	if reqData.Sum <= 0 {
		http.Error(res, "Sum must be positive", http.StatusBadRequest)
		return
	}

	err = h.service.ProcessWithdraw(ctx, userCtx.ID, reqData.Order, reqData.Sum)
	switch err {
	case errs.OrderIsNotValid:
		http.Error(res, "Order number is not valid", http.StatusUnprocessableEntity)
		return

	case errs.WithdrawAlreadyProcessed:
		log.Warn("Attempt withdrawal order that already has been withdrawed", logger.F.Any("request data", reqData))
		http.Error(res, "For this order already exists withdrawal", http.StatusAlreadyReported)
		return

	case errs.BalanceInsufficient:
		http.Error(res, "Not enough points on balance", http.StatusPaymentRequired)
		return
	}

	res.WriteHeader(http.StatusOK)
}
