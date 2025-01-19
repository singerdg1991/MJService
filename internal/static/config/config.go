package config

type ConfigType struct {
	StaticURIPath string `mapstructure:"STATIC_URI_PATH"`
	StaticDir     string `mapstructure:"STATIC_DIR"`
}

var StaticConfig *ConfigType
