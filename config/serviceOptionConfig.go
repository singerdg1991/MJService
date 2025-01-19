package config

import (
	"github.com/hoitek/Maja-Service/internal/serviceoption/config"
)

// ServiceOptionConfig is a global variable for the ServiceOption domain config
var ServiceOptionConfig config.ConfigType

// LoadServiceOptionConfig loads the ServiceOption domain config
func LoadServiceOptionConfig() config.ConfigType {
	ServiceOptionConfig = config.ConfigType{
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
	return ServiceOptionConfig
}
