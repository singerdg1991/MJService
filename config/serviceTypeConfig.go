package config

import (
	"github.com/hoitek/Maja-Service/internal/servicetype/config"
)

// ServiceTypeConfig is a global variable for the serviceType domain config
var ServiceTypeConfig config.ConfigType

// LoadServiceTypeConfig loads the serviceType domain config
func LoadServiceTypeConfig() config.ConfigType {
	ServiceTypeConfig = config.ConfigType{
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
	return ServiceTypeConfig
}
