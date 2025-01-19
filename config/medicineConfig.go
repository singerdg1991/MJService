package config

import (
	"github.com/hoitek/Maja-Service/internal/medicine/config"
)

// MedicineConfig is a global variable for the medicine domain config
var MedicineConfig config.ConfigType

// LoadMedicineConfig loads the medicine domain config
func LoadMedicineConfig() config.ConfigType {
	MedicineConfig = config.ConfigType{
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
	return MedicineConfig
}
