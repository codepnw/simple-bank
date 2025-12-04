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
	"github.com/codepnw/simple-bank/pkg/jwt"
	"github.com/codepnw/simple-bank/pkg/utils/response"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type routesConfig struct {
	db     *sql.DB
	router *gin.Engine
	token  *jwt.JWTToken
	tx     database.TxManager
	mid    *middleware.AuthMiddleware
}

func RunHTTPServer(cfg *config.EnvConfig, db *sql.DB, tx database.TxManager, token *jwt.JWTToken) error {
	// New Middleware
	mid := middleware.NewMiddleware(token)

	// New Router
	gin.SetMode(gin.ReleaseMode)
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

	log.Printf("start HTTP server at %v", cfg.Server.HTTPAddr)
	return r.Run()
}

// ================ gRPC Server ====================

func RunGrpcServer(cfg *config.EnvConfig, db *sql.DB, tx database.TxManager, token *jwt.JWTToken) error {
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

func unaryServerInterceptor(token *jwt.JWTToken) grpc.UnaryServerInterceptor {
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
