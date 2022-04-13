package process

import (
	"github.com/lejenome/lro/pkg/config"
)

type ProcessExecutorConfig struct {
	App      config.AppConfig `mapstructure:",squash"`
	Redis    config.RedisConfig
	Database config.DatabaseConfig
	Nats     config.NatsConfig
}
