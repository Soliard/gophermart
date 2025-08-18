package workers

import (
	"context"
	"time"

	"github.com/Soliard/gophermart/internal/logger"
	"github.com/Soliard/gophermart/internal/services"
)

type accrualUpdater struct {
	accrual  services.AccrualServiceInterface
	interval time.Duration
}

func NewAccrualUpdater(accrual services.AccrualServiceInterface, interval time.Duration) *accrualUpdater {
	return &accrualUpdater{
		accrual:  accrual,
		interval: interval,
	}
}

func (u *accrualUpdater) Start(ctx context.Context) {
	log := logger.FromContext(ctx)
	ticker := time.NewTicker(u.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Warn("Accrual updater stopped")
			return
		case <-ticker.C:
			err := u.accrual.UpdateOrders(ctx)
			if err != nil {
				log.Error("Failed to update accural orders", logger.F.Error(err))
			}
		}
	}
}
