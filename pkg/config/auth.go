package config

import "time"

type AuthConfig struct {
	JWT JWTConfig
}
type JWTConfig struct {
	SecretKey string        `mapstructure:"secret_key" validate:"required,min=32"`
	Duration  time.Duration `validate:"min=1h"`
}
