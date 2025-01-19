package config

import (
	"github.com/hoitek/Maja-Service/internal/limitation/config"
)

// LimitationConfig is a global variable for the limitation domain config
var LimitationConfig config.ConfigType

// LoadLimitationConfig loads the limitation domain config
func LoadLimitationConfig() config.ConfigType {
	LimitationConfig = config.ConfigType{
		Debug:            AppConfig.Debug,
		Environment:      AppConfig.Environment,
		ApiPrefix:        AppConfig.ApiPrefix,
		ApiVersion1:      AppConfig.ApiVersion1,
		ApiVersion2:      AppConfig.ApiVersion2,
		DatabasePort:     AppConfig.DatabasePort,
		DatabaseName:     AppConfig.DatabaseName,
		DatabaseHost:     AppConfig.DatabaseHost,
		DatabaseUser:     AppConfig.DatabaseUser,
		DatabasePassword: AppConfig.DatabasePassword,
		DatabaseSslMode:  AppConfig.DatabaseSslMode,
		DatabaseTimeZone: AppConfig.DatabaseTimeZone,
	}
	return LimitationConfig
}
