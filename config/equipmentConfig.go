package config

import (
	"github.com/hoitek/Maja-Service/internal/equipment/config"
)

// EquipmentConfig is a global variable for the equipment domain config
var EquipmentConfig config.ConfigType

// LoadEquipmentConfig loads the equipment domain config
func LoadEquipmentConfig() config.ConfigType {
	EquipmentConfig = config.ConfigType{
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
	return EquipmentConfig
}
