package config

import (
	"github.com/hoitek/Maja-Service/internal/cycle/config"
)

// CycleConfig is a global variable for the cycle domain config
var CycleConfig config.ConfigType

// LoadCycleConfig loads the cycle domain config
func LoadCycleConfig() config.ConfigType {
	CycleConfig = config.ConfigType{
		Debug:               AppConfig.Debug,
		Environment:         AppConfig.Environment,
		ApiPrefix:           AppConfig.ApiPrefix,
		ApiVersion1:         AppConfig.ApiVersion1,
		ApiVersion2:         AppConfig.ApiVersion2,
		DatabasePort:        AppConfig.DatabasePort,
		DatabaseName:        AppConfig.DatabaseName,
		DatabaseHost:        AppConfig.DatabaseHost,
		DatabaseUser:        AppConfig.DatabaseUser,
		DatabasePassword:    AppConfig.DatabasePassword,
		DatabaseSslMode:     AppConfig.DatabaseSslMode,
		DatabaseTimeZone:    AppConfig.DatabaseTimeZone,
		DatabaseMongoDBHost: AppConfig.DatabaseMongoDBHost,
		DatabaseMongoDBPort: AppConfig.DatabaseMongoDBPort,
		DatabaseMongoDBName: AppConfig.DatabaseMongoDBName,
		DatabaseMongoDBUser: AppConfig.DatabaseMongoDBUser,
		DatabaseMongoDBPass: AppConfig.DatabaseMongoDBPass,
		RabbitMQHost:        AppConfig.RabbitMQHost,
		RabbitMQPort:        AppConfig.RabbitMQPort,
		RabbitMQUser:        AppConfig.RabbitMQUser,
		RabbitMQPassword:    AppConfig.RabbitMQPassword,
		MaxUploadSizeLimit:  AppConfig.MaxUploadSizeLimit,
	}
	return CycleConfig
}
