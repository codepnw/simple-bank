package main

import (
	"log"

	"github.com/codepnw/simple-bank/internal/server"
	"github.com/codepnw/simple-bank/pkg/config"
)

func main() {
	cfg, err := config.LoadEnv("dev.env")
	if err != nil {
		log.Fatal(err)
	}

	if err = server.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
