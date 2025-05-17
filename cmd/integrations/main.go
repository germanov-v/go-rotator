package main

import (
	"context"
	"github.com/germanov-v/go-rotator/internal/config"
	"github.com/germanov-v/go-rotator/internal/integrations"
	"log"
	"os"
	"strings"
)

func main() {
	ctx := context.Background()
	cfg, err := config.ReadConsoleConfigParameter()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
		os.Exit(1)
	}

	log.Printf("VARIABLES count: %d", len(os.Environ()))

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		log.Printf("%s = %s", parts[0], parts[1])
	}

	if err := integrations.RunAll(ctx, cfg); err != nil {
		log.Fatalf("Integration tests failed: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
