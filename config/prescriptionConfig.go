package config

import (
	"github.com/hoitek/Maja-Service/internal/prescription/config"
)

// PrescriptionConfig is a global variable for the prescription domain config
var PrescriptionConfig config.ConfigType

// LoadPrescriptionConfig loads the prescription domain config
func LoadPrescriptionConfig() config.ConfigType {
	PrescriptionConfig = config.ConfigType{
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
	return PrescriptionConfig
}
