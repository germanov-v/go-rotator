package main

import (
	"flag"
	"github.com/germanov-v/go-rotator/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"log"
)

func main() {
	conf, _ := config.ReadConsoleConfigParameter()
	log.Println("CONNECT: " + conf.DbBaseConfig.MigrateConnection)
	//showFolders(dir, "")
	migrationsPath := flag.String("migration", "file://migrations", "path to config file")
	m, err := migrate.New(*migrationsPath, conf.DbBaseConfig.MigrateConnection)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
}
