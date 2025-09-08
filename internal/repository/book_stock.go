package repository

import (
	"belajar-golang-rest-api/lat/domain"
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
)

type bookStockRepository struct {
	db *goqu.Database
}


func NewBookStock(con *sql.DB) domain.BookStockRepository {
	return &bookStockRepository{
		db: goqu.New("default", con),
	}
}

// FindByBookId implements domain.BookStockRepository.
func (b *bookStockRepository) FindByBookId(ctx context.Context, id string) (result []domain.BookStock, err error) {
	dataset := b.db.From("book_stocks").Where(goqu.C("book_id").Eq(id))
	err = dataset.ScanStructsContext(ctx, &result)
	return
}

// FindByBookAndCode implements domain.BookStockRepository.
func (b *bookStockRepository) FindByBookAndCode(ctx context.Context, id string, code string) (result domain.BookStock, err error) {
	dataset := b.db.From("book_stocks").Where(goqu.C("book_id").Eq(id), goqu.C("code").Eq(code))
	_, err = dataset.ScanStructContext(ctx, &result)
	return
}

// Save implements domain.BookStockRepository.
func (b *bookStockRepository) Save(ctx context.Context, data []domain.BookStock) error {
	executor := b.db.Insert("book_stocks").Rows(data).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}

// Update implements domain.BookStockRepository.
func (b *bookStockRepository) Update(ctx context.Context, stock *domain.BookStock) error {
	executor := b.db.Update("book_stocks").
	Where(goqu.C("code").Eq(stock.Code)).Set(stock).
	Executor()
	_, err := executor.ExecContext(ctx)
	return err
}
// DeleteByBookId implements domain.BookStockRepository.
func (b *bookStockRepository) DeleteByBookId(ctx context.Context, id string) error {
	executor := b.db.Delete("book_stocks").
	Where(goqu.C("book_id").Eq(id)).
	Executor()
	_, err := executor.ExecContext(ctx)
	return err
}

// DeleteBycodes implements domain.BookStockRepository.
func (b *bookStockRepository) DeleteBycodes(ctx context.Context, codes []string) error {
	executor := b.db.Delete("book_stocks").
	Where(goqu.C("code").In(codes)).
	Executor()
	_, err := executor.ExecContext(ctx)
	return err
}


