//go:build wireinject
// +build wireinject

package contractor

import (
	"context"

	"github.com/google/wire"
	"github.com/mercor/payment-service/pkg/db/sql/postgres"
)

func Wire(ctx context.Context, db *postgres.DbCluster) (*Controller, error) {
	panic(wire.Build(ProviderSet))
}
