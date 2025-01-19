package config

import (
	"github.com/hoitek/Maja-Service/internal/shifttype/config"
)

// ShiftTypeConfig is a global variable for the shiftType domain config
var ShiftTypeConfig config.ConfigType

// LoadShiftTypeConfig loads the shiftType domain config
func LoadShiftTypeConfig() config.ConfigType {
	ShiftTypeConfig = config.ConfigType{
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
	return ShiftTypeConfig
}
