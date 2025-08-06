package main

import (
	"log"

	"github.com/codepnw/simple-bank/config"
	"github.com/codepnw/simple-bank/internal/server"
)

const envFile = "dev.config.env"

func main() {
	cfg, err := config.LoadEnvConfig(envFile)
	if err != nil {
		log.Fatalf("load env failed: %v", err)
	}

	if err = server.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
