package domain

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mercor/payment-service/internal/job/request"
	"github.com/mercor/payment-service/pkg/repository/scd"
)

// Job represents a job entity in the system
type Job struct {
	*scd.SCDModel
	Status       string  `gorm:"column:status;not null"`
	Rate         float64 `gorm:"column:rate;not null"`
	Title        string  `gorm:"column:title;not null"`
	CompanyID    string  `gorm:"column:company_id;not null"`
	ContractorID string  `gorm:"column:contractor_id;not null"`
}

// TableName specifies the table name for the Job model
func (Job) TableName() string {
	return "job"
}

func NewJob(status string, rate float64, title string, companyID string, contractorID string) *Job {
	return &Job{
		SCDModel:     &scd.SCDModel{},
		Status:       status,
		Rate:         rate,
		Title:        title,
		CompanyID:    companyID,
		ContractorID: contractorID,
	}
}

type JobRepositoryInterface interface {
	scd.SCDRepository[Job]
}

type JobServiceInterface interface {
	CreateJob(ctx context.Context, req *request.CreateJobSvcReq) error
	GetJobsByStatus(ctx context.Context, status string) ([]Job, error)
	GetActiveJobsForContractor(ctx context.Context, contractorID string) ([]Job, error)
}

type JobControllerInterface interface {
	CreateJob(ctx *gin.Context)
	GetJobsByStatus(ctx *gin.Context)
	GetActiveJobsForContractor(ctx *gin.Context)
}
