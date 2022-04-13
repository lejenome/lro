package config

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type JWTConfig struct {
	SecretKey string `validate:"required,min=32"`
}
type DatabaseConfig struct {
	Driver string `validate:"oneof='postgresql'"`
	URL    string `validate:"required"`
}
type RedisConfig struct {
	Network  string `json:"network" validate:"oneof='tcp' 'unix'"`
	URL      string `json:"url" validate:"required,unix_addr|hostname_port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB       int    `json:"db" validate:"min=0"`
}
type NatsConfig struct {
	URL string `json:"url" validate:"required,hostname_port"`
}
type Config struct {
	Env      string `validate:"oneof='production' 'dev'"`
	Name     string `validate:"required"`
	Database DatabaseConfig
	/* TODO
	Auth     struct {
		JWT JWTConfig
	}
	*/
	Redis RedisConfig
	Nats  NatsConfig
}

var Conf Config

func Defaults() map[string]interface{} {
	return map[string]interface{}{
		"Env":             "dev",
		"Database.Driver": "postgresql",
		"Redis.Network":   "tcp",
		"Redis.URL":       "localhost:6379",
		"Redis.DB":        0,
		"Nats.URL":        "localhost:4222",
	}
}

func validateConfig(config *Config) error {
	validate := validator.New()
	err := validate.Struct(config)
	return err
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("app")
	viper.AutomaticEnv()
	for k, v := range Defaults() {
		viper.SetDefault(k, v)
	}
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Panicf("Fatal error config file: %s \n", err)
	}
	err = viper.Unmarshal(&Conf)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	log.Printf("Loaded config: %+v\n", Conf)

	if err = validateConfig(&Conf); err != nil {
		log.Fatalf("Unvalid config struct content: %s\n", err.Error())
	}
	return &Conf, nil
}

func init() {
	Load()
}
