package storage

import (
	"apibgo/pkg/univenv"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	PgSql Clusters `yaml:"pgsql"`
	MySql Clusters `yaml:"mysql"`
}

type Clusters struct {
	Master Cluster `yaml:"master"`
	Slave  Cluster `yaml:"slave"`
}

type Cluster struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath + "database.yaml"); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	withEnv, err := univenv.YamlWithEnv(configPath + "database.yaml")

	if err != nil {
		log.Fatalf("cannot set env variables to yaml a file: %s", err)
	}

	var cfg Config

	if err := cleanenv.ParseYAML(withEnv, &cfg); err != nil {
		log.Fatalf("DATABASE.YAML -> cannot read config: %s", err)
	}

	return &cfg
}
