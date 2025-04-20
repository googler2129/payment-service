package scd

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mercor/payment-service/pkg/db/sql/postgres"
	"gorm.io/gorm"
)

// scdRepositoryImpl is the implementation of SCDRepository
type scdRepositoryImpl[T SCDRecord] struct {
	db        *postgres.DbCluster
	modelType T
}

func NewSCDRepository[T SCDRecord](db *postgres.DbCluster, modelType T) SCDRepository[T] {
	return &scdRepositoryImpl[T]{
		db:        db,
		modelType: modelType,
	}
}

// FindByID returns the latest version of a record by ID
func (r *scdRepositoryImpl[T]) FindByID(ctx context.Context, id string) (*T, error) {
	var result T
	err := r.db.GetSlaveDB(ctx).
		Where("id = ? AND is_latest = ?", id, true).
		First(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find record by ID: %w", err)
	}

	return &result, nil
}

// FindByUID returns a specific version of a record by UID
func (r *scdRepositoryImpl[T]) FindByUID(ctx context.Context, uid string) (*T, error) {
	var result T
	err := r.db.GetSlaveDB(ctx).
		Where("uid = ?", uid).
		First(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find record by UID: %w", err)
	}

	return &result, nil
}

// FindAllLatest returns all latest versions of records
func (r *scdRepositoryImpl[T]) FindAllLatest(ctx context.Context) ([]T, error) {
	var results []T
	err := r.db.GetSlaveDB(ctx).
		Where("is_latest = ?", true).
		Find(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find all latest records: %w", err)
	}

	return results, nil
}

// FindLatestWithFilter returns latest versions that match the filter
func (r *scdRepositoryImpl[T]) FindLatestWithFilter(ctx context.Context, filter map[string]interface{}) ([]T, error) {
	var results []T
	err := r.db.GetSlaveDB(ctx).
		Where("is_latest = ?", true).
		Where(filter).
		Find(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find latest records with filter: %w", err)
	}

	return results, nil
}

// FindVersionsForID returns all versions of a record by ID
func (r *scdRepositoryImpl[T]) FindVersionsForID(ctx context.Context, id string) ([]T, error) {
	var results []T
	err := r.db.GetSlaveDB(ctx).
		Where("id = ?", id).
		Order("version ASC").
		Find(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find versions for ID: %w", err)
	}

	return results, nil
}

// Create creates a new record with version 1
func (r *scdRepositoryImpl[T]) Create(ctx context.Context, record *T) error {
	(*record).SetVersion(1)
	(*record).SetID(uuid.New().String())
	(*record).SetIsLatest(true)
	(*record).SetUID(uuid.New().String())

	err := r.db.GetMasterDB(ctx).Create(record).Error
	if err != nil {
		return fmt.Errorf("failed to create record: %w", err)
	}

	return nil
}

// Update creates a new version of an existing record
func (r *scdRepositoryImpl[T]) Update(ctx context.Context, id string, record *T) error {
	// Find the latest version
	latestRecord, err := r.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find latest version: %w", err)
	}

	if latestRecord == nil {
		return errors.New("record not found")
	}

	// Set the new version
	(*record).SetVersion((*latestRecord).GetVersion() + 1)
	(*record).SetIsLatest(true)
	(*record).SetUID(uuid.New().String())
	(*record).SetID((*latestRecord).GetID())

	// Execute operations within a transaction
	err = r.db.GetMasterDB(ctx).Transaction(func(tx *gorm.DB) error {
		// Create a new record with incremented version
		if err := tx.Create(record).Error; err != nil {
			return fmt.Errorf("failed to update record: %w", err)
		}

		// Update the latest flag for the old version
		if err := tx.Model(r.modelType).
			Where("uid = ? AND is_latest = ?", (*latestRecord).GetUID(), true).
			Update("is_latest", false).Error; err != nil {
			return fmt.Errorf("failed to update latest flag: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// CustomQuery executes a custom query with SCD handling
func (r *scdRepositoryImpl[T]) CustomQuery(ctx context.Context, queryBuilder func(*gorm.DB) *gorm.DB) ([]T, error) {
	var results []T

	// Get master DB
	db := r.db.GetMasterDB(ctx)

	stmt := &gorm.Statement{DB: db}
	err := stmt.Parse(r.modelType)
	if err != nil {
		return nil, fmt.Errorf("failed to parse model type: %w", err)
	}
	tableName := stmt.Table

	// Start with a base query that includes only the latest versions
	// Use qualified column name to avoid ambiguity
	baseQuery := db.Model(r.modelType).
		Where(fmt.Sprintf("%s.is_latest = ?", tableName), true)

	// Apply the user's custom query function to the base query
	customQuery := queryBuilder(baseQuery)

	// Execute the combined query
	err = customQuery.Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to execute custom query: %w", err)
	}

	return results, nil
}
