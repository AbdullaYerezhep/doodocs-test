package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	HTTPServer  `yaml:"http_server"`
	SmtpServer  `yaml:"smtp_server"`
	Gmail 		`yaml:"gmail"`
}

type HTTPServer struct {
	Address     string        `yaml: "address" env_default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env_default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env_default:"60s"`
}

type SmtpServer struct {
	Host string	`yaml:"host" env_default:"smtp.gmail.com"`
	Port string	`yaml:"port" env_default:"587"`
}

type Gmail struct {
	Login string `yaml:"login" env_default:""`
	Password string `yaml:"password" env_default:""`
}

func MustLoad() *Config {
	if isRunningTests() {
		err := godotenv.Load("../../local.env")
		if err != nil {
			log.Fatalf("Config for test not loaded: %s does not exist", err)
		}
		configPath := os.Getenv("CONFIG_PATH_TEST")

		if configPath == "" {
			configPath = "../../config/config.yaml"
		}
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Fatalf("Config file %s does not exist", configPath)
		}
		var cfg Config
		configData, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatal("Could not read config file")
		}
		err = yaml.Unmarshal(configData, &cfg)
		if err != nil {
			log.Fatal(err)
		}
		// Set the storage path relative to the test environment
		return &cfg
	}

	err := godotenv.Load("./local.env")
	if err != nil {
		log.Fatalf("Config not loaded: %s does not exist", err)
	}

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		configPath = "./config/config.yaml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file %s does not exist", configPath)
	}

	var cfg Config
	configData, err := os.ReadFile(configPath)

	if err != nil {
		log.Fatal("Could not read config file")
	}

	err = yaml.Unmarshal(configData, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	// Set the storage path relative to the application's execution directory
	return &cfg
}


func isRunningTests() bool {
	wd, _ := os.Getwd()
	return strings.HasSuffix(wd, "/internal/handler")
}
