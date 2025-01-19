package config

import (
	"github.com/hoitek/Maja-Service/internal/ticket/config"
)

// TicketConfig is a global variable for the ticket domain config
var TicketConfig config.ConfigType

// LoadTicketConfig loads the ticket domain config
func LoadTicketConfig() config.ConfigType {
	TicketConfig = config.ConfigType{
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
	return TicketConfig
}
