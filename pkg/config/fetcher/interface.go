package fetcher

import (
	"context"

	"github.com/mercor/payment-service/pkg/config/model"
)

type Fetcher interface {
	GetConfig(ctx context.Context) (*model.Config, error)
}
