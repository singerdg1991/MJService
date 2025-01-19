package config

import (
	"github.com/hoitek/Maja-Service/internal/email/config"
)

// EmailConfig is a global variable for the email domain config
var EmailConfig config.ConfigType

// LoadEmailConfig loads the email domain config
func LoadEmailConfig() config.ConfigType {
	EmailConfig = config.ConfigType{
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
	return EmailConfig
}
