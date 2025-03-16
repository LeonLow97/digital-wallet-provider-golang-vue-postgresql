package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LeonLow97/go-clean-architecture/delivery/http/app"
	"github.com/LeonLow97/go-clean-architecture/infrastructure"
	logger "github.com/LeonLow97/go-clean-architecture/infrastructure/logger"
	"github.com/LeonLow97/go-clean-architecture/repository"
	"github.com/LeonLow97/go-clean-architecture/usecase"
	"github.com/LeonLow97/go-clean-architecture/utils/constants/headers"
	"github.com/rs/cors"
)

func main() {
	// loading config file
	cfg, err := infrastructure.LoadConfig()
	if err != nil {
		log.Fatalln("error loading config file", err)
	}

	l, err := logger.NewZapLogger()
	if err != nil {
		log.Fatalln("error loading zap logger", err)
	}
	defer l.Sync()
	l.Info("Hello Test", l.String("error", "Testing error message"))

	// connecting to Postgres
	conn, err := infrastructure.NewPostgresConn(cfg)
	if err != nil {
		log.Fatalln("error connecting to db", err)
	}
	defer conn.Close()
	dbConn := conn.DB

	// connecting to Redis
	redisClient := infrastructure.NewRedisClient(cfg)
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Fatalln("error closing redis client", err)
		}
	}()

	// setting up SMTP instance
	smtpClient := infrastructure.NewSMTPInstance(cfg)

	// setting up TOTP instance
	totpInstance := infrastructure.NewTOTPMultiFactor(cfg)

	// Initiating handlers, service, and repository
	userRepo := repository.NewUserRepository(dbConn)
	userUsecase := usecase.NewUserUsecase(*cfg, userRepo, redisClient, *smtpClient, totpInstance)

	balanceRepo := repository.NewBalanceRepository(dbConn)
	balanceUsecase := usecase.NewBalanceUsecase(dbConn, userRepo, balanceRepo)

	beneficiaryRepo := repository.NewBeneficiaryRepository(dbConn)
	beneficiaryUsecase := usecase.NewBeneficiaryUsecase(beneficiaryRepo)

	walletRepo := repository.NewWalletRepository(dbConn)
	walletUsecase := usecase.NewWalletUsecase(dbConn, walletRepo, balanceRepo)

	transactionRepo := repository.NewTransactionRepository(dbConn)
	transactionUsecase := usecase.NewTransactionUsecase(dbConn, transactionRepo, walletRepo, balanceRepo, userRepo)

	application := &app.Application{
		Cfg:                cfg,
		RedisClient:        redisClient,
		UserUsecase:        userUsecase,
		BalanceUsecase:     balanceUsecase,
		BeneficiaryUsecase: beneficiaryUsecase,
		WalletUsecase:      walletUsecase,
		TransactionUsecase: transactionUsecase,
	}

	apiRouter, err := application.CreateRouter()
	if err != nil {
		log.Fatalf("failed to create router with error: %v\n", err)
	}

	// Create CORS options
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{
			cfg.Frontend.FrontendURL,
		},
		AllowCredentials: true,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		// AllowedHeaders specifies which headers are allowed to be sent in requests from client (browser) to server
		AllowedHeaders: []string{
			headers.Accept,
			headers.CacheControl,
			headers.ContentType,
			headers.Origin,
			headers.XCSRFToken,
		},
		// ExposedHeaders specifies which response headers are exposed to client (browser) and can be accessed by JavaScript
		ExposedHeaders: []string{
			headers.CacheControl,
			headers.ContentType,
			headers.ContentSecurityPolicy,
			headers.Pragma,
			headers.ReferrerPolicy,
			headers.StrictTransportSecurity,
			headers.XContentTypeOptions,
			headers.XCSRFToken,
			headers.XFrameOptions,
			headers.XHasNextPage,
			headers.XHasPreviousPage,
			headers.XPage,
			headers.XPageSize,
			headers.XTotal,
			headers.XTotalPages,
		},
		MaxAge:             86400, // cache HTTP headers set by CORS for 1 day (86400 seconds)
		OptionsPassthrough: false,
		Debug:              false,
	})
	wrappedRouter := corsHandler.Handler(apiRouter)

	port := cfg.Server.Port
	log.Println("Server is running on port", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), wrappedRouter); err != nil {
		log.Fatalf("Failed to listen to server with error: %v\n", err)
	}
}
