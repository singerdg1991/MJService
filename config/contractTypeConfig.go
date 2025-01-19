package config

import (
	"github.com/hoitek/Maja-Service/internal/contracttype/config"
)

// ContractTypeConfig is a global variable for the contractType domain config
var ContractTypeConfig config.ConfigType

// LoadContractTypeConfig loads the contractType domain config
func LoadContractTypeConfig() config.ConfigType {
	ContractTypeConfig = config.ConfigType{
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
	return ContractTypeConfig
}
