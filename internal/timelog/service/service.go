package service

import (
	"context"

	"github.com/mercor/payment-service/internal/domain"
)

type TimelogService struct {
	repo domain.TimeLogRepositoryInterface
}

func NewTimelogService(repo domain.TimeLogRepositoryInterface) *TimelogService {
	return &TimelogService{repo: repo}
}

func (s *TimelogService) GetTimelogsForContractorPeriod(ctx context.Context, contractorID string, startDate, endDate int64) ([]domain.Timelog, error) {
	return s.repo.FindByContractorAndPeriod(ctx, contractorID, startDate, endDate)
}
