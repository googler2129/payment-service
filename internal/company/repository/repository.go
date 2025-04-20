package repository

import (
	"sync"

	"github.com/mercor/payment-service/internal/domain"
	"github.com/mercor/payment-service/pkg/db/sql/postgres"
	"github.com/mercor/payment-service/pkg/repository/static"
)

var (
	repo     *CompanyRepository
	repoOnce sync.Once
)

type CompanyRepository struct {
	static.StaticRepository[domain.Company]
	db *postgres.DbCluster
}

func NewCompanyRepository(db *postgres.DbCluster) domain.CompanyRepository {
	repoOnce.Do(func() {
		repo = &CompanyRepository{
			db:               db,
			StaticRepository: static.NewStaticRepository[domain.Company](db),
		}
	})

	return repo
}
