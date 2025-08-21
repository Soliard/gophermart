package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Soliard/gophermart/internal/services"
)

type Handlers struct {
	User       *userHandler
	Order      *orderHandler
	Balance    *balanceHandler
	Withdrawal *withdrawalHandler
}

func New(services *services.Services) *Handlers {
	return &Handlers{
		User:       NewUserHandler(services.Reg, services.Auth),
		Order:      NewOrderHandler(services.Order),
		Balance:    NewBalanceHandler(services.Balance),
		Withdrawal: NewWithdrawalHandler(services.Withdrawal),
	}
}

func validateContentType(r *http.Request, expected string) bool {
	ct := r.Header.Get("Content-Type")
	return strings.HasPrefix(ct, expected)
}

func validateJSONContentType(req *http.Request) bool {
	return validateContentType(req, "application/json")
}

func validateTextContentType(req *http.Request) bool {
	return validateContentType(req, "text/plain")
}

func handleJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}
