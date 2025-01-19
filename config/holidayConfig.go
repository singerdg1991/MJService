package config

import (
	"github.com/hoitek/Maja-Service/internal/staffclub/holiday/config"
)

// HolidayConfig is a global variable for the holiday domain config
var HolidayConfig config.ConfigType

// LoadHolidayConfig loads the holiday domain config
func LoadHolidayConfig() config.ConfigType {
	HolidayConfig = config.ConfigType{
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
	return HolidayConfig
}
