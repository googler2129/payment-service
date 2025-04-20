package domain

import (
	"time"

	"github.com/mercor/payment-service/pkg/repository/static"
	"gorm.io/gorm"
)

type Company struct {
	*static.Model
	Name      string         `gorm:"column:name;not null"`
	CreatedAt time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func NewCompany(name string, createdAt time.Time, updatedAt time.Time, deletedAt gorm.DeletedAt) *Company {
	return &Company{
		Model:     &static.Model{},
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}
}

func (Company) TableName() string {
	return "company"
}

type CompanyRepository interface {
	static.StaticRepository[Company]
}
