package model

import (
	"testing"
)

func TestNewEnvironment(t *testing.T) {
	tests := []struct {
		input    string
		expected Environment
	}{
		{"development", Development},
		{"dev", Dev},
		{"staging", Staging},
		{"production", Production},
		{"prod", Prod},
		{"local", Local},
		{"invalid", Invalid},
		{"unknown", Invalid},
	}

	for _, test := range tests {
		result := NewEnvironment(test.input)
		if result != test.expected {
			t.Errorf("NewEnvironment(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}
func TestEnvironmentString(t *testing.T) {
	tests := []struct {
		env      Environment
		expected string
	}{
		{Development, "development"},
		{Dev, "dev"},
		{Staging, "staging"},
		{Production, "production"},
		{Prod, "prod"},
		{Local, "local"},
		{Invalid, "invalid"},
	}

	for _, test := range tests {
		result := test.env.String()
		if result != test.expected {
			t.Errorf("Environment(%q).String() = %v; want %v", test.env, result, test.expected)
		}
	}
}
func TestEnvFunctions(t *testing.T) {
	tests := []struct {
		env      Environment
		envFn    func() bool
		expected bool
		envStr   func() string
	}{
		{Development, Development.IsDevelopment, true, Development.String},
		{Dev, Dev.IsDevelopment, true, Dev.String},
		{Staging, Staging.IsStaging, true, Staging.String},
		{Production, Production.IsProduction, true, Production.String},
		{Prod, Prod.IsProduction, true, Prod.String},
		{Local, Local.IsLocal, true, Local.String},
	}

	for _, test := range tests {
		result := test.envFn()
		if result != test.expected {
			t.Errorf("Environment(%q).IsDevelopment() = %v; want %v", test.env, result, test.expected)
		}

		resultStr := test.envStr()
		if resultStr != test.env.String() {
			t.Errorf("Environment(%q).String() = %v; want %v", test.env, resultStr, test.env.String())
		}
	}
}
