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
	repo     *TimelogRepository
	repoOnce sync.Once
)

type TimelogRepository struct {
	db *postgres.DbCluster
	scd.SCDRepository[domain.Timelog]
}

func NewTimelogRepository(db *postgres.DbCluster) domain.TimeLogRepositoryInterface {
	repoOnce.Do(func() {
		repo = &TimelogRepository{
			db:            db,
			SCDRepository: scd.NewSCDRepository(db, domain.Timelog{}),
		}
	})

	return repo
}

func (r *TimelogRepository) FindByContractorAndPeriod(ctx context.Context, contractorID string, startDate, endDate int64) ([]domain.Timelog, error) {

	timelogs, err := r.CustomQuery(ctx, func(db *gorm.DB) *gorm.DB {
		return db.
			Joins("JOIN job ON job.uid = timelog.job_uid").
			Where("job.contractor_id = ? AND timelog.time_start > ? AND timelog.time_end < ?", contractorID, startDate, endDate)
	})
	if err != nil {
		return nil, err
	}

	return timelogs, nil
}
