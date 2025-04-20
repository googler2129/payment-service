package static

import "context"

type Static interface {
	SetID(id string)
}

type StaticRepository[T Static] interface {
	Create(ctx context.Context, record *T) error
	CreateInBatch(ctx context.Context, records []*T, batchSize int) error
	UpdateByCondition(ctx context.Context, filter map[string]interface{}, record *T) error
	UpdatesByConditions(ctx context.Context, filter map[string]interface{}, updates map[string]interface{}) error
	Delete(ctx context.Context, record *T) error
	DeleteByConditions(ctx context.Context, filter map[string]interface{}) error
	GetByConditions(ctx context.Context, filter map[string]interface{}) (*T, error)
	GetAllByConditions(ctx context.Context, filter map[string]interface{}) ([]T, error)
}
