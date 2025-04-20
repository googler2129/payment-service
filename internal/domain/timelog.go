package domain

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mercor/payment-service/pkg/repository/scd"
)

// Timelog represents a time logging entity in the system
type Timelog struct {
	*scd.SCDModel
	Duration  int64  `gorm:"column:duration;not null"`
	TimeStart int64  `gorm:"column:time_start;not null"`
	TimeEnd   int64  `gorm:"column:time_end;not null"`
	Type      string `gorm:"column:type;not null"`
	JobUID    string `gorm:"column:job_uid;not null;index"`
}

// TableName specifies the table name for the Timelog model
func (Timelog) TableName() string {
	return "timelog"
}

type TimeLogRepositoryInterface interface {
	scd.SCDRepository[Timelog]
	FindByContractorAndPeriod(ctx context.Context, contractorID string, startDate, endDate int64) ([]Timelog, error)
}

type TimelogServiceInterface interface {
	GetTimelogsForContractorPeriod(ctx context.Context, contractorID string, startDate, endDate int64) ([]Timelog, error)
}

type TimelogControllerInterface interface {
	GetTimelogsForContractorPeriod(ctx *gin.Context)
}
