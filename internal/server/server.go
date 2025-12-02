package server

import (
	"database/sql"

	"github.com/codepnw/simple-bank/internal/middleware"
	"github.com/codepnw/simple-bank/pkg/config"
	"github.com/codepnw/simple-bank/pkg/database"
	"github.com/codepnw/simple-bank/pkg/jwt"
	"github.com/codepnw/simple-bank/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type routesConfig struct {
	db     *sql.DB
	router *gin.Engine
	token  *jwt.JWTToken
	tx     database.TxManager
	mid    *middleware.AuthMiddleware
}

func Run(cfg *config.EnvConfig) error {
	// Connect Database
	db, err := database.ConnectPostgres(&cfg.DB)
	if err != nil {
		return err
	}
	defer db.Close()

	// Init JWT
	token, err := jwt.InitJWT(&cfg.JWT)
	if err != nil {
		return err
	}

	// New Transaction
	tx := database.NewTransaction(db)

	// New Middleware
	mid := middleware.NewMiddleware(token)

	// New Router
	r := gin.New()
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, "server running...", nil)
	})

	// Config Routes
	routes := &routesConfig{
		router: r,
		token:  token,
		db:     db,
		tx:     tx,
		mid:    mid,
	}
	routes.registerUserRoutes()
	routes.registerAccountRoutes()
	routes.registerTransferRoutes()

	return r.Run()
}
