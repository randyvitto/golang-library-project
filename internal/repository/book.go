package repository

import (
	"belajar-golang-rest-api/lat/domain"
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
)

type bookRepository struct {
	db *goqu.Database
}





func NewBook(con *sql.DB) domain.BookRepository {
	return &bookRepository{
		db: goqu.New("default", con),
	}
}
// FindAll implements domain.BookRepository.
func (b bookRepository) FindAll(ctx context.Context) (books []domain.Book, err error) {
	dataset := b.db.From("books").Where(goqu.C("deleted_at").IsNull())
	err = dataset.ScanStructsContext(ctx, &books)
	return
}

// FindById implements domain.BookRepository.
func (b bookRepository) FindById(ctx context.Context, id string) (book domain.Book, err error) {
	dataset := b.db.From("books").
	Where(
		goqu.C("id").Eq(id),
		 goqu.C("deleted_at").IsNull(),)
		 _, err = dataset.ScanStructContext(ctx, &book)
		 return
}

func (b bookRepository) FindByIds(ctx context.Context, ids []string) (books []domain.Book, err error) {
	dataset := b.db.From("books").
	Where(
		goqu.C("id").Eq(ids),
		 goqu.C("deleted_at").IsNull(),)
		 err = dataset.ScanStructsContext(ctx, &books)
		 return
}

// Save implements domain.BookRepository.
func (b bookRepository) Save(ctx context.Context, book *domain.Book) error {
	executor := b.db.Insert("books").Rows(book).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}

// Update implements domain.BookRepository.
func (b bookRepository) Update(ctx context.Context, book *domain.Book) error {
	executor := b.db.Update("books").
	Where(goqu.C("id").Eq(book.Id)).Set(book).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}

// Delete implements domain.BookRepository.
func (b bookRepository) Delete(ctx context.Context, id string) error {
	executor := b.db.Update("books").Where(goqu.C("id").Eq(id)).Set(goqu.Record{"deleted_at": sql.NullTime{Valid : true, Time: time.Now()}}).
	Executor()
	_, err := executor.ExecContext(ctx)
	return err
}