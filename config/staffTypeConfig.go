package config

import (
	"github.com/hoitek/Maja-Service/internal/stafftype/config"
)

// StaffTypeConfig is a global variable for the staffType domain config
var StaffTypeConfig config.ConfigType

// LoadStaffTypeConfig loads the staffType domain config
func LoadStaffTypeConfig() config.ConfigType {
	StaffTypeConfig = config.ConfigType{
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
	return StaffTypeConfig
}
