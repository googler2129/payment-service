package payment

import (
	"github.com/google/wire"
	"github.com/mercor/payment-service/internal/domain"
	"github.com/mercor/payment-service/internal/payment/repository"
	"github.com/mercor/payment-service/internal/payment/service"
)

var ProviderSet wire.ProviderSet = wire.NewSet(
	NewPaymentController,
	service.NewPaymentService,
	repository.NewPaymentRepository,

	wire.Bind(new(domain.PaymentLineControllerInterface), new(*PaymentController)),
	wire.Bind(new(domain.PaymentLineServiceInterface), new(*service.PaymentService)),
)
