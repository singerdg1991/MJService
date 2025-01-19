package config

type ConfigType struct {
	Debug            bool   `mapstructure:"DEBUG"`
	Environment      string `mapstructure:"ENVIRONMENT"`
	ApiPrefix        string `mapstructure:"API_PREFIX"`
	ApiVersion1      string `mapstructure:"API_VERSION_1"`
	ApiVersion2      string `mapstructure:"API_VERSION_2"`
	DatabasePort     int    `mapstructure:"DATABASE_PORT"`
	DatabaseName     string `mapstructure:"DATABASE_NAME"`
	DatabaseHost     string `mapstructure:"DATABASE_HOST"`
	DatabaseUser     string `mapstructure:"DATABASE_USER"`
	DatabasePassword string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseSslMode  string `mapstructure:"DATABASE_SSL_MODE"`
	DatabaseTimeZone string `mapstructure:"DATABASE_TIMEZONE"`
}

var PaymentTypeConfig *ConfigType
