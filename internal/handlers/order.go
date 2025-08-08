package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/logger"
	"github.com/Soliard/gophermart/internal/services"
)

type orderHandler struct {
	orderService services.OrderServiceInterface
}

func NewOrderHandler(orderService services.OrderServiceInterface) *orderHandler {
	return &orderHandler{
		orderService: orderService,
	}
}

func (h *orderHandler) UploadOrder(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := logger.FromContext(ctx)

	if req.Header.Get("Content-Type") != "text/plain" {
		http.Error(res, "Incorrect body format", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.Error("Failed to read body", logger.F.Error(err))
		http.Error(res, "Failed to read body", http.StatusInternalServerError)
		return
	}

	orderNumber := string(body)
	isValid := h.orderService.ValidateOrderNumber(ctx, orderNumber)
	if !isValid {
		http.Error(res, "Order number is not valid", http.StatusUnprocessableEntity)
		return
	}

	userCtx, err := services.GetUserFromContext(ctx)
	if err != nil {
		log.Error("Failed to get user context from ctx after authentication", logger.F.Error(err))
		http.Error(res, "Failed to get user context", http.StatusInternalServerError)
		return
	}

	_, err = h.orderService.UploadOrder(ctx, userCtx.ID, orderNumber)
	if err != nil {
		if errors.Is(err, errs.OrderAlreadyUploadedByOtherUser) {
			http.Error(res, "Order already uploaded by other user", http.StatusConflict)
			return
		}
		if errors.Is(err, errs.OrderAlreadyUploadedByThisUser) {
			res.WriteHeader(http.StatusOK)
			return
		}
		log.Error("Failed to upload order", logger.F.Error(err))
		http.Error(res, "Failed to upload order", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusAccepted)
}
