package services

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Soliard/gophermart/internal/dto"
	"github.com/Soliard/gophermart/internal/errs"
	"github.com/Soliard/gophermart/internal/logger"
	"github.com/Soliard/gophermart/internal/models"
	"github.com/go-resty/resty/v2"
)

type AccrualUpdater interface {
	GetOrdersToAccrualUpdate(ctx context.Context) ([]*models.Order, error)
	UpdateStatusAndAccural(
		ctx context.Context, numberOrder string,
		status models.OrderStatus, accrual *float64) error
}

type accrualService struct {
	updater AccrualUpdater
	client  *resty.Client
	baseURL string
}

func NewAccrualService(orders AccrualUpdater, accrualURL string) *accrualService {
	if !strings.HasPrefix(accrualURL, "http://") && !strings.HasPrefix(accrualURL, "https://") {
		accrualURL = "http://" + accrualURL
	}

	return &accrualService{
		updater: orders,
		client:  resty.New(),
		baseURL: accrualURL,
	}
}

func (s *accrualService) UpdateOrders(ctx context.Context) error {
	log := logger.FromContext(ctx)
	orders, err := s.updater.GetOrdersToAccrualUpdate(ctx)
	if err != nil {
		return err
	}
	for _, v := range orders {
		err := s.updateOrder(ctx, v.Number)
		if err != nil {
			log.Error("Failed to update order", logger.F.Any("order", v), logger.F.Error(err))
		}
	}
	return nil
}

func (s *accrualService) updateOrder(ctx context.Context, number string) error {
	fullURL, err := url.JoinPath(s.baseURL, "api", "orders", number)
	if err != nil {
		return err
	}

	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err := s.client.R().Get(fullURL)
		if err != nil {
			return err
		}
		status := resp.StatusCode()
		switch status {
		case http.StatusOK:
		case http.StatusNoContent:
			return nil
		case http.StatusTooManyRequests:
			time.Sleep(time.Second * 60)
			continue
		default:
			return errs.ErrUnexpectedStatusAccrualService
		}

		var recievedOrder dto.AccrualOrder
		err = json.Unmarshal(resp.Body(), &recievedOrder)
		if err != nil {
			return err
		}
		return s.updateStatusAndAccural(ctx, recievedOrder.Order, models.OrderStatus(recievedOrder.Status), recievedOrder.Accrual)
	}

	return nil
}

func (s *accrualService) updateStatusAndAccural(
	ctx context.Context, number string,
	status models.OrderStatus, accural *float64) error {

	return s.updater.UpdateStatusAndAccural(ctx, number, status, accural)
}
