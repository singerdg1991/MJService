package config

import (
	"github.com/hoitek/Maja-Service/internal/trash/config"
)

// TrashConfig is a global variable for the trash domain config
var TrashConfig config.ConfigType

// LoadTrashConfig loads the trash domain config
func LoadTrashConfig() config.ConfigType {
	TrashConfig = config.ConfigType{
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
	return TrashConfig
}
