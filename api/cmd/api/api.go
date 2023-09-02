package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Open Database Connection
	db, err := openDB()
	if err != nil {
		log.Fatalln(err)
	}

	router := routes(db)

	log.Println("Application has started. Listening port is 4000")
	http.ListenAndServe(":4000", router)
}

func openDB() (*sqlx.DB, error) {
	connConfig, err := pgx.ParseConfig("postgres://postgres:postgres@db:5432/leon-db?sslmode=disable")
	if err != nil {
		errMsg := err.Error()
		errMsg = regexp.MustCompile(`(://[^:]+:).+(@.+)`).ReplaceAllString(errMsg, "$1*****$2")
		errMsg = regexp.MustCompile(`(password=).+(\s+)`).ReplaceAllString(errMsg, "$1*****$2")
		return nil, fmt.Errorf("parsing DSN failed: %s", errMsg)
	}
	connectionStr := stdlib.RegisterConnConfig(connConfig)
	db, err := sqlx.Open("pgx", connectionStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	instance, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"leon-db",
		instance,
	)
	if err != nil {
		return nil, err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	seed(db)

	return db, nil
}

func seed(db *sqlx.DB) {
	users := []struct {
		Username     string
		Password     string
		Active       int
		Admin        int
		MobileNumber string
	}{
		{Username: "Alice", Password: "$2a$10$CerQd299qowq2ck8k/EqQeB7Jpjd/4Cut/Df.f8jnq9kYsuG0W7zG", Active: 1, Admin: 1, MobileNumber: "+65 90399012"},
		{Username: "Bob", Password: "$2a$10$MVLL5BT/nIQKk6OYbgzK7.fbT0XKMBtNdeoy64ihYUUhr8Ag6358u", Active: 1, Admin: 1, MobileNumber: "+65 89230122"},
		{Username: "Charlie", Password: "$2a$10$yKz0rguTzykTec4Bgke7LempFl/GQVTw9w9qEXfGUpI/XGK97VHFq", Active: 1, Admin: 1, MobileNumber: "+1 555-123-4567"},
		{Username: "David", Password: "$2a$10$p444biF49.py2HOTVe5TSuUNAhSKqelEtlbLtZXghUh3o21Et7DNO", Active: 1, Admin: 1, MobileNumber: "+49 1234567890"},
	}
	for _, user := range users {
		db.Exec("INSERT INTO users(username, password, active, admin, mobile_number) VALUES ($1,$2,$3,$4,$5) ON CONFLICT DO NOTHING", user.Username, user.Password, user.Active, user.Admin, user.MobileNumber)
	}

	userBalances := []struct {
		UserId         int
		Balance        float64
		Currency       string
		CountryIsoCode string
	}{
		{UserId: 1, Balance: 70000.00, Currency: "SGD", CountryIsoCode: "SG"},
		{UserId: 2, Balance: 3000.00, Currency: "SGD", CountryIsoCode: "SG"},
		{UserId: 3, Balance: 6000.00, Currency: "USD", CountryIsoCode: "US"},
		{UserId: 4, Balance: 2000.00, Currency: "EUR", CountryIsoCode: "FR"},
	}
	for _, userBalance := range userBalances {
		db.Exec("INSERT INTO user_balance (user_id, balance, currency, country_iso_code) VALUES ($1,$2,$3,$4) ON CONFLICT DO NOTHING", userBalance.UserId, userBalance.Balance, userBalance.Currency, userBalance.CountryIsoCode)
	}

	beneficiaries := []struct {
		BeneficiaryName string
		MobileNumber    string
		Currency        string
		IsInternal      int
	}{
		{BeneficiaryName: "Alice", MobileNumber: "+65 90399012", Currency: "SGD", IsInternal: 1},
		{BeneficiaryName: "Bob", MobileNumber: "+65 89230122", Currency: "SGD", IsInternal: 1},
		{BeneficiaryName: "Charlie", MobileNumber: "+1 555-123-4567", Currency: "USD", IsInternal: 1},
		{BeneficiaryName: "David", MobileNumber: "+49 1234567890", Currency: "EUR", IsInternal: 1},
	}
	for _, beneficiary := range beneficiaries {
		db.Exec("INSERT INTO beneficiaries (beneficiary_name, mobile_number, currency, is_internal) VALUES ($1,$2,$3,$4) ON CONFLICT DO NOTHING", beneficiary.BeneficiaryName, beneficiary.MobileNumber, beneficiary.Currency, beneficiary.IsInternal)
	}

	userBeneficiaryMapping := []struct {
		UserId        int
		BeneficiaryId int
	}{
		{UserId: 1, BeneficiaryId: 2},
		{UserId: 1, BeneficiaryId: 3},
		{UserId: 1, BeneficiaryId: 4},
		{UserId: 2, BeneficiaryId: 1},
		{UserId: 2, BeneficiaryId: 3},
		{UserId: 3, BeneficiaryId: 2},
		{UserId: 3, BeneficiaryId: 4},
		{UserId: 4, BeneficiaryId: 1},
	}
	for _, userBeneficiary := range userBeneficiaryMapping {
		db.Exec("INSERT INTO user_beneficiary (user_id, beneficiary_id) VALUES ($1,$2) ON CONFLICT DO NOTHING", userBeneficiary.UserId, userBeneficiary.BeneficiaryId)
	}

	log.Println("Database Seeding Completed.")
}
