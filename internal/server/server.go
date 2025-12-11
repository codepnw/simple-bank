package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"

	accountrepository "github.com/codepnw/simple-bank/internal/features/account/repository"
	entryrepository "github.com/codepnw/simple-bank/internal/features/entry/repository"
	transfergrpc "github.com/codepnw/simple-bank/internal/features/transfer/grpc"
	transferrepository "github.com/codepnw/simple-bank/internal/features/transfer/repository"
	transferusecase "github.com/codepnw/simple-bank/internal/features/transfer/usecase"
	"github.com/codepnw/simple-bank/internal/middleware"
	pb "github.com/codepnw/simple-bank/pb/proto"
	"github.com/codepnw/simple-bank/pkg/auth"
	"github.com/codepnw/simple-bank/pkg/config"
	"github.com/codepnw/simple-bank/pkg/database"
	"github.com/codepnw/simple-bank/pkg/token"
	"github.com/codepnw/simple-bank/pkg/utils/response"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	_ "github.com/codepnw/simple-bank/docs/swagger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type routesConfig struct {
	db     *sql.DB
	router *gin.Engine
	prefix string
	token  token.TokenMaker
	tx     database.TxManager
	mid    *middleware.AuthMiddleware
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	// New Router
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS", "DELETE", "PATCH"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true

	r.Use(cors.New(corsConfig))

	r.GET("/health", func(c *gin.Context) {
		response.Success(c, "server running...", nil)
	})
	// API Docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func RunHTTPServer(cfg *config.EnvConfig, db *sql.DB, tx database.TxManager, token token.TokenMaker) error {
	// New Middleware
	mid := middleware.NewMiddleware(token)

	// Setup Router
	router := setupRouter()

	// Config Routes
	routes := &routesConfig{
		router: router,
		prefix: cfg.Server.HTTPPrefix,
		token:  token,
		db:     db,
		tx:     tx,
		mid:    mid,
	}
	routes.registerUserRoutes()
	routes.registerAccountRoutes()
	routes.registerTransferRoutes()

	addr := cfg.Server.HTTPAddr
	log.Printf("start HTTP docs at %s/swagger/index.html", addr)
	log.Printf("start HTTP server at %s%s", addr, routes.prefix)
	return router.Run(addr)
}

// ================ gRPC Server ====================

func RunGrpcServer(cfg *config.EnvConfig, db *sql.DB, tx database.TxManager, token token.TokenMaker) error {
	tranRepo := transferrepository.NewTransferRepository(db)
	accRepo := accountrepository.NewAccountRepository(db)
	entRepo := entryrepository.NewEntryRepository(db)

	uc := transferusecase.NewTransferUsecase(tranRepo, accRepo, entRepo, tx)

	server := transfergrpc.NewTransferServer(uc)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(unaryServerInterceptor(token)),
	)

	pb.RegisterSimpleBankServer(grpcServer, server)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", cfg.Server.GrpcAddr)
	if err != nil {
		return fmt.Errorf("cannot create listener: %v", err)
	}

	log.Printf("start gRPC server at %v", cfg.Server.GrpcAddr)
	if err = grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("cannot start grpc server: %v", err)
	}
	return nil
}

func unaryServerInterceptor(token token.TokenMaker) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		authHeader := values[0]
		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			return nil, status.Errorf(codes.Unauthenticated, "invalid authorization header format")
		}
		authType := strings.ToLower(fields[0])
		if authType != "bearer" {
			return nil, status.Errorf(codes.Unauthenticated, "unsupported authorization type: %s", authType)
		}
		accessToken := fields[1]

		payload, err := token.VerifyAccessToken(accessToken)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
		}

		ctx = auth.SetUserID(ctx, payload.UserID)
		return handler(ctx, req)
	}
}
