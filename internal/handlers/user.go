package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Soliard/gophermart/internal/dto"
	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/logger"
	"github.com/Soliard/gophermart/internal/services"
)

type userHandler struct {
	reg  services.RegistrationServiceInterface
	auth services.AuthServiceInterface
}

func NewUserHandler(reg services.RegistrationServiceInterface, auth services.AuthServiceInterface) *userHandler {
	return &userHandler{
		reg:  reg,
		auth: auth,
	}
}

func (h *userHandler) Register(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := logger.FromContext(ctx)

	if !validateJSONContentType(req) {
		http.Error(res, "Incorrect body format", http.StatusBadRequest)
		return
	}

	regData := &dto.RegisterRequest{}
	err := json.NewDecoder(req.Body).Decode(regData)
	if err != nil {
		log.Error("Failed to decode body", logger.F.Error(err))
		http.Error(res, "Failed to decode body", http.StatusBadRequest)
		return
	}

	u, err := h.reg.Register(ctx, regData)
	if err != nil {
		if errors.Is(err, errs.ErrUserAlreadyExists) {
			http.Error(res, "User with this login already exists", http.StatusConflict)
			return
		}
		if errors.Is(err, errs.ErrEmptyLoginOrPassword) {
			http.Error(res, "Login and password must be not empty", http.StatusBadRequest)
			return
		}
		log.Error("Failed to register user", logger.F.Error(err))
		http.Error(res, "Failed to register user", http.StatusInternalServerError)
		return
	}

	logData := &dto.LoginRequest{Login: regData.Login, Password: regData.Password}
	token, err := h.auth.Login(ctx, logData)
	if err != nil {
		log.Error("Failed to login after register (shouldnt happen)",
			logger.F.Error(err),
			logger.F.Any("regData", regData),
			logger.F.Any("user", u),
			logger.F.Any("logData", logData),
		)
		http.Error(res, "Failed to login after registration", http.StatusInternalServerError)
		return
	}
	res.Header().Add("Authorization", token)
	err = handleJSONResponse(res, http.StatusOK, u)
	if err != nil {
		log.Error("Failed to marshal body", logger.F.Error(err))
	}
}

func (h *userHandler) Login(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := logger.FromContext(ctx)

	if !validateJSONContentType(req) {
		http.Error(res, "Only application/json is allowed", http.StatusBadRequest)
		return
	}

	logData := &dto.LoginRequest{}
	err := json.NewDecoder(req.Body).Decode(logData)
	if err != nil {
		log.Error("Failed to decode body", logger.F.Error(err))
		http.Error(res, "Failed to decode body", http.StatusBadRequest)
		return
	}
	if logData.Login == "" || logData.Password == "" {
		http.Error(res, "Login and password must be not empty", http.StatusBadRequest)
		return
	}

	token, err := h.auth.Login(ctx, logData)
	if err != nil {
		if errors.Is(err, errs.ErrWrongLoginOrPassword) {
			http.Error(res, "Wrong login or password", http.StatusUnauthorized)
			return
		}
		log.Error("Failed to login user", logger.F.Error(err))
		http.Error(res, "Failed to login, try later", http.StatusInternalServerError)
		return
	}

	res.Header().Add("Authorization", token)
	res.WriteHeader(http.StatusOK)
}
