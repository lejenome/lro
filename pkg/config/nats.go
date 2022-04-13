package config

type NatsConfig struct {
	URL string `mapstructure:"url" validate:"required,hostname_port"`
}
