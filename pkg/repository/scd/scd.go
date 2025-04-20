package scd

import (
	"context"

	"gorm.io/gorm"
)

// SCDRecord is an interface that must be implemented by all SCD models
type SCDRecord interface {
	GetID() string
	GetVersion() int
	GetUID() string
	SetVersion(version int)
	SetID(id string)
	SetIsLatest(isLatest bool)
	SetUID(id string)
}

// SCDRepository is a generic repository for SCD tables
type SCDRepository[T SCDRecord] interface {
	FindByID(ctx context.Context, id string) (*T, error)

	FindByUID(ctx context.Context, uid string) (*T, error)

	FindAllLatest(ctx context.Context) ([]T, error)

	FindLatestWithFilter(ctx context.Context, filter map[string]interface{}) ([]T, error)

	FindVersionsForID(ctx context.Context, id string) ([]T, error)

	Create(ctx context.Context, record *T) error

	Update(ctx context.Context, id string, record *T) error

	CustomQuery(ctx context.Context, queryBuilder func(*gorm.DB) *gorm.DB) ([]T, error)
}
