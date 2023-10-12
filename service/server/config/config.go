package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string           `yaml:"env" env-default:"local"`
	Database   DatabaseConfig   `yaml:"database"`
	HttpServer HttpServerConfig `yaml:"http_server"`
}

type DatabaseConfig struct {
	Host    string `yaml:"host" env-default:"localhost"`
	Port    int    `yaml:"port" env-default:"5432"`
	User    string `yaml:"user" env-required:"true"`
	DBName  string `yaml:"dbname" env-default:"postgres"`
	SSLMode string `yaml:"sslmode" env-default:"disable"`
}

type HttpServerConfig struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func Load(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, err
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, err
	}

	log.Println("Config read successful!")

	return &cfg, nil
}
