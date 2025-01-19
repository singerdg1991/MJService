package config

import (
	"github.com/hoitek/Maja-Service/internal/geartype/config"
)

// GearTypeConfig is a global variable for the geartype domain config
var GearTypeConfig config.ConfigType

// LoadGearTypeConfig loads the geartype domain config
func LoadGearTypeConfig() config.ConfigType {
	GearTypeConfig = config.ConfigType{
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
	return GearTypeConfig
}
