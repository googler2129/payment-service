package router

import (
	"context"

	"github.com/mercor/payment-service/pkg/config"
	"github.com/mercor/payment-service/pkg/http"
	"github.com/mercor/payment-service/pkg/log"
)

const (
	DefaultPerPageLimit = 100
)

func Initialize(ctx context.Context, s *http.Server) (err error) {
	//Middleware for adding config to ctx
	s.Engine.Use(config.Middleware())

	s.Engine.Use(log.RequestLogMiddleware(log.MiddlewareOptions{
		Format:      config.GetString(ctx, "log.format"),
		Level:       config.GetString(ctx, "log.level"),
		LogRequest:  config.GetBool(ctx, "log.request"),
		LogResponse: config.GetBool(ctx, "log.response"),
	}))

	err = PublicRoutes(ctx, s)
	if err != nil {
		return
	}

	return
}
