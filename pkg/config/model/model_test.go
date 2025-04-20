package model

import (
	"strings"
	"testing"

	"github.com/mercor/payment-service/constants"
)

func TestNewConfig(t *testing.T) {
	t.Run("should create config with environment set", func(t *testing.T) {
		data := map[string]interface{}{
			strings.ToLower(constants.DeployedEnv): "production",
		}
		config := NewConfig(data)

		if config == nil {
			t.Fatalf("expected config to be non-nil")
		}

		if config.environment != Environment("production") {
			t.Errorf("expected environment to be 'production', got '%v'", config.environment)
		}
	})

	t.Run("should create config without environment set", func(t *testing.T) {
		data := map[string]interface{}{}
		config := NewConfig(data)

		if config == nil {
			t.Fatalf("expected config to be non-nil")
		}

		if config.GetEnvironment() != "" {
			t.Errorf("expected environment to be empty, got '%v'", config.environment)
		}
	})

	t.Run("should create config with invalid environment", func(t *testing.T) {
		data := map[string]interface{}{
			strings.ToLower(constants.DeployedEnv): "invalid",
		}
		config := NewConfig(data)

		if config == nil {
			t.Fatalf("expected config to be non-nil")
		}

		if config.GetEnvironment() != "" {
			t.Errorf("expected environment to be empty, got '%v'", config.environment)
		}
	})

	t.Run("should run with nil config", func(t *testing.T) {
		config := NewConfig(nil)

		config = nil

		if config != nil {
			t.Fatalf("expected config to be nil")
		}

		d, ok := config.GetValueForKey("key")

		if d != nil || ok {
			t.Fatalf("expected value to be nil and ok to be false")
		}
	})

	t.Run("should run with non string environment", func(t *testing.T) {
		data := map[string]interface{}{
			strings.ToLower(constants.DeployedEnv): 123,
		}
		config := NewConfig(data)

		if config == nil {
			t.Fatalf("expected config to be non-nil")
		}

		if config.GetEnvironment() != "" {
			t.Errorf("expected environment to be empty, got '%v'", config.environment)
		}
	})
}
