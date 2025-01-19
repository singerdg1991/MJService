package config

import (
	"github.com/hoitek/Maja-Service/internal/vehicletype/config"
)

// VehicleTypeConfig is a global variable for the vehicletype domain config
var VehicleTypeConfig config.ConfigType

// LoadVehicleTypeConfig loads the vehicletype domain config
func LoadVehicleTypeConfig() config.ConfigType {
	VehicleTypeConfig = config.ConfigType{
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
	return VehicleTypeConfig
}
