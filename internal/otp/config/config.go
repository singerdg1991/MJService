package config

type ConfigType struct {
	Debug                     bool   `mapstructure:"DEBUG"`
	Environment               string `mapstructure:"ENVIRONMENT"`
	ApiPrefix                 string `mapstructure:"API_PREFIX"`
	ApiVersion1               string `mapstructure:"API_VERSION_1"`
	ApiVersion2               string `mapstructure:"API_VERSION_2"`
	DatabasePort              int    `mapstructure:"DATABASE_PORT"`
	DatabaseName              string `mapstructure:"DATABASE_NAME"`
	DatabaseHost              string `mapstructure:"DATABASE_HOST"`
	DatabaseUser              string `mapstructure:"DATABASE_USER"`
	DatabasePassword          string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseSslMode           string `mapstructure:"DATABASE_SSL_MODE"`
	DatabaseTimeZone          string `mapstructure:"DATABASE_TIMEZONE"`
	JwtTokenExpiration        int64  `mapstructure:"JWT_TOKEN_EXPIRATION"`
	JwtRefreshTokenExpiration int64  `mapstructure:"JWT_REFRESH_TOKEN_EXPIRATION"`
	JwtSigningKey             string `mapstructure:"JWT_SIGNING_KEY"`
	OTPCodeLength             int    `mapstructure:"OTP_CODE_LENGTH"`
	OTPCodeExpirationSeconds  int    `mapstructure:"OTP_CODE_EXPIRATION_SECONDS"`
	OTPGRPCHost               string `mapstructure:"OTP_GRPC_HOST"`
	OTPGRPCPort               int    `mapstructure:"OTP_GRPC_PORT"`
	OTPGRPCTimeoutSeconds     int    `mapstructure:"OTP_GRPC_TIMEOUT_SECONDS"`
	OTPTestMode               bool   `mapstructure:"OTP_TEST_MODE"`
	OTPEnable                 bool   `mapstructure:"OTP_ENABLE"`
}

var OTPConfig *ConfigType
