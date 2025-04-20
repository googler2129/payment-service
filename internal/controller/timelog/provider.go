package timelog

import (
	"github.com/google/wire"
	"github.com/mercor/payment-service/internal/domain"
	repository "github.com/mercor/payment-service/internal/timelog/repository"
	service "github.com/mercor/payment-service/internal/timelog/service"
)

var ProviderSet wire.ProviderSet = wire.NewSet(
	NewTimelogController,
	service.NewTimelogService,
	repository.NewTimelogRepository,

	wire.Bind(new(domain.TimelogControllerInterface), new(*TimelogController)),
	wire.Bind(new(domain.TimelogServiceInterface), new(*service.TimelogService)),
)
