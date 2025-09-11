package repository

import (
	"belajar-golang-rest-api/lat/domain"
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
)

type chargeRepository struct {
	db *goqu.Database
}

func NewCharge(con *sql.DB) domain.ChargeRepository {
	return &chargeRepository{
		db: goqu.New("defaul", con),
	}
}

func (c *chargeRepository) Save(ctx context.Context, charge *domain.Charge) error {
	executor := c.db.Insert("charges").Rows(charge).Executor()
	_, err := executor.ExecContext(ctx)
	return err 
}