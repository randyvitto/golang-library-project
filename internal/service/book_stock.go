package service

import (
	"belajar-golang-rest-api/lat/domain"
	"belajar-golang-rest-api/lat/dto"
	"context"
)

type bookStockService struct {
	bookRepository      domain.BookRepository
	bookStockRepository domain.BookStockRepository
}



func NewBookStock(bookRepository domain.BookRepository,
	bookStockrRepository domain.BookStockRepository) domain.BookStockService {
	return &bookStockService{
		bookRepository:      bookRepository,
		bookStockRepository: bookStockrRepository,
	}
}

// Create implements domain.BookStockService.
func (b *bookStockService) Create(ctx context.Context, req dto.CreateBookStockRequest) error {
	book, err := b.bookRepository.FindById(ctx, req.BookId)
	if err != nil{
		return err
	}
	if book.Id == ""{
		return domain.BookNotFound
	}
	

	stocks := make([]domain.BookStock, 0 )
	for _, v := range req.Codes {
		stocks = append(stocks, domain.BookStock{
			Code: v,
			BookId: req.BookId,
			Status: domain.BookStockStatusAvailable,
		})
	}
	return b.bookStockRepository.Save(ctx, stocks)
}

// Delete implements domain.BookStockService.
func (b *bookStockService) Delete(ctx context.Context, req dto.DeleteBookStockRequest) error {
	return b.bookStockRepository.DeleteBycodes(ctx, req.Codes )
}