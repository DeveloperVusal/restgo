package config

import (
	"apibgo/pkg/univenv"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath + "main.yaml"); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	withEnv, err := univenv.YamlWithEnv(configPath + "main.yaml")

	if err != nil {
		log.Fatalf("cannot set env variables to yaml a file: %s", err)
	}

	var cfg Config

	if err := cleanenv.ParseYAML(withEnv, &cfg); err != nil {
		log.Fatalf("MAIN.YAML -> cannot read config: %s", err)
	}

	return &cfg
}
