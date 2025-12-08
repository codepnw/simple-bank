package server

import (
	accountrepository "github.com/codepnw/simple-bank/internal/features/account/repository"
	entryrepository "github.com/codepnw/simple-bank/internal/features/entry/repository"
	transferhandler "github.com/codepnw/simple-bank/internal/features/transfer/handler"
	transferrepository "github.com/codepnw/simple-bank/internal/features/transfer/repository"
	transferusecase "github.com/codepnw/simple-bank/internal/features/transfer/usecase"
)

func (cfg *routesConfig) registerTransferRoutes() {
	tranRepo := transferrepository.NewTransferRepository(cfg.db)
	accRepo := accountrepository.NewAccountRepository(cfg.db)
	entRepo := entryrepository.NewEntryRepository(cfg.db)

	uc := transferusecase.NewTransferUsecase(tranRepo, accRepo, entRepo, cfg.tx)
	handler := transferhandler.NewTransferHandler(uc)

	r := cfg.router.Group(cfg.prefix+"/transfers", cfg.mid.Authorized())
	{
		r.POST("", handler.CreateTransfer)
	}
}
