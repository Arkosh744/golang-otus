package config

import (
	"context"
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var configFile string

const defaultConfigPath = "config.yaml"

const (
	StorageMemory   = "memory"
	StoragePostgres = "postgres"
)

//nolint:gochecknoinits // it is ok to use init to get config file path from flags
func init() {
	flag.StringVar(&configFile, "config", defaultConfigPath, "Path to configuration file")
}

type Config struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`

	Log struct {
		Preset string `yaml:"preset"`
	} `yaml:"log"`

	Storage  string `yaml:"storage"`
	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"postgres"`
}

var AppConfig = Config{}

func Init(_ context.Context) error {
	rawYaml, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	err = yaml.Unmarshal(rawYaml, &AppConfig)
	if err != nil {
		return fmt.Errorf("parse config file: %w", err)
	}

	return nil
}

func (c *Config) GetPostgresDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Postgres.Host, c.Postgres.Port, c.Postgres.User, c.Postgres.Password, c.Postgres.Database)
}

func (c *Config) GetStorage() string {
	switch c.Storage {
	case StorageMemory:
		return StorageMemory
	case StoragePostgres:
		return StoragePostgres
	default:
		return StorageMemory
	}
}
