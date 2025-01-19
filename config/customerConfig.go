package config

import (
	"github.com/hoitek/Maja-Service/internal/customer/config"
)

// CustomerConfig is a global variable for the customer domain config
var CustomerConfig config.ConfigType

// LoadCustomerConfig loads the customer domain config
func LoadCustomerConfig() config.ConfigType {
	CustomerConfig = config.ConfigType{
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
	return CustomerConfig
}
