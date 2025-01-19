package config

import "github.com/hoitek/Maja-Service/internal/static/config"

// StaticConfig is a global variable for the static domain config
var StaticConfig config.ConfigType

// LoadStaticConfig loads the static domain config
func LoadStaticConfig() config.ConfigType {
	StaticConfig = config.ConfigType{
		StaticURIPath: "/",
		StaticDir:     "public/",
	}
	return StaticConfig
}
