package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DbBaseConfig `json:"Database"`
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

	overrideString(&cfg.DbBaseConfig.ConnectionString, "DATABASE_CONNECTIONSTRING")
	overrideString(&cfg.DbBaseConfig.MigrateConnection, "DATABASE_MIGRATECONNECTION")

	return cfg, nil
}

func overrideString(field *string, envName string) {
	if v, ok := os.LookupEnv(envName); ok && v != "" {
		*field = v
	}
}

func overrideInt(field *int, envName string) {
	if v, ok := os.LookupEnv(envName); ok && v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			*field = i
		}
	}
}
