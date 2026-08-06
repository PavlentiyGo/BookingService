package core_config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// LOGGER
	Level  string `envconfig:"LOGGER_LEVEL"    default:"DEBUG"`
	Folder string `envconfig:"LOGGER_FOLDER"` // exports via make command
	// SERVER
	Addr            string        `envconfig:"SERVER_ADDR"     default:"8080"`
	ShutdownTimeout time.Duration `envconfig:"SERVER_SHUTDOWN_TIME" default:"5s"`
	// JWT
	SigningKey string        `envconfig:"JWT_SIGNING_KEY" default:"secret_key"`
	LifeTime   time.Duration `envconfig:"JWT_LIFE_TIME"  default:"24h"`
	// HASH
	BcryptCost int `envconfig:"JWT_BCRYPT_COST" default:"12"`
	DbConfig   DatabaseConfig
}
type DatabaseConfig struct {
	User     string        `envconfig:"POSTGRES_USER" default:"postgres"`
	Password string        `envconfig:"POSTGRES_PASSWORD" default:"password"`
	Name     string        `envconfig:"POSTGRES_DB" default:"postgres"`
	Host     string        `envconfig:"POSTGRES_HOST" default:"localhost"`
	Port     string        `envconfig:"POSTGRES_PORT" default:"5432"`
	Timeout  time.Duration `evnconfig:"DATABASE_TIMEOUtT" default:"5s"`
}

func NewConfig() (Config, error) {

	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		panic(err)
	}
	return config
}
