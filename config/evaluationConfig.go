package config

import (
	"github.com/hoitek/Maja-Service/internal/evaluation/config"
)

// EvaluationConfig is a global variable for the evaluation domain config
var EvaluationConfig config.ConfigType

// LoadEvaluationConfig loads the evaluation domain config
func LoadEvaluationConfig() config.ConfigType {
	EvaluationConfig = config.ConfigType{
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
	return EvaluationConfig
}
