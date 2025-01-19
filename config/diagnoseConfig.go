package config

import (
	"github.com/hoitek/Maja-Service/internal/diagnose/config"
)

// DiagnoseConfig is a global variable for the diagnose domain config
var DiagnoseConfig config.ConfigType

// LoadDiagnoseConfig loads the diagnose domain config
func LoadDiagnoseConfig() config.ConfigType {
	DiagnoseConfig = config.ConfigType{
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
	return DiagnoseConfig
}
