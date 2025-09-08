package service

import (
	"belajar-golang-rest-api/lat/domain"
	"belajar-golang-rest-api/lat/dto"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type BookService struct {
	bookRepository domain.BookRepository
	bookStockRepository domain.BookStockRepository
}

func NewBook(bookRepository domain.BookRepository,
	bookStockRepository domain.BookStockRepository) domain.BookService {
	return &BookService{
		bookRepository:      bookRepository,
		bookStockRepository: bookStockRepository,
	}
}

// Index implements domain.BookService.
func (b *BookService) Index(ctx context.Context) ([]dto.BookData, error) {
	result, err := b.bookRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	var data []dto.BookData
	for _, v := range result {
		data = append(data, dto.BookData{
			Id:          v.Id,
			Isbn:        v.Isbn,
			Title:       v.Title,
			Description: v.Description,
		})
	}
	return data, nil
}

// Show implements domain.BookService.
func (b *BookService) Show(ctx context.Context, id string) (dto.BookShowData, error) {
	data, err := b.bookRepository.FindById(ctx, id)
	if err != nil {
		return dto.BookShowData{}, err
	}
	if data.Id == "" {
		return dto.BookShowData{}, domain.BookNotFound
	}
	stocks, err:= b.bookStockRepository.FindByBookId(ctx, data.Id)
	if err != nil{
		return dto.BookShowData{} ,err
	}

	stocksData := make([]dto.BookStockData, 0)
	for _, v := range stocks{
		stocksData = append(stocksData, dto.BookStockData{
			Code: v.Code,
			Status: v.Status,
		} )
	}
	return dto.BookShowData{
		BookData: dto.BookData{
		Id:          data.Id,
		Isbn:        data.Isbn,
		Title:       data.Title,
		Description: data.Description,
		},
		Stock: stocksData,
	}, nil
}

func (b *BookService) Create(ctx context.Context, req dto.CreateBookRequest) error {
	book := domain.Book{
		Id: uuid.NewString(),
		Isbn: req.Isbn,
		Title: req.Title,
		Description: req.Description,
		CreatedAt: sql.NullTime{Valid: true, Time: time.Now()},
	}
	return b.bookRepository.Save(ctx, &book)
}

// Update implements domain.BookService.
func (b *BookService) Update(ctx context.Context, req dto.UpdateBookRequest) error {
	exist, err := b.bookRepository.FindById(ctx, req.Id)
	if err != nil{
		return err
	}
	if exist.Id == "" {
		return errors.New("Book not Found") 
	}
	exist.Isbn = req.Isbn
	exist.Title = req.Title
	exist.Description = req.Description
	exist.UpdatedAt = sql.NullTime{Valid:true, Time: time.Now()}
	return b.bookRepository.Update(ctx, &exist)
	
}

func (b *BookService) Delete(ctx context.Context, id string) error {
	exist, err := b.bookRepository.FindById(ctx,id)
	if err != nil{
		return err
	}
	if exist.Id == "" {
		return errors.New("Book not Found") 
	}
	err = b.bookRepository.Delete(ctx, exist.Id)
	if err != nil{
		return err
	}
	return b.bookStockRepository.DeleteByBookId(ctx, exist.Id)
}
