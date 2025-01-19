package config

import (
	"github.com/hoitek/Maja-Service/internal/city/config"
)

// CityConfig is a global variable for the city domain config
var CityConfig config.ConfigType

// LoadCityConfig loads the city domain config
func LoadCityConfig() config.ConfigType {
	CityConfig = config.ConfigType{
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
	return CityConfig
}
