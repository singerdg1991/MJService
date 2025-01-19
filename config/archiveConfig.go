package config

import (
	"github.com/hoitek/Maja-Service/internal/archive/config"
)

// ArchiveConfig is a global variable for the archive domain config
var ArchiveConfig config.ConfigType

// LoadArchiveConfig loads the archive domain config
func LoadArchiveConfig() config.ConfigType {
	ArchiveConfig = config.ConfigType{
		Debug:              AppConfig.Debug,
		Environment:        AppConfig.Environment,
		ApiPrefix:          AppConfig.ApiPrefix,
		ApiVersion1:        AppConfig.ApiVersion1,
		ApiVersion2:        AppConfig.ApiVersion2,
		DatabasePort:       AppConfig.DatabasePort,
		DatabaseName:       AppConfig.DatabaseName,
		DatabaseHost:       AppConfig.DatabaseHost,
		DatabaseUser:       AppConfig.DatabaseUser,
		DatabasePassword:   AppConfig.DatabasePassword,
		DatabaseSslMode:    AppConfig.DatabaseSslMode,
		DatabaseTimeZone:   AppConfig.DatabaseTimeZone,
		MaxUploadSizeLimit: AppConfig.MaxUploadSizeLimit,
	}
	return ArchiveConfig
}
