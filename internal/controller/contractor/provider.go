package contractor

import (
	"github.com/google/wire"
	repository "github.com/mercor/payment-service/internal/contractor/repository"
	service "github.com/mercor/payment-service/internal/contractor/service"
	"github.com/mercor/payment-service/internal/domain"
)

var ProviderSet wire.ProviderSet = wire.NewSet(
	NewController,
	service.NewContractorService,
	repository.NewContractorRepository,

	wire.Bind(new(domain.ContractorControllerInterface), new(*Controller)),
	wire.Bind(new(domain.ContractorServiceInterface), new(*service.ContractorService)),
)
