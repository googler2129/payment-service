package job

import (
	"context"
	"sync"

	"github.com/mercor/payment-service/internal/domain"
	"github.com/mercor/payment-service/internal/job/request"
)

type Service struct {
	jobRepo        domain.JobRepositoryInterface
	companyRepo    domain.CompanyRepository
	contractorRepo domain.ContractorRepository
}

var (
	svc     *Service
	svcOnce sync.Once
)

func NewService(repo domain.JobRepositoryInterface, companyRepo domain.CompanyRepository, contractorRepo domain.ContractorRepository) *Service {
	svcOnce.Do(func() {
		svc = &Service{jobRepo: repo, companyRepo: companyRepo, contractorRepo: contractorRepo}
	})
	return svc
}

func (s *Service) CreateJob(ctx context.Context, req *request.CreateJobSvcReq) error {
	err := s.validateCreateJobRequest(ctx, req)
	if err != nil {
		return err
	}

	j := domain.NewJob(
		req.Status,
		req.Rate,
		req.Title,
		req.CompanyID,
		req.ContractorID,
	)

	return s.jobRepo.Create(ctx, j)
}

func (s *Service) validateCreateJobRequest(ctx context.Context, req *request.CreateJobSvcReq) error {
	_, err := s.companyRepo.GetByConditions(ctx, map[string]interface{}{"id": req.CompanyID})
	if err != nil {
		return err
	}
	_, err = s.contractorRepo.GetByConditions(ctx, map[string]interface{}{"id": req.ContractorID})
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetJobsByStatus(ctx context.Context, status string) ([]domain.Job, error) {
	return s.jobRepo.FindLatestWithFilter(ctx, map[string]interface{}{"status": status})
}

func (s *Service) GetActiveJobsForContractor(ctx context.Context, contractorID string) ([]domain.Job, error) {
	return s.jobRepo.FindLatestWithFilter(ctx, map[string]interface{}{"contractor_id": contractorID, "status": "active"})
}
