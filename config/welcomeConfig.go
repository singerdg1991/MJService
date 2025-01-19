package config

import "github.com/hoitek/Maja-Service/internal/welcome/config"

// WelcomeConfig is a global variable for the welcome domain config
var WelcomeConfig config.ConfigType

// LoadWelcomeConfig loads the welcome domain config
func LoadWelcomeConfig() config.ConfigType {
	WelcomeConfig = config.ConfigType{}
	return WelcomeConfig
}
