package config

type APIServerConfig struct {
	Address string `mapstructure:"address" validate:"required,hostname_port"`
}
