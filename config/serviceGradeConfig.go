package config

import (
	"github.com/hoitek/Maja-Service/internal/servicegrade/config"
)

// ServiceGradeConfig is a global variable for the serviceGrade domain config
var ServiceGradeConfig config.ConfigType

// LoadServiceGradeConfig loads the serviceGrade domain config
func LoadServiceGradeConfig() config.ConfigType {
	ServiceGradeConfig = config.ConfigType{
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
	return ServiceGradeConfig
}
