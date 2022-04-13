package config

type RedisConfig struct {
	Network  string `mapstructure:"network" validate:"oneof='tcp' 'unix'"`
	URL      string `mapstructure:"url" validate:"required,unix_addr|hostname_port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db" validate:"min=0"`
}
