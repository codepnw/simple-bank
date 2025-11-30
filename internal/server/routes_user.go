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

	public := cfg.router.Group("/auth")
	{
		public.POST("/register", handler.Register)
		public.POST("/login", handler.Login)
		public.POST("/refresh-token", handler.RefreshToken)
		public.POST("/logout", handler.Logout)
	}
}
