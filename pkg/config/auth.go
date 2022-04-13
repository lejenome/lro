package config

type AuthConfig struct {
	JWT JWTConfig
}
type JWTConfig struct {
	SecretKey string `mapstructure:"secret_key" validate:"required,min=32"`
}
