package domain

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mercor/payment-service/internal/contractor/request"
	"github.com/mercor/payment-service/pkg/repository/static"
	"gorm.io/gorm"
)

type Contractor struct {
	*static.Model
	Name      string         `gorm:"column:name;not null"`
	Email     string         `gorm:"column:email;not null"`
	Phone     string         `gorm:"column:phone;not null"`
	CreatedAt time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (Contractor) TableName() string {
	return "contractor"
}

func NewContractor(name string, email string, phone string) *Contractor {
	return &Contractor{
		Model: &static.Model{},
		Name:  name,
		Email: email,
		Phone: phone,
	}
}

type ContractorRepository interface {
	static.StaticRepository[Contractor]
}

type ContractorServiceInterface interface {
	CreateContractor(ctx context.Context, req *request.CreateContractorSvcReq) error
}

type ContractorControllerInterface interface {
	CreateContractor(ctx *gin.Context)
}
