package config

import (
	"github.com/hoitek/Maja-Service/internal/punishment/config"
)

// PunishmentConfig is a global variable for the punishment domain config
var PunishmentConfig config.ConfigType

// LoadPunishmentConfig loads the punishment domain config
func LoadPunishmentConfig() config.ConfigType {
	PunishmentConfig = config.ConfigType{
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
	return PunishmentConfig
}
