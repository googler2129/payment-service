package service

import (
	"context"

	"github.com/mercor/payment-service/internal/contractor/request"
	"github.com/mercor/payment-service/internal/domain"
	"github.com/mercor/payment-service/pkg/repository/static"
)

type ContractorService struct {
	repo domain.ContractorRepository
}

func NewContractorService(repo domain.ContractorRepository) *ContractorService {
	return &ContractorService{repo: repo}
}

func (s *ContractorService) CreateContractor(ctx context.Context, req *request.CreateContractorSvcReq) error {
	contractor := &domain.Contractor{
		Model: &static.Model{},
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}
	return s.repo.Create(ctx, contractor)
}
