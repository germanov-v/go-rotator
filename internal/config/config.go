package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Config struct {
	DbBaseConfig `json:"Database"`
	ServerConfig `json:"Server"`
}

type ServerConfig struct {
	Host           string `json:"Host"`
	Port           int    `json:"Port"`
	ApiGatewayPort int    `json:"ApiGatewayPort"`
	ApiGatewayHost string `json:"ApiGatewayHost"`
}

type DbBaseConfig struct {
	ConnectionString  string `json:"ConnectionString"`
	MigrateConnection string `json:"MigrateConnection"`
}

func ReadConsoleConfigParameter() (*Config, error) {
	path := flag.String("config", "config.json", "path to config file")
	flag.Parse()
	config, err := LoadConfig(*path)
	if err != nil {
		panic(err)
	}
	////overrideFromEnvDynamic(config)
	return config, err
}

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	//overrideString(&cfg.DbBaseConfig.ConnectionString, "DATABASE_CONNECTIONSTRING")
	//overrideString(&cfg.DbBaseConfig.MigrateConnection, "DATABASE_MIGRATECONNECTION")

	return cfg, nil
}
