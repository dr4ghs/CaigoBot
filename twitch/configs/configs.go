package configs

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

type TwitchConfig struct {
	ClientID       string `yaml:"client_id" envconfig:"TWITCH_CLIENT_ID"`
	AppAccessToken string `yaml:"app_access_token" envconfig:"TWITCH_APP_SECRET_TOKEN"`
}

const (
	CONFIG_FILE_PATH_VAR_ENV = "CONFIG_FILE_PATH"

	DEFAULT_CONFIG_FILE_PATH = "config.yml"
)

var Config *TwitchConfig

func init() {
	Config = &TwitchConfig{}

	readFile(Config)
	readEnv(Config)
}

func readFile(c *TwitchConfig) {
	configFilePath := os.Getenv(CONFIG_FILE_PATH_VAR_ENV)
	if len(configFilePath) == 0 {
		configFilePath = DEFAULT_CONFIG_FILE_PATH
	}

	file, err := os.Open(configFilePath)
	if err != nil {
		log.Fatalf("Cannot open config file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		log.Fatalf("Error while decoding config file: %v", err)
	}
}

func readEnv(c *TwitchConfig) {
	err := envconfig.Process("", c)
	if err != nil {
		log.Fatalf("Cannot process environment variables: %v", err)
	}
}
