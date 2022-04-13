package config

import (
	"errors"
	"log"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Env  string `validate:"oneof='production' 'dev'"`
	Name string `validate:"required"`
}
type DefaultConfig struct {
	App      AppConfig `mapstructure:",squash"`
	Database DatabaseConfig
	Redis    RedisConfig
	Nats     NatsConfig
}

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

func validateConfig(conf interface{}) error {
	validate := validator.New()
	err := validate.Struct(conf)
	return err
}

func Load(conf interface{}) error {

	val := reflect.ValueOf(conf)

	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return errors.New("Object should be a struct")
	}

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
	err = viper.Unmarshal(conf)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	log.Printf("Loaded config: %+v\n", conf)

	if err = validateConfig(conf); err != nil {
		log.Fatalf("Unvalid config struct content: %s\n", err.Error())
	}
	return nil
}
