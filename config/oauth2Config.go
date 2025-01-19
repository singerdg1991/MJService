package config

import (
	"github.com/hoitek/Maja-Service/internal/oauth2/config"
)

// OAuth2Config is a global variable for the oauth2 domain config
var OAuth2Config config.ConfigType

// LoadOAuth2Config loads the oauth2 domain config
func LoadOAuth2Config() config.ConfigType {
	OAuth2Config = config.ConfigType{
		Debug:                     AppConfig.Debug,
		Environment:               AppConfig.Environment,
		ApiPrefix:                 AppConfig.ApiPrefix,
		ApiVersion1:               AppConfig.ApiVersion1,
		ApiVersion2:               AppConfig.ApiVersion2,
		DatabasePort:              AppConfig.DatabasePort,
		DatabaseName:              AppConfig.DatabaseName,
		DatabaseHost:              AppConfig.DatabaseHost,
		DatabaseUser:              AppConfig.DatabaseUser,
		DatabasePassword:          AppConfig.DatabasePassword,
		DatabaseSslMode:           AppConfig.DatabaseSslMode,
		DatabaseTimeZone:          AppConfig.DatabaseTimeZone,
		JwtTokenExpiration:        AppConfig.JwtTokenExpiration,
		JwtRefreshTokenExpiration: AppConfig.JwtRefreshTokenExpiration,
		JwtSigningKey:             AppConfig.JwtSigningKey,
		OTPCodeLength:             AppConfig.OTPCodeLength,
		OTPCodeExpirationSeconds:  AppConfig.OTPCodeExpirationSeconds,
		OTPGRPCHost:               AppConfig.OTPGRPCHost,
		OTPGRPCPort:               AppConfig.OTPGRPCPort,
		OTPGRPCTimeoutSeconds:     AppConfig.OTPGRPCTimeoutSeconds,
		OTPTestMode:               AppConfig.OTPTestMode,
		OTPEnable:                 AppConfig.OTPEnable,
	}
	return OAuth2Config
}
