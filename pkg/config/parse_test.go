package config

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mercor/payment-service/constants"
	"github.com/mercor/payment-service/pkg/config/model"
	"github.com/stretchr/testify/assert"
)

// Setup config to constants.LocalFreeFormPath for testing
// Add sample yaml file to the path
func setupConfig() error {
	// sample yaml config
	config := `
testKey: testValue
testBoolKey: true
testIntKey: 123
testInt64Key: 123456789
testUintKey: 123
testUint64Key: 123456789
testFloat32Key: 123.45
testFloat64Key: 123.456
testTimeKey: 2023-10-10T10:10:10Z
testDurationKey: 10s
testSliceKey:
  - 1
  - two
  - 3.0
testBoolSliceKey:
  - true
  - false
  - true
testIntSliceKey:
  - 1
  - 2
  - 3
testStringSliceKey:
  - one
  - two
  - three
`

	// create directory if not exists
	err := os.MkdirAll(strings.Join(strings.Split(constants.LocalFreeFormPath, "/")[:2], "/"), 0755)
	if err != nil {
		fmt.Println("error creating directory")
		return err
	}
	// write config to file
	err = os.WriteFile(constants.LocalFreeFormPath, []byte(config), 0644)
	if err != nil {
		fmt.Println("error writing config file")
		return err
	}
	os.Setenv("CONFIG_SOURCE", "local")
	return nil
}

// delete config file
func teardownConfig() error {
	err := os.Remove(constants.LocalFreeFormPath)
	if err != nil {
		fmt.Println("error deleting config file")
		return err
	}
	return nil
}

func setupSuite(tb testing.TB, withConfig bool) func(tb testing.TB) {
	fmt.Println("setup suite")
	if withConfig {
		err := setupConfig()
		if err != nil {
			tb.Fatal(err)
		}
	}
	Init(time.Second * 10)
	return func(tb testing.TB) {
		fmt.Println("teardown suite")
		if withConfig {
			err := teardownConfig()
			if err != nil {
				tb.Fatal(err)
			}
		}
	}
}

func TestGetConfigFromContext(t *testing.T) {
	teardownSuite := setupSuite(t, false)
	defer teardownSuite(t)
	tests := []struct {
		name        string
		ctx         context.Context
		setupConfig bool
		wantNil     bool
	}{
		{
			name:    "context with config",
			ctx:     context.WithValue(context.Background(), constants.Config, &model.Config{}),
			wantNil: false,
		},
		{
			name:    "context without config",
			ctx:     context.Background(),
			wantNil: true,
		},
		{
			name:    "context with invalid config type",
			ctx:     context.WithValue(context.Background(), constants.Config, "invalid"),
			wantNil: true,
		},
		{
			name:        "context with config from application",
			ctx:         context.WithValue(context.Background(), constants.Config, nil),
			wantNil:     true,
			setupConfig: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupConfig {
				if err := setupConfig(); err != nil {
					t.Fatal(err)
				}
				defer teardownConfig()
			}
			got := getConfigFromContext(tt.ctx, "testKey")
			if tt.wantNil {
				assert.Nil(t, got)
			} else {
				assert.NotNil(t, got)
			}
		})
	}
}
func TestGet(t *testing.T) {
	teardownSuite := setupSuite(t, false)
	defer teardownSuite(t)
	tests := []struct {
		name    string
		ctx     context.Context
		key     string
		want    interface{}
		wantNil bool
	}{
		{
			name: "valid key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testkey": "testValue"}),
			),
			key:     "testkey",
			want:    "testValue",
			wantNil: false,
		},
		{
			name: "invalid key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testKey": "testValue"}),
			),
			key:     "invalidKey",
			want:    nil,
			wantNil: true,
		},
		{
			name:    "no config in context",
			ctx:     context.Background(),
			key:     "testKey",
			want:    nil,
			wantNil: true,
		},
		{
			name:    "invalid config type in context",
			ctx:     context.WithValue(context.Background(), constants.Config, "invalid"),
			key:     "testKey",
			want:    nil,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Get(tt.ctx, tt.key)
			if tt.wantNil {
				assert.Nil(t, got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
func TestGetBool(t *testing.T) {
	teardownSuite := setupSuite(t, false)
	defer teardownSuite(t)
	tests := []struct {
		name    string
		ctx     context.Context
		key     string
		want    bool
		wantErr bool
	}{
		{
			name: "valid bool key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testboolkey": true}),
			),
			key:     "testBoolKey",
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid bool key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testboolkey": "notBool"}),
			),
			key:     "testBoolKey",
			want:    false,
			wantErr: true,
		},
		{
			name:    "no config in context",
			ctx:     context.Background(),
			key:     "testBoolKey",
			want:    false,
			wantErr: true,
		},
		{
			name:    "invalid config type in context",
			ctx:     context.WithValue(context.Background(), constants.Config, "invalid"),
			key:     "testBoolKey",
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetBool(tt.ctx, tt.key)
			if tt.wantErr {
				assert.False(t, got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
func TestGetString(t *testing.T) {
	teardownSuite := setupSuite(t, false)
	defer teardownSuite(t)
	tests := []struct {
		name    string
		ctx     context.Context
		key     string
		want    string
		wantErr bool
	}{
		{
			name: "valid string key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"teststringkey": "testValue"}),
			),
			key:     "testStringKey",
			want:    "testValue",
			wantErr: false,
		},
		{
			name: "invalid string key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"teststringkey": 123}),
			),
			key:     "testStringKey",
			want:    "123",
			wantErr: false,
		},
		{
			name: "invalid data type in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"teststringkeyerror": struct{ name string }{name: "test"}})),
			key:     "teststringkeyerror",
			want:    "",
			wantErr: true,
		},
		{
			name:    "no config in context",
			ctx:     context.Background(),
			key:     "testStringKey",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid config type in context",
			ctx:     context.WithValue(context.Background(), constants.Config, "invalid"),
			key:     "testStringKey",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetString(tt.ctx, tt.key)
			if tt.wantErr {
				assert.Empty(t, got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
func TestGetInt(t *testing.T) {
	teardownSuite := setupSuite(t, false)
	defer teardownSuite(t)
	tests := []struct {
		name    string
		ctx     context.Context
		key     string
		want    int
		wantErr bool
	}{
		{
			name: "valid int key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testintkey": 123}),
			),
			key:     "testIntKey",
			want:    123,
			wantErr: false,
		},
		{
			name: "invalid int key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testintkey": "notInt"}),
			),
			key:     "testIntKey",
			want:    0,
			wantErr: true,
		},
		{
			name:    "no config in context",
			ctx:     context.Background(),
			key:     "testIntKey",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid config type in context",
			ctx:     context.WithValue(context.Background(), constants.Config, "invalid"),
			key:     "testIntKey",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetInt(tt.ctx, tt.key)
			if tt.wantErr {
				assert.Equal(t, 0, got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
func TestGetInt64(t *testing.T) {
	teardownSuite := setupSuite(t, false)
	defer teardownSuite(t)
	tests := []struct {
		name    string
		ctx     context.Context
		key     string
		want    int64
		wantErr bool
	}{
		{
			name: "valid int64 key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testint64key": int64(123456789)}),
			),
			key:     "testInt64Key",
			want:    int64(123456789),
			wantErr: false,
		},
		{
			name: "invalid int64 key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testint64key": "notInt64"}),
			),
			key:     "testInt64Key",
			want:    0,
			wantErr: true,
		},
		{
			name:    "no config in context",
			ctx:     context.Background(),
			key:     "testInt64Key",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid config type in context",
			ctx:     context.WithValue(context.Background(), constants.Config, "invalid"),
			key:     "testInt64Key",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetInt64(tt.ctx, tt.key)
			if tt.wantErr {
				assert.Equal(t, int64(0), got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
func TestGetUint(t *testing.T) {
	teardownSuite := setupSuite(t, false)
	defer teardownSuite(t)
	tests := []struct {
		name    string
		ctx     context.Context
		key     string
		want    uint
		wantErr bool
	}{
		{
			name: "valid uint key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testuintkey": uint(123)}),
			),
			key:     "testUintKey",
			want:    uint(123),
			wantErr: false,
		},
		{
			name: "invalid uint key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testuintkey": "notUint"}),
			),
			key:     "testUintKey",
			want:    0,
			wantErr: true,
		},
		{
			name:    "no config in context",
			ctx:     context.Background(),
			key:     "testUintKey",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid config type in context",
			ctx:     context.WithValue(context.Background(), constants.Config, "invalid"),
			key:     "testUintKey",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetUint(tt.ctx, tt.key)
			if tt.wantErr {
				assert.Equal(t, uint(0), got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
func TestGetUint64(t *testing.T) {
	teardownSuite := setupSuite(t, false)
	defer teardownSuite(t)
	tests := []struct {
		name    string
		ctx     context.Context
		key     string
		want    uint64
		wantErr bool
	}{
		{
			name: "valid uint64 key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testuint64key": uint64(123456789)}),
			),
			key:     "testUint64Key",
			want:    uint64(123456789),
			wantErr: false,
		},
		{
			name: "invalid uint64 key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testuint64key": "notUint64"}),
			),
			key:     "testUint64Key",
			want:    0,
			wantErr: true,
		},
		{
			name:    "no config in context",
			ctx:     context.Background(),
			key:     "testUint64Key",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid config type in context",
			ctx:     context.WithValue(context.Background(), constants.Config, "invalid"),
			key:     "testUint64Key",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetUint64(tt.ctx, tt.key)
			if tt.wantErr {
				assert.Equal(t, uint64(0), got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
func TestGetFloat32(t *testing.T) {
	teardownSuite := setupSuite(t, false)
	defer teardownSuite(t)
	tests := []struct {
		name    string
		ctx     context.Context
		key     string
		want    float32
		wantErr bool
	}{
		{
			name: "valid float32 key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testfloat32key": float32(123.45)}),
			),
			key:     "testFloat32Key",
			want:    float32(123.45),
			wantErr: false,
		},
		{
			name: "invalid float32 key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testfloat32key": "notFloat32"}),
			),
			key:     "testFloat32Key",
			want:    0,
			wantErr: true,
		},
		{
			name:    "no config in context",
			ctx:     context.Background(),
			key:     "testFloat32Key",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid config type in context",
			ctx:     context.WithValue(context.Background(), constants.Config, "invalid"),
			key:     "testFloat32Key",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFloat32(tt.ctx, tt.key)
			if tt.wantErr {
				assert.Equal(t, float32(0), got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
func TestGetFloat64(t *testing.T) {
	teardownSuite := setupSuite(t, false)
	defer teardownSuite(t)
	tests := []struct {
		name    string
		ctx     context.Context
		key     string
		want    float64
		wantErr bool
	}{
		{
			name: "valid float64 key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testfloat64key": float64(123.456)}),
			),
			key:     "testFloat64Key",
			want:    float64(123.456),
			wantErr: false,
		},
		{
			name: "invalid float64 key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testfloat64key": "notFloat64"}),
			),
			key:     "testFloat64Key",
			want:    0,
			wantErr: true,
		},
		{
			name:    "no config in context",
			ctx:     context.Background(),
			key:     "testFloat64Key",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid config type in context",
			ctx:     context.WithValue(context.Background(), constants.Config, "invalid"),
			key:     "testFloat64Key",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFloat64(tt.ctx, tt.key)
			if tt.wantErr {
				assert.Equal(t, float64(0), got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
func TestGetTime(t *testing.T) {
	teardownSuite := setupSuite(t, false)
	defer teardownSuite(t)
	tests := []struct {
		name    string
		ctx     context.Context
		key     string
		want    time.Time
		wantErr bool
	}{
		{
			name: "valid time key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testtimekey": "2023-10-10T10:10:10Z"}),
			),
			key:     "testTimeKey",
			want:    time.Date(2023, 10, 10, 10, 10, 10, 0, time.UTC),
			wantErr: false,
		},
		{
			name: "invalid time key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testtimekey": "notTime"}),
			),
			key:     "testTimeKey",
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "no config in context",
			ctx:     context.Background(),
			key:     "testTimeKey",
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "invalid config type in context",
			ctx:     context.WithValue(context.Background(), constants.Config, "invalid"),
			key:     "testTimeKey",
			want:    time.Time{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTime(tt.ctx, tt.key)
			if tt.wantErr {
				assert.True(t, got.IsZero())
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
func TestGetDuration(t *testing.T) {
	teardownSuite := setupSuite(t, false)
	defer teardownSuite(t)
	tests := []struct {
		name    string
		ctx     context.Context
		key     string
		want    time.Duration
		wantErr bool
	}{
		{
			name: "valid duration key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testdurationkey": "10s"}),
			),
			key:     "testDurationKey",
			want:    10 * time.Second,
			wantErr: false,
		},
		{
			name: "invalid duration key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testdurationkey": "notDuration"}),
			),
			key:     "testDurationKey",
			want:    0,
			wantErr: true,
		},
		{
			name:    "no config in context",
			ctx:     context.Background(),
			key:     "testDurationKey",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid config type in context",
			ctx:     context.WithValue(context.Background(), constants.Config, "invalid"),
			key:     "testDurationKey",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDuration(tt.ctx, tt.key)
			if tt.wantErr {
				assert.Equal(t, time.Duration(0), got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
func TestGetSlice(t *testing.T) {
	teardownSuite := setupSuite(t, false)
	defer teardownSuite(t)
	tests := []struct {
		name    string
		ctx     context.Context
		key     string
		want    []interface{}
		wantErr bool
	}{
		{
			name: "valid slice key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testslicekey": []interface{}{1, "two", 3.0}}),
			),
			key:     "testSliceKey",
			want:    []interface{}{1, "two", 3.0},
			wantErr: false,
		},
		{
			name: "invalid slice key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testslicekey": "notSlice"}),
			),
			key:     "testSliceKey",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "no config in context",
			ctx:     context.Background(),
			key:     "testSliceKey",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid config type in context",
			ctx:     context.WithValue(context.Background(), constants.Config, "invalid"),
			key:     "testSliceKey",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSlice(tt.ctx, tt.key)
			if tt.wantErr {
				assert.Nil(t, got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
func TestGetBoolSlice(t *testing.T) {
	teardownSuite := setupSuite(t, false)
	defer teardownSuite(t)
	tests := []struct {
		name    string
		ctx     context.Context
		key     string
		want    []bool
		wantErr bool
	}{
		{
			name: "valid bool slice key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testboolslicekey": []bool{true, false, true}}),
			),
			key:     "testBoolSliceKey",
			want:    []bool{true, false, true},
			wantErr: false,
		},
		{
			name: "invalid bool slice key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testboolslicekey": "notBoolSlice"}),
			),
			key:     "testBoolSliceKey",
			want:    []bool{},
			wantErr: false,
		},
		{
			name:    "no config in context",
			ctx:     context.Background(),
			key:     "testBoolSliceKey",
			want:    []bool{},
			wantErr: false,
		},
		{
			name:    "invalid config type in context",
			ctx:     context.WithValue(context.Background(), constants.Config, "invalid"),
			key:     "testBoolSliceKey",
			want:    []bool{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetBoolSlice(tt.ctx, tt.key)
			if tt.wantErr {
				assert.Nil(t, got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
func TestGetIntSlice(t *testing.T) {
	teardownSuite := setupSuite(t, false)
	defer teardownSuite(t)
	tests := []struct {
		name    string
		ctx     context.Context
		key     string
		want    []int
		wantErr bool
	}{
		{
			name: "valid int slice key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testintslicekey": []int{1, 2, 3}}),
			),
			key:     "testIntSliceKey",
			want:    []int{1, 2, 3},
			wantErr: false,
		},
		{
			name: "invalid int slice key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"testintslicekey": "notIntSlice"}),
			),
			key:     "testIntSliceKey",
			want:    []int{},
			wantErr: false,
		},
		{
			name:    "no config in context",
			ctx:     context.Background(),
			key:     "testIntSliceKey",
			want:    []int{},
			wantErr: false,
		},
		{
			name:    "invalid config type in context",
			ctx:     context.WithValue(context.Background(), constants.Config, "invalid"),
			key:     "testIntSliceKey",
			want:    []int{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetIntSlice(tt.ctx, tt.key)
			if tt.wantErr {
				assert.Nil(t, got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
func TestGetStringSlice(t *testing.T) {
	teardownSuite := setupSuite(t, false)
	defer teardownSuite(t)
	tests := []struct {
		name    string
		ctx     context.Context
		key     string
		want    []string
		wantErr bool
	}{
		{
			name: "valid string slice key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"teststringslicekey": []string{"one", "two", "three"}}),
			),
			key:     "testStringSliceKey",
			want:    []string{"one", "two", "three"},
			wantErr: false,
		},
		{
			name: "invalid string slice key in context config",
			ctx: context.WithValue(
				context.Background(),
				constants.Config,
				model.NewConfig(map[string]interface{}{"teststringslicekey": "notStringSlice"}),
			),
			key:     "testStringSliceKey",
			want:    []string{"notStringSlice"},
			wantErr: false,
		},
		{
			name:    "no config in context",
			ctx:     context.Background(),
			key:     "testStringSliceKey",
			want:    nil,
			wantErr: false,
		},
		{
			name:    "invalid config type in context",
			ctx:     context.WithValue(context.Background(), constants.Config, "invalid"),
			key:     "testStringSliceKey",
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetStringSlice(tt.ctx, tt.key)
			if tt.wantErr {
				assert.Nil(t, got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
