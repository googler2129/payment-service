package config

import (
	"context"
	"time"

	"github.com/mercor/payment-service/constants"
	"github.com/mercor/payment-service/pkg/config/model"
	"github.com/mercor/payment-service/pkg/log"

	"github.com/spf13/cast"
)

func getConfigFromContext(ctx context.Context, key string) *model.Config {
	contextConfig := ctx.Value(constants.Config)
	if contextConfig == nil {
		tempApp := getApplication()
		if tempApp == nil {
			return nil
		}
		return tempApp.observer.GetConfig()
	}

	conf, ok := contextConfig.(*model.Config)
	if !ok {
		return nil
	}

	return conf
}

func Get(ctx context.Context, key string) interface{} {
	conf := getConfigFromContext(ctx, key)

	if conf == nil {
		log.Errorf("error while parsing configuration for key %s", key)
		return nil
	}

	val, ok := conf.GetValueForKey(key)
	if !ok {
		log.Errorf("error while parsing configuration for key %s", key)
		return nil
	}

	return val
}

func GetBool(ctx context.Context, key string) bool {
	val, err := cast.ToBoolE(Get(ctx, key))
	if err != nil {
		log.Errorf("error while parsing configuration for key %s", key)
	}
	return val
}

func GetString(ctx context.Context, key string) string {
	val, err := cast.ToStringE(Get(ctx, key))
	if err != nil {
		log.Errorf("error while parsing configuration for key %s", key)
	}
	return val
}

func GetInt(ctx context.Context, key string) int {
	val, err := cast.ToIntE(Get(ctx, key))
	if err != nil {
		log.Errorf("error while parsing configuration for key %s", key)
	}
	return val
}

func GetInt64(ctx context.Context, key string) int64 {
	val, err := cast.ToInt64E(Get(ctx, key))
	if err != nil {
		log.Errorf("error while parsing configuration for key %s", key)
	}
	return val
}

func GetUint(ctx context.Context, key string) uint {
	val, err := cast.ToUintE(Get(ctx, key))
	if err != nil {
		log.Errorf("error while parsing configuration for key %s", key)
	}
	return val
}

func GetUint64(ctx context.Context, key string) uint64 {
	val, err := cast.ToUint64E(Get(ctx, key))
	if err != nil {
		log.Errorf("error while parsing configuration for key %s", key)
	}
	return val
}

func GetFloat32(ctx context.Context, key string) float32 {
	val, err := cast.ToFloat32E(Get(ctx, key))
	if err != nil {
		log.Errorf("error while parsing configuration for key %s", key)
	}
	return val
}

func GetFloat64(ctx context.Context, key string) float64 {
	val, err := cast.ToFloat64E(Get(ctx, key))
	if err != nil {
		log.Errorf("error while parsing configuration for key %s", key)
	}
	return val
}

func GetTime(ctx context.Context, key string) time.Time {
	val, err := cast.ToTimeE(Get(ctx, key))
	if err != nil {
		log.Errorf("error while parsing configuration for key %s", key)
	}
	return val
}

func GetDuration(ctx context.Context, key string) time.Duration {
	val, err := cast.ToDurationE(Get(ctx, key))
	if err != nil {
		log.Errorf("error while parsing configuration for key %s", key)
	}
	return val
}

func GetSlice(ctx context.Context, key string) []interface{} {
	val, err := cast.ToSliceE(Get(ctx, key))
	if err != nil {
		log.Errorf("error while parsing configuration for key %s", key)
	}
	return val
}

func GetBoolSlice(ctx context.Context, key string) []bool {
	val, err := cast.ToBoolSliceE(Get(ctx, key))
	if err != nil {
		log.Errorf("error while parsing configuration for key %s", key)
	}
	return val
}

func GetIntSlice(ctx context.Context, key string) []int {
	val, err := cast.ToIntSliceE(Get(ctx, key))
	if err != nil {
		log.Errorf("error while parsing configuration for key %s", key)
	}
	return val
}

func GetStringSlice(ctx context.Context, key string) []string {
	val, err := cast.ToStringSliceE(Get(ctx, key))
	if err != nil {
		log.Errorf("error while parsing configuration for key %s", key)
	}
	return val
}

func IsDevelopment(ctx context.Context) bool {
	conf := getConfigFromContext(ctx, constants.DeployedEnv)

	if conf == nil {
		log.Errorf("error while parsing configuration for key %s", constants.DeployedEnv)
		return false
	}

	return conf.GetEnvironment().IsDevelopment()
}

func IsStaging(ctx context.Context) bool {
	conf := getConfigFromContext(ctx, constants.DeployedEnv)

	if conf == nil {
		log.Errorf("error while parsing configuration for key %s", constants.DeployedEnv)
		return false
	}

	return conf.GetEnvironment().IsStaging()
}

func IsProduction(ctx context.Context) bool {
	conf := getConfigFromContext(ctx, constants.DeployedEnv)

	if conf == nil {
		log.Errorf("error while parsing configuration for key %s", constants.DeployedEnv)
		return false
	}

	return conf.GetEnvironment().IsProduction()
}

func IsLocal(ctx context.Context) bool {
	conf := getConfigFromContext(ctx, constants.DeployedEnv)

	if conf == nil {
		log.Errorf("error while parsing configuration for key %s", constants.DeployedEnv)
		return false
	}

	return conf.GetEnvironment().IsLocal()
}

func SwaggerEnabled(ctx context.Context) bool {
	return !IsProduction(ctx)
}

func ProfilingEnabled(ctx context.Context) bool {
	profilingEnabled := GetBool(ctx, constants.ProfilingEnabledKey)
	return !IsProduction(ctx) && profilingEnabled
}
