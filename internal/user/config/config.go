package config

type ConfigType struct {
	Debug               bool   `mapstructure:"DEBUG"`
	Environment         string `mapstructure:"ENVIRONMENT"`
	ApiPrefix           string `mapstructure:"API_PREFIX"`
	ApiVersion1         string `mapstructure:"API_VERSION_1"`
	ApiVersion2         string `mapstructure:"API_VERSION_2"`
	DatabasePort        int    `mapstructure:"DATABASE_PORT"`
	DatabaseName        string `mapstructure:"DATABASE_NAME"`
	DatabaseHost        string `mapstructure:"DATABASE_HOST"`
	DatabaseUser        string `mapstructure:"DATABASE_USER"`
	DatabasePassword    string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseSslMode     string `mapstructure:"DATABASE_SSL_MODE"`
	DatabaseTimeZone    string `mapstructure:"DATABASE_TIMEZONE"`
	DatabaseMongoDBHost string `mapstructure:"DATABASE_MONGODB_HOST"`
	DatabaseMongoDBPort int    `mapstructure:"DATABASE_MONGODB_PORT"`
	DatabaseMongoDBName string `mapstructure:"DATABASE_MONGODB_NAME"`
	DatabaseMongoDBUser string `mapstructure:"DATABASE_MONGODB_USER"`
	DatabaseMongoDBPass string `mapstructure:"DATABASE_MONGODB_PASS"`
}

var UserConfig *ConfigType
