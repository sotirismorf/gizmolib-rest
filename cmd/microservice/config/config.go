package config

import (
	"bytes"
	_ "embed"
	// "strings"

	"github.com/spf13/viper"
)

//go:embed config.yml
var defaultConfiguration []byte

type Admin struct {
	Username string
	Password string
}

type Postgres struct {
	Host     string
	User     string
	Password string
}

type Config struct {
	Postgres *Postgres
	Admin    *Admin
}

func Read() (*Config, error) {
	viper.AllowEmptyEnv(false)
	viper.SetConfigType("yml")

	// Defaults
	viper.SetDefault("admin.password", "password")
	viper.SetDefault("admin.username", "admin")
	viper.SetDefault("postgres.database", "postgres")
	viper.SetDefault("postgres.host", "localhost")
	viper.SetDefault("postgres.password", "password")
	viper.SetDefault("postgres.user", "postgres")

	// Environment variables
	viper.BindEnv("admin.password", "GL_ADMIN_PASSWORD")
	viper.BindEnv("admin.username", "GL_ADMIN_USER")
	viper.BindEnv("postgres.database", "PG_DATABASE")
	viper.BindEnv("postgres.host", "PG_HOST")
	viper.BindEnv("postgres.password", "PG_PASSWORD")
	viper.BindEnv("postgres.user", "PG_USER")

	// Read configuration
	if err := viper.ReadConfig(bytes.NewBuffer(defaultConfiguration)); err != nil {
		return nil, err
	}

	// Unmarshal the configuration
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
