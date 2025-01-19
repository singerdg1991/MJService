package config

import "github.com/hoitek/Maja-Service/internal/healthcheck/config"

// HealthCheckConfig is a global variable for the health check domain config
var HealthCheckConfig config.ConfigType

// LoadHealthCheckConfig loads the health check domain config
func LoadHealthCheckConfig() config.ConfigType {
	HealthCheckConfig = config.ConfigType{}
	return HealthCheckConfig
}
