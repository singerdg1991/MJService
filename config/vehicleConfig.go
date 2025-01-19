package config

import (
	"github.com/hoitek/Maja-Service/internal/vehicle/config"
)

// VehicleConfig is a global variable for the vehicle domain config
var VehicleConfig config.ConfigType

// LoadVehicleConfig loads the vehicle domain config
func LoadVehicleConfig() config.ConfigType {
	VehicleConfig = config.ConfigType{
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
	return VehicleConfig
}
