package server

import (
	userhandler "github.com/codepnw/simple-bank/internal/features/user/handler"
	userrepository "github.com/codepnw/simple-bank/internal/features/user/repository"
	userusecase "github.com/codepnw/simple-bank/internal/features/user/usecase"
)

func (cfg *routesConfig) registerUserRoutes() {
	repo := userrepository.NewUserRepository(cfg.db)
	uc := userusecase.NewUserUsecase(repo, cfg.token, cfg.tx)
	handler := userhandler.NewUserHandler(uc)

	auth := cfg.router.Group("/auth")
	{
		// Public
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)

		// Private
		auth.POST("/refresh-token", cfg.mid.Authorized(), handler.RefreshToken)
		auth.POST("/logout", cfg.mid.Authorized(), handler.Logout)
	}
}
