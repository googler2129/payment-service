package model

import (
	"strings"

	"github.com/mercor/payment-service/constants"
)

type Config struct {
	data        map[string]interface{}
	environment Environment
}

func (conf *Config) GetValueForKey(key string) (interface{}, bool) {
	if conf == nil {
		return nil, false
	}
	val, ok := conf.data[strings.ToLower(key)]
	return val, ok
}

func (conf *Config) SetEnvironment(envVal interface{}) {
	env, ok := envVal.(string)
	if !ok {
		return
	}
	environment := Environment(env)

	if environment.IsValid() {
		conf.environment = environment
	}
}

func (conf *Config) GetEnvironment() Environment {
	return conf.environment
}

func NewConfig(data map[string]interface{}) *Config {
	config := Config{data: data}
	envVal, ok := config.GetValueForKey(constants.DeployedEnv)
	if ok {
		config.SetEnvironment(envVal)
	}
	return &config
}
