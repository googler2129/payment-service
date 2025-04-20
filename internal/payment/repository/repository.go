package repository

import (
	"context"
	"sync"

	"github.com/mercor/payment-service/internal/domain"
	"github.com/mercor/payment-service/pkg/db/sql/postgres"
	"github.com/mercor/payment-service/pkg/repository/scd"
	"gorm.io/gorm"
)

var (
	repo     *PaymentRepository
	repoOnce sync.Once
)

type PaymentRepository struct {
	db *postgres.DbCluster
	scd.SCDRepository[domain.PaymentLineItem]
}

func NewPaymentRepository(db *postgres.DbCluster) domain.PaymentLineRepository {
	repoOnce.Do(func() {
		repo = &PaymentRepository{
			db:            db,
			SCDRepository: scd.NewSCDRepository(db, domain.PaymentLineItem{}),
		}
	})

	return repo
}

func (r *PaymentRepository) FindByContractorAndPeriod(ctx context.Context, contractorID string, startTime, endTime int64) ([]domain.PaymentLineItem, error) {
	paymentItems, err := r.CustomQuery(ctx, func(db *gorm.DB) *gorm.DB {
		return db.
			Joins("JOIN job ON job.uid = payment_line_items.job_uid").
			Joins("JOIN timelog ON timelog.uid = payment_line_items.timelog_uid").
			Where("job.contractor_id = ? AND timelog.time_start > ? AND timelog.time_end < ?", contractorID, startTime, endTime)
	})
	if err != nil {
		return nil, err
	}
	return paymentItems, nil
}
