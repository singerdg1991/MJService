package config

import (
	"github.com/hoitek/Maja-Service/internal/ability/config"
)

// AbilityConfig is a global variable for the ability domain config
var AbilityConfig config.ConfigType

// LoadAbilityConfig loads the ability domain config
func LoadAbilityConfig() config.ConfigType {
	AbilityConfig = config.ConfigType{
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
	return AbilityConfig
}
