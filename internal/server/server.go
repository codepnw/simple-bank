package server

import (
	"net/http"

	"github.com/codepnw/simple-bank/pkg/config"
	"github.com/codepnw/simple-bank/pkg/database"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.EnvConfig) error {
	db, err := database.ConnectPostgres(&cfg.DB)
	if err != nil {
		return err
	}
	defer db.Close()

	r := gin.New()
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "server running...",
		})
	})

	return r.Run()
}
