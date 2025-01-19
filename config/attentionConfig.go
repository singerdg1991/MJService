package config

import (
	"github.com/hoitek/Maja-Service/internal/staffclub/attention/config"
)

// AttentionConfig is a global variable for the attention domain config
var AttentionConfig config.ConfigType

// LoadAttentionConfig loads the attention domain config
func LoadAttentionConfig() config.ConfigType {
	AttentionConfig = config.ConfigType{
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
	return AttentionConfig
}
