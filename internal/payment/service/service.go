package service

import (
	"context"

	"github.com/mercor/payment-service/internal/controller/payment/request"
	"github.com/mercor/payment-service/internal/domain"
)

type PaymentService struct {
	repo domain.PaymentLineRepository
}

func NewPaymentService(repo domain.PaymentLineRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) GetPaymentLineItemsForContractorPeriod(ctx context.Context, contractorID string, startTime, endTime int64) ([]domain.PaymentLineItem, error) {
	return s.repo.FindByContractorAndPeriod(ctx, contractorID, startTime, endTime)
}

func (s *PaymentService) UpdatePaymentLineItemByID(ctx context.Context, id string, paymentLineItemReq *request.UpdatePaymentRequest) error {
	paymentLineItem := domain.NewPaymentLineItem(
		paymentLineItemReq.JobUID,
		paymentLineItemReq.TimelogUID,
		paymentLineItemReq.Amount,
		paymentLineItemReq.Status,
	)
	return s.repo.Update(ctx, id, paymentLineItem)
}
