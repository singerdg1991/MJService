package config

import (
	"github.com/hoitek/Maja-Service/internal/report/config"
)

// ReportConfig is a global variable for the customer domain config
var ReportConfig config.ConfigType

// LoadReportConfig loads the customer domain config
func LoadReportConfig() config.ConfigType {
	ReportConfig = config.ConfigType{
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
		ForceMigrateAndSeed: false,
		MaxUploadSizeLimit:  AppConfig.MaxUploadSizeLimit,
	}
	return ReportConfig
}
