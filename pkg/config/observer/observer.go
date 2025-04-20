package observer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mercor/payment-service/pkg/config/fetcher"
	"github.com/mercor/payment-service/pkg/config/model"
	"github.com/mercor/payment-service/pkg/log"
)

type Observer struct {
	config *model.Config
	mu     sync.RWMutex
}

func NewObserver(ctx context.Context, fetcher fetcher.Fetcher, pollInterval time.Duration) (*Observer, error) {
	observer := &Observer{}

	// Initialize the observer by fetching the initial configuration and setting up polling
	if err := observer.startPolling(ctx, fetcher, pollInterval); err != nil {
		return nil, err
	}

	return observer, nil
}

// GetConfig safely returns the current configuration using a read lock for concurrency
func (o *Observer) GetConfig() *model.Config {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.config
}

// startPolling fetches the initial configuration and sets up periodic updates
func (o *Observer) startPolling(ctx context.Context, fetcher fetcher.Fetcher, pollInterval time.Duration) error {
	initialConfig, err := fetcher.GetConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch initial configuration: %w", err)
	}

	o.updateConfig(initialConfig)

	// Start a goroutine to periodically fetch and update the configuration
	go o.pollUpdates(ctx, fetcher, pollInterval)

	return nil
}

// pollUpdates continuously fetches and updates the configuration at the specified interval
func (o *Observer) pollUpdates(ctx context.Context, fetcher fetcher.Fetcher, pollInterval time.Duration) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c, err := fetcher.GetConfig(ctx)
			if err != nil {
				log.Error("Failed to fetch configuration: ", err)
				continue
			}
			o.updateConfig(c)

		case <-ctx.Done():
			log.Info("Stopping configuration polling")
			return
		}
	}
}

// updateConfig safely updates the configuration using a write lock for concurrency
func (o *Observer) updateConfig(newConfig *model.Config) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.config = newConfig
}
