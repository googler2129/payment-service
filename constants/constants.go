package constants

import "time"

const (
	EventualConsistency  = "eventual"
	StrongConsistency    = "strong"
	LocalConsistency     = "local"
	LocalFreeFormPath    = "./configs/config.yaml"
	CloudFreeFormPath    = "./config/cloud.yaml"
	CloudProfileName     = "default"
	CloudApplicationName = "payment-service"

	HeaderXMercorRequestID = "X-Mercor-Request-ID"
	Env                    = "env"
	DeployedEnv            = "deployed_env"
	DefaultTimeout         = 10 * time.Second

	LocalStorePath        = "./config/local.yaml"
	CloudStorePath        = "./config/cloud.yaml"
	RemoteStorePath       = "./config/remote.yaml"
	ConfigStorePath       = "./config/config.yaml"
	ProfilingEnabledKey   = "profiling_enabled"
	ProfilingEnabledValue = "true"
	Config                = "config"
	LocalStore            = "local"
	RemoteFreeformProfile = "remote_freeform"
	LocalSource           = "local"

	Consistency  = "consistency"
	DBPreference = "db_preference"
	SlaveDB      = "slave_db"

	Authorization = "Authorization"
	Bearer        = "Bearer"
	AccessToken   = "access_token"

	UserDetails = "user_details"
)
