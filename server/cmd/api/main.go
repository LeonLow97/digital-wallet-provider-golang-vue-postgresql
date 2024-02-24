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
)

func main() {
	conn, err := infrastructure.NewPostgresConn()
	if err != nil {
		log.Fatalln("error connecting to db", err)
	}
	defer conn.Close()
	dbConn := conn.DB

	redisClient := infrastructure.NewRedisClient()
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Fatalln("error closing redis client", err)
		}
	}()

	router := mux.NewRouter()
	router = router.PathPrefix("/api/v1").Subrouter() // api versioning v1

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Healthy!"))
	}).Methods(http.MethodGet)

	// Initiating handlers, service, and repository
	userRepo := repository.NewUserRepository(dbConn)
	authUsecase := usecase.NewAuthUsecase(userRepo, redisClient)
	handlers.NewAuthHandler(router, authUsecase, redisClient)

	balanceRepo := repository.NewBalanceRepository(dbConn)
	balanceUsecase := usecase.NewBalanceUsecase(balanceRepo)
	handlers.NewBalanceHandler(router, balanceUsecase)

	walletRepo := repository.NewWalletRepository(dbConn)
	walletUsecase := usecase.NewWalletUsecase(walletRepo, balanceRepo)
	handlers.NewWalletHandler(router, walletUsecase)

	transactionRepo := repository.NewTransactionRepository(dbConn)
	transactionUsecase := usecase.NewTransactionUsecase(transactionRepo, walletRepo, userRepo, balanceRepo)
	handlers.NewTransactionHandler(router, transactionUsecase)

	beneficiaryRepo := repository.NewBeneficiaryRepository(dbConn)
	beneficiaryUsecase := usecase.NewBeneficiaryUsecase(beneficiaryRepo)
	handlers.NewBeneficiaryHandler(router, beneficiaryUsecase)

	// skipping endpoints
	skipperFunc := middleware.NewSkipperFunc(
		"/api/v1/login",
		"/api/v1/signup",
		"/api/v1/health",
	)

	router.Use(
		middleware.NewCorsMiddleware,
		middleware.NewAuthenticationMiddleware(skipperFunc, redisClient, authUsecase).Middleware,
	)

	port := os.Getenv("SERVICE_PORT")
	log.Println("Server is running on port", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		log.Fatalf("Failed to listen to server with error: %v\n", err)
	}
}

// connectToDB connects to PostgreSQL
// func connectToDB() (*sqlx.DB, error) {
// 	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
// 		os.Getenv("POSTGRES_USER"),
// 		os.Getenv("POSTGRES_PASSWORD"),
// 		os.Getenv("POSTGRES_HOST"),
// 		os.Getenv("POSTGRES_PORT"),
// 		os.Getenv("POSTGRES_DB"),
// 	)

// 	connConfig, err := pgx.ParseConfig(dsn)
// 	if err != nil {
// 		errMsg := err.Error()
// 		errMsg = regexp.MustCompile(`(://[^:]+:).+(@.+)`).ReplaceAllString(errMsg, "$1*****$2")
// 		errMsg = regexp.MustCompile(`(password=).+(\s+)`).ReplaceAllString(errMsg, "$1*****$2")
// 		return nil, fmt.Errorf("parsing DSN failed: %s", errMsg)
// 	}
// 	connStr := stdlib.RegisterConnConfig(connConfig)
// 	db, err := sqlx.Open("pgx", connStr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := db.Ping(); err != nil {
// 		return nil, err
// 	}

// 	return db, nil
// }
