package repository

import (
	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/jmoiron/sqlx"
)

type transactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) domain.TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}
