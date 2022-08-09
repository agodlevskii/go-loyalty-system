package jobs

import (
	"github.com/robfig/cron"
	"go-loyalty-system/internal/services"
	"go-loyalty-system/internal/storage"
)

func OrderStatusUpdate(db storage.Repo, accrual services.AccrualClient) error {
	c := cron.New()
	return c.AddFunc("@hourly", func() {
		services.UpdateOrdersStatus(db, accrual)
	})
}
