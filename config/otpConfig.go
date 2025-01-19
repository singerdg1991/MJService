package config

import (
	"github.com/hoitek/Maja-Service/internal/otp/config"
)

// OTPConfig is a global variable for the otp domain config
var OTPConfig config.ConfigType

// LoadOTPConfig loads the otp domain config
func LoadOTPConfig() config.ConfigType {
	OTPConfig = config.ConfigType{
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
	return OTPConfig
}
