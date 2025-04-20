package repository

import (
	"sync"

	"github.com/mercor/payment-service/internal/domain"
	"github.com/mercor/payment-service/pkg/db/sql/postgres"
	"github.com/mercor/payment-service/pkg/repository/static"
)

var (
	repo     *ContractorRepository
	repoOnce sync.Once
)

type ContractorRepository struct {
	static.StaticRepository[domain.Contractor]
	db *postgres.DbCluster
}

func NewContractorRepository(db *postgres.DbCluster) domain.ContractorRepository {
	repoOnce.Do(func() {
		repo = &ContractorRepository{
			db:               db,
			StaticRepository: static.NewStaticRepository[domain.Contractor](db),
		}
	})

	return repo
}
