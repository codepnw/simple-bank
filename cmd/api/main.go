package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/codepnw/simple-bank/docs/swagger"
	"github.com/codepnw/simple-bank/internal/server"
	"github.com/codepnw/simple-bank/pkg/config"
	"github.com/codepnw/simple-bank/pkg/database"
	"github.com/codepnw/simple-bank/pkg/token"
	"github.com/codepnw/simple-bank/pkg/token/jwtmaker"
	"github.com/codepnw/simple-bank/pkg/token/pasetomaker"
	"golang.org/x/sync/errgroup"
)

// @title  Simple Bank API
// @version 1.0
// @description A robust banking API managing accounts and transfers.

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.LoadEnv(".env")
	if err != nil {
		log.Fatal(err)
	}

	swagger.SwaggerInfo.Host = cfg.Server.HTTPAddr
	swagger.SwaggerInfo.BasePath = cfg.Server.HTTPPrefix

	app, cleanup, err := initialize(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	var g errgroup.Group

	// gRPC Server
	g.Go(func() error {
		return server.RunGrpcServer(cfg, app.db, app.tx, app.token)
	})

	// HTTP Server
	g.Go(func() error {
		return server.RunHTTPServer(cfg, app.db, app.tx, app.token)
	})

	if err := g.Wait(); err != nil {
		log.Fatalf("server crashed: %v", err)
	}
}

type appContainer struct {
	db    *sql.DB
	tx    database.TxManager
	token token.TokenMaker
}

func initialize(cfg *config.EnvConfig) (*appContainer, func(), error) {
	// Connect Database
	db, err := database.ConnectPostgres(&cfg.DB)
	if err != nil {
		return nil, nil, fmt.Errorf("failed connect db: %v", err)
	}

	// New Transaction
	tx, err := database.NewTransaction(db)
	if err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed init tx: %v", err)
	}

	// New Token : JWT
	jwtToken, err := jwtmaker.NewJWTMaker(&cfg.JWT)
	if err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed init jwt: %v", err)
	}
	_ = jwtToken

	// New Token : PASETO
	pasetoToken, err := pasetomaker.NewPasetoMaker(cfg.Paseto.SymmetricKey)
	if err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed init paseto: %v", err)
	}

	cleanup := func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}

	return &appContainer{
		db:    db,
		tx:    tx,
		token: pasetoToken, // Change Token PASETO or JWT
	}, cleanup, nil
}
