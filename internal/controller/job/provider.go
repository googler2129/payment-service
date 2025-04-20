package job

import (
	"github.com/google/wire"
	companyRepository "github.com/mercor/payment-service/internal/company/repository"
	contractorRepository "github.com/mercor/payment-service/internal/contractor/repository"
	"github.com/mercor/payment-service/internal/domain"
	"github.com/mercor/payment-service/internal/job/repository"
	service "github.com/mercor/payment-service/internal/job/service"
)

var ProviderSet wire.ProviderSet = wire.NewSet(
	NewController,
	service.NewService,
	repository.NewJobRepository,
	companyRepository.NewCompanyRepository,
	contractorRepository.NewContractorRepository,

	wire.Bind(new(domain.JobControllerInterface), new(*Controller)),
	wire.Bind(new(domain.JobServiceInterface), new(*service.Service)),
)
