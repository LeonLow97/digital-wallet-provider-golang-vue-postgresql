package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	handlers "github.com/LeonLow97/go-clean-architecture/delivery/http/handler"
	"github.com/LeonLow97/go-clean-architecture/delivery/http/middleware"
	"github.com/LeonLow97/go-clean-architecture/infrastructure"
	"github.com/LeonLow97/go-clean-architecture/repository"
	"github.com/LeonLow97/go-clean-architecture/usecase"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// connecting to Postgres
	conn, err := infrastructure.NewPostgresConn()
	if err != nil {
		log.Fatalln("error connecting to db", err)
	}
	defer conn.Close()
	dbConn := conn.DB

	// connecting to Redis
	redisClient := infrastructure.NewRedisClient()
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Fatalln("error closing redis client", err)
		}
	}()

	// loading config file
	cfg, err := infrastructure.LoadConfig()
	if err != nil {
		log.Fatalln("error loading config file", err)
	}

	router := mux.NewRouter()
	router = router.PathPrefix("/api/v1").Subrouter() // api versioning v1

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Healthy!"))
	}).Methods(http.MethodGet)

	// setting up SMTP instance
	smtpClient := infrastructure.NewSMTPInstance(cfg)

	// setting up TOTP instance
	totpInstance := infrastructure.NewTOTPMultiFactor(cfg)

	// Initiating handlers, service, and repository
	userRepo := repository.NewUserRepository(dbConn)
	userUsecase := usecase.NewUserUsecase(*cfg, userRepo, redisClient, *smtpClient, totpInstance)
	handlers.NewUserHandler(router, userUsecase, redisClient)

	balanceRepo := repository.NewBalanceRepository(dbConn)
	balanceUsecase := usecase.NewBalanceUsecase(dbConn, userRepo, balanceRepo)
	handlers.NewBalanceHandler(router, balanceUsecase)

	walletRepo := repository.NewWalletRepository(dbConn)
	walletUsecase := usecase.NewWalletUsecase(dbConn, walletRepo, balanceRepo)
	handlers.NewWalletHandler(router, walletUsecase)

	transactionRepo := repository.NewTransactionRepository(dbConn)
	transactionUsecase := usecase.NewTransactionUsecase(dbConn, transactionRepo, walletRepo, userRepo, balanceRepo)
	handlers.NewTransactionHandler(router, transactionUsecase)

	beneficiaryRepo := repository.NewBeneficiaryRepository(dbConn)
	beneficiaryUsecase := usecase.NewBeneficiaryUsecase(beneficiaryRepo)
	handlers.NewBeneficiaryHandler(router, beneficiaryUsecase)

	// skipping endpoints
	skipperFunc := middleware.NewSkipperFunc(
		"/api/v1/login",
		"/api/v1/password-reset/send",
		"/api/v1/password-reset/reset",
		"/api/v1/signup",
		"/api/v1/health",
		"/api/v1/configure-mfa",
		"/api/v1/verify-mfa",
	)

	router.Use(
		// TODO: Add SessionMiddleware to inject user object and session details into context
		middleware.NewAuthenticationMiddleware(*cfg, skipperFunc, redisClient, userUsecase).Middleware,
		middleware.NewCSRFMiddleware(*cfg, skipperFunc, redisClient).Middleware,
		middleware.NewCSPMiddleware(),
	)

	// Create CORS options
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{
			cfg.Env.BackendURL,
			cfg.Env.FrontendURL,
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
		AllowedHeaders: []string{"Accept", "Origin", "Content-Type", "Authorization", "X-CSRF-Token"},
		// ExposedHeaders specifies which response headers are exposed to client (browser) and can be accessed by JavaScript
		ExposedHeaders:     []string{"X-CSRF-Token", "Content-Security-Policy"},
		MaxAge:             86400, // cache for 1 day (86400 seconds)
		OptionsPassthrough: false,
		Debug:              false,
	})
	wrappedRouter := corsHandler.Handler(router)

	port := os.Getenv("SERVICE_PORT")
	log.Println("Server is running on port", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), wrappedRouter); err != nil {
		log.Fatalf("Failed to listen to server with error: %v\n", err)
	}
}
