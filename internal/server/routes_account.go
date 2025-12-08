package server

import (
	"github.com/codepnw/simple-bank/internal/consts"
	accounthandler "github.com/codepnw/simple-bank/internal/features/account/handler"
	accountrepository "github.com/codepnw/simple-bank/internal/features/account/repository"
	accountusecase "github.com/codepnw/simple-bank/internal/features/account/usecase"
)

func (cfg *routesConfig) registerAccountRoutes() {
	repo := accountrepository.NewAccountRepository(cfg.db)
	uc := accountusecase.NewAccountUsecase(repo)
	handler := accounthandler.NewAccountHandler(uc)

	r := cfg.router.Group(cfg.prefix+"/accounts", cfg.mid.Authorized())
	{
		r.POST("", handler.CreateAccount)
		r.GET("/:"+consts.ParamAccountID, handler.GetAccount)
		r.GET("", handler.ListAccounts)
	}
}
