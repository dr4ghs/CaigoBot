package configs

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// BotConfig structs that represents the service configs
type BotConfig struct {
	BotToken      string `yaml:"bot_token" envconfig:"DISCORD_BOT_TOKEN"`
	RemoveCommand bool   `yaml:"rm_cmd"`
	Intents       int    `yaml:"intents"`
	Guilds        []struct {
		ID       string   `yaml:"id"`
		Commands []string `yaml:"commands"`
	} `yaml:"guilds"`
}

const (
	// Config file path environment variable name
	CONFIG_FILE_PATH_VAR_ENV = "CONFIG_FILE_PATH"

	// Default config file path
	DEFAULT_CONFIG_FILE_PATH = "config.yml"
)

// Config instance
var Config *BotConfig

func init() {
	Config = &BotConfig{}

	readFile(Config)
	readEnv(Config)
}

// readFile reads the config YAML file and populate the Config variable
func readFile(c *BotConfig) {
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

// readEnv parses the environment variables
func readEnv(c *BotConfig) {
	err := envconfig.Process("", c)
	if err != nil {
		log.Fatalf("Cannot process environment variables: %v", err)
	}
}
