package config

import (
	"github.com/hoitek/Maja-Service/internal/push/config"
)

// PushConfig is a global variable for the push domain config
var PushConfig config.ConfigType

// LoadPushConfig loads the push domain config
func LoadPushConfig() config.ConfigType {
	PushConfig = config.ConfigType{
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
		VAPIDPublicKey:   AppConfig.VAPIDPublicKey,
		VAPIDPrivateKey:  AppConfig.VAPIDPrivateKey,
	}
	return PushConfig
}
