package domain

import (
	"belajar-golang-rest-api/lat/dto"
	"context"
	"database/sql"
)

type Book struct {
	Id          string       `db:"id"`
	Isbn        string       `db:"isbn"`
	Title       string       `db:"title"`
	Description string       `db:"description"`
	CreatedAt  sql.NullTime `db:"created_at"`
	UpdatedAt   sql.NullTime `db:"updated_at"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
}

type BookRepository interface {
	FindAll(ctx context.Context) ([]Book, error)
	FindById(ctx context.Context, id string) (Book, error)
	FindByIds(ctx context.Context, ids []string) ([]Book, error)
	Save(ctx context.Context, book *Book) error
	Update(ctx context.Context, book *Book) error
	Delete(ctx context.Context, id string) error
}

type BookService interface {
	Index(ctx context.Context) ([]dto.BookData, error)
	Show(ctx context.Context, id string) (dto.BookShowData, error)
	Create(ctx context.Context, req dto.CreateBookRequest) error
	Update(ctx context.Context, req dto.UpdateBookRequest) error
	Delete(ctx context.Context, id string) error
}
