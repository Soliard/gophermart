package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

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
	ct := req.Header.Get("Content-Type")

	if !strings.HasPrefix(ct, "application/json") {
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
		http.Error(res, "Failed to decode body", http.StatusBadRequest)
		return
	}

	if reqData.Sum <= 0 {
		http.Error(res, "Sum must be positive", http.StatusBadRequest)
		return
	}

	err = h.service.ProcessWithdraw(ctx, userCtx.ID, reqData.Order, reqData.Sum)
	switch err {
	case errs.ErrOrderIsNotValid:
		http.Error(res, "Order number is not valid", http.StatusUnprocessableEntity)
		return

	case errs.ErrWithdrawalAlreadyProcessed:
		log.Warn("Attempt withdrawal order that already has been withdrawed", logger.F.Any("request data", reqData))
		res.WriteHeader(http.StatusOK)
		return

	case errs.ErrBalanceInsufficient:
		http.Error(res, "Not enough points on balance", http.StatusPaymentRequired)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func (h *withdrawalHandler) GetWithdrawals(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := logger.FromContext(ctx)

	userCtx, err := services.GetUserFromContext(ctx)
	if err != nil {
		log.Error("Failed to get user context from ctx after authentication", logger.F.Error(err))
		http.Error(res, "Failed to get user context", http.StatusInternalServerError)
		return
	}

	withdrawals, err := h.service.GetWithdrawals(ctx, userCtx.ID)
	if err != nil {
		log.Error("Failed to get user withdrawals", logger.F.Error(err), logger.F.Any("user", userCtx))
		http.Error(res, "Failed to get withdrawals", http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	body, err := json.Marshal(withdrawals)
	if err != nil {
		log.Error("Failed to marshal withdrawals", logger.F.Error(err), logger.F.Any("user", userCtx))
		http.Error(res, "Failed to marshal withdrawals", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(body)
}
