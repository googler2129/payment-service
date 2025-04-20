package repository

import (
	"sync"

	"github.com/mercor/payment-service/internal/domain"
	"github.com/mercor/payment-service/pkg/db/sql/postgres"
	"github.com/mercor/payment-service/pkg/repository/scd"
)

var (
	repo     *JobRepository
	repoOnce sync.Once
)

type JobRepository struct {
	scd.SCDRepository[domain.Job]
	db *postgres.DbCluster
}

func NewJobRepository(db *postgres.DbCluster) domain.JobRepositoryInterface {
	repoOnce.Do(func() {
		repo = &JobRepository{
			db:            db,
			SCDRepository: scd.NewSCDRepository(db, domain.Job{}),
		}
	})

	return repo
}
