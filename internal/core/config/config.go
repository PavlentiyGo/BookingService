package core_config

import (
	"os/exec"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// LOGGER
	Level  string `envconfig:"LOGGER_LEVEL"    default:"DEBUG"`
	Folder string `envconfig:"LOGGER_FOLDER"`
	// SERVER
	Addr string `envconfig:"SERVER_ADDR"     default:"8080"`
	// JWT
	SigningKey string        `envconfig:"JWT_SIGNING_KEY" default:"secret_key"`
	LifeTime   time.Duration `envconfig:"JWT_LIFE_TIME"  default:"24h"`
	// HASH
	BcryptCost int `envconfig:"JWT_BCRYPT_COST" default:"12"`
}

func NewConfig() (Config, error) {

	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		return Config{}, err
	}

	cmd := exec.Command("bash", "-c", "pwd")
	output, _ := cmd.Output()
	currentFolderStr := strings.TrimSpace(string(output)) + "/log"
	config.Folder = currentFolderStr

	return config, nil
}
func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		panic(err)
	}
	return config
}
