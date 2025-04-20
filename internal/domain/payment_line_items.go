package domain

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mercor/payment-service/pkg/repository/scd"
)

// PaymentLineItem represents a payment line item entity in the system
type PaymentLineItem struct {
	*scd.SCDModel
	JobUID     string  `gorm:"column:job_uid;not null;index"`
	TimelogUID string  `gorm:"column:timelog_uid;not null;index"`
	Amount     float64 `gorm:"column:amount;not null"`
	Status     string  `gorm:"column:status;not null"`
}

// TableName specifies the table name for the PaymentLineItem model
func (PaymentLineItem) TableName() string {
	return "payment_line_items"
}

func NewPaymentLineItem(jobUID, timelogUID string, amount float64, status string) *PaymentLineItem {
	return &PaymentLineItem{
		SCDModel:   &scd.SCDModel{},
		JobUID:     jobUID,
		TimelogUID: timelogUID,
		Amount:     amount,
		Status:     status,
	}
}

type PaymentLineRepository interface {
	scd.SCDRepository[PaymentLineItem]
	FindByContractorAndPeriod(ctx context.Context, contractorID string, startTime, endTime int64) ([]PaymentLineItem, error)
}

type PaymentLineServiceInterface interface {
	GetPaymentLineItemsForContractorPeriod(ctx context.Context, contractorID string, startTime, endTime int64) ([]PaymentLineItem, error)
	UpdatePaymentLineItemByID(ctx context.Context, id string, updates map[string]interface{}) error
}

type PaymentLineControllerInterface interface {
	GetPaymentLineItemsForContractorPeriod(ctx *gin.Context)
	UpdatePaymentLineItemByID(ctx *gin.Context)
}
