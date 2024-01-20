package repository

import (
	"github.com/LeonLow97/go-clean-architecture/domain"
	"github.com/jmoiron/sqlx"
)

type balanceRepository struct {
	db *sqlx.DB
}

func NewBalanceRepository(db *sqlx.DB) domain.BalanceRepository {
	return &balanceRepository{
		db: db,
	}
}
