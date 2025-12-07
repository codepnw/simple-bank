package main

import (
	"log"

	"github.com/codepnw/simple-bank/internal/server"
	"github.com/codepnw/simple-bank/pkg/config"
	"github.com/codepnw/simple-bank/pkg/database"
	"github.com/codepnw/simple-bank/pkg/jwt"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg, err := config.LoadEnv(".env")
	if err != nil {
		log.Fatal(err)
	}

	// Connect Database
	db, err := database.ConnectPostgres(&cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// New Transaction
	tx, err := database.NewTransaction(db)
	if err != nil {
		log.Fatal(err)
	}

	// Init JWT
	token, err := jwt.InitJWT(&cfg.JWT)
	if err != nil {
		log.Fatal(err)
	}

	var g errgroup.Group

	// gRPC Server
	g.Go(func() error {
		return server.RunGrpcServer(cfg, db, tx, token)
	})

	// HTTP Server
	g.Go(func() error {
		return server.RunHTTPServer(cfg, db, tx, token)
	})

	if err := g.Wait(); err != nil {
		log.Fatalf("server crashed: %v", err)
	}
}
