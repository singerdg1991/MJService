package config

import (
	"github.com/hoitek/Maja-Service/internal/languageskill/config"
)

// LanguageSkillConfig is a global variable for the languageskill domain config
var LanguageSkillConfig config.ConfigType

// LoadLanguageSkillConfig loads the languageskill domain config
func LoadLanguageSkillConfig() config.ConfigType {
	LanguageSkillConfig = config.ConfigType{
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
	return LanguageSkillConfig
}
