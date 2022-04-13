package config

type DatabaseConfig struct {
	Driver string `validate:"oneof='postgresql'"`
	URL    string `validate:"required"`
}
