package repository

import (
	"belajar-golang-rest-api/lat/domain"
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
)

type CustomerRepository struct {
	db *goqu.Database
}

func NewCustomer(con *sql.DB) domain.CustomerRepository {
	return &CustomerRepository{
		db: goqu.New("default", con),
	}
}

// Delete implements domain.CustomerRepository.
func (cr *CustomerRepository) Delete(ctx context.Context, id string) error {
	executor := cr.db.Update("customers").Where(goqu.C("id").Eq(id)).Set(goqu.Record{"deleted_at": sql.NullTime{Valid: true, Time: time.Now()}}).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}

// FindAll implements domain.CustomerRepository.
func (cr CustomerRepository) FindAll(ctx context.Context) (result []domain.Customer, err error) {
	dataset := cr.db.From("customers").Where(goqu.C("deleted_at").IsNull())
	err = dataset.ScanStructsContext(ctx, &result)
	return
}

// FindById implements domain.CustomerRepository.
func (cr *CustomerRepository) FindById(ctx context.Context, id string) (result domain.Customer, err error) {
	dataset := cr.db.From("customers").Where(goqu.C("deleted_at").IsNull(), goqu.C("id").Eq(id))
	_, err = dataset.ScanStructContext(ctx, &result)
	return
}

func (cr *CustomerRepository) FindByIds(ctx context.Context, id []string) (result []domain.Customer, err error) {
	dataset := cr.db.From("customers").Where(goqu.C("deleted_at").IsNull(), goqu.C("id").Eq(id))
	err = dataset.ScanStructsContext(ctx, &result)
	return
}

// Save implements domain.CustomerRepository.
func (cr CustomerRepository) Save(ctx context.Context, c *domain.Customer) error {
	executor := cr.db.Insert("customers").Rows(c).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}

// Update implements domain.CustomerRepository.
func (cr CustomerRepository) Update(ctx context.Context, c *domain.Customer) error {
	executor := cr.db.Update("customers").Where(goqu.C("id").Eq(c.ID)).Set(c).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}
