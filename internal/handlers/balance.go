package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Soliard/gophermart/internal/logger"
	"github.com/Soliard/gophermart/internal/services"
)

type balanceHandler struct {
	service services.BalanceServiceInterface
}

func NewBalanceHandler(balanceService services.BalanceServiceInterface) *balanceHandler {
	return &balanceHandler{
		service: balanceService,
	}
}

func (h *balanceHandler) GetBalance(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := logger.FromContext(ctx)

	userCtx, err := services.GetUserFromContext(ctx)
	if err != nil {
		log.Error("Failed to get user context from ctx after authentication", logger.F.Error(err))
		http.Error(res, "Failed to get user context", http.StatusInternalServerError)
		return
	}

	balance, err := h.service.GetBalance(ctx, userCtx.ID)
	if err != nil {
		log.Error("Failed to get users balance", logger.F.Error(err), logger.F.Any("user", userCtx))
		http.Error(res, "Failed to get balance", http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(balance)
	if err != nil {
		log.Error("Failed to masrshal balance", logger.F.Error(err), logger.F.Any("user", userCtx))
		http.Error(res, "Failed to masrshal balance", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(body)
}
