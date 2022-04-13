package process

import (
	"github.com/lejenome/lro/pkg/config"
)

type ProcessApiConfig struct {
	App      config.AppConfig `mapstructure:",squash"`
	Server   config.APIServerConfig
	Redis    config.RedisConfig
	Database config.DatabaseConfig
	Nats     config.NatsConfig
	Auth     config.AuthConfig
}
