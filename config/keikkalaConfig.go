package config

import (
	"github.com/hoitek/Maja-Service/internal/keikkala/config"
)

// KeikkalaConfig is a global variable for the keikkala domain config
var KeikkalaConfig config.ConfigType

// LoadKeikkalaConfig loads the keikkala domain config
func LoadKeikkalaConfig() config.ConfigType {
	KeikkalaConfig = config.ConfigType{
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
	return KeikkalaConfig
}
