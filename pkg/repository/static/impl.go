package static

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/mercor/payment-service/pkg/db/sql/postgres"
)

type staticRepositoryImpl[T Static] struct {
	db *postgres.DbCluster
}

func NewStaticRepository[T Static](db *postgres.DbCluster) StaticRepository[T] {
	return &staticRepositoryImpl[T]{
		db: db,
	}
}

func (r *staticRepositoryImpl[T]) Create(ctx context.Context, record *T) error {
	if record == nil {
		return errors.New("record cannot be nil")
	}
	(*record).SetID(uuid.New().String())
	return r.db.GetMasterDB(ctx).Create(record).Error
}

func (r *staticRepositoryImpl[T]) CreateInBatch(ctx context.Context, records []*T, batchSize int) error {
	for _, record := range records {
		(*record).SetID(uuid.New().String())
	}
	return r.db.GetMasterDB(ctx).CreateInBatches(records, batchSize).Error
}

func (r *staticRepositoryImpl[T]) UpdateByCondition(ctx context.Context, filter map[string]interface{}, record *T) error {
	return r.db.GetMasterDB(ctx).Where(filter).Updates(record).Error
}

func (r *staticRepositoryImpl[T]) UpdatesByConditions(ctx context.Context, filter map[string]interface{}, updates map[string]interface{}) error {
	return r.db.GetMasterDB(ctx).Where(filter).Updates(updates).Error
}

func (r *staticRepositoryImpl[T]) Delete(ctx context.Context, record *T) error {
	return r.db.GetMasterDB(ctx).Delete(record).Error
}

func (r *staticRepositoryImpl[T]) DeleteByConditions(ctx context.Context, filter map[string]interface{}) error {
	var t T
	return r.db.GetMasterDB(ctx).Where(filter).Delete(&t).Error
}

func (r *staticRepositoryImpl[T]) GetByConditions(ctx context.Context, filter map[string]interface{}) (*T, error) {
	var result T
	err := r.db.GetSlaveDB(ctx).Where(filter).First(&result).Error
	return &result, err
}

func (r *staticRepositoryImpl[T]) GetAllByConditions(ctx context.Context, filter map[string]interface{}) ([]T, error) {
	var results []T
	err := r.db.GetSlaveDB(ctx).Where(filter).Find(&results).Error
	return results, err
}
