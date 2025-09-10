package service

import (
	"belajar-golang-rest-api/lat/domain"
	"belajar-golang-rest-api/lat/dto"
	"belajar-golang-rest-api/lat/internal/config"
	"context"
	"database/sql"
	"errors"
	"path"
	"time"

	"github.com/google/uuid"
)

type BookService struct {
	cnf                 *config.Config
	bookRepository      domain.BookRepository
	bookStockRepository domain.BookStockRepository
	mediaRepository     domain.MediaRepository
}

func NewBook(cnf *config.Config,
	bookRepository domain.BookRepository,
	bookStockRepository domain.BookStockRepository,
	mediaRepository domain.MediaRepository) domain.BookService {
	return &BookService{
		cnf:                 cnf,
		bookRepository:      bookRepository,
		bookStockRepository: bookStockRepository,
		mediaRepository:     mediaRepository,
	}
}

// Index implements domain.BookService.
func (b *BookService) Index(ctx context.Context) ([]dto.BookData, error) {
	result, err := b.bookRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	coverId := make([]string, 0)
	for _, v := range result {
		if v.CoverId.Valid {
			coverId = append(coverId, v.CoverId.String)
		}
	}
	covers := make(map[string]string)
	if len(coverId) > 0 {
		coversDb, _ := b.mediaRepository.FindByIds(ctx, coverId)
		for _, v := range coversDb {
			covers[v.Id] = path.Join(b.cnf.Server.Asset, v.Path)
		}
	}

	var data []dto.BookData
	for _, v := range result {
		var coverUrl string
		if v2, e := covers[v.CoverId.String]; e {
			coverUrl = v2
		}
		data = append(data, dto.BookData{
			Id:          v.Id,
			Isbn:        v.Isbn,
			Title:       v.Title,
			CoverUrl:    coverUrl,
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
	stocks, err := b.bookStockRepository.FindByBookId(ctx, data.Id)
	if err != nil {
		return dto.BookShowData{}, err
	}

	stocksData := make([]dto.BookStockData, 0)
	for _, v := range stocks {
		stocksData = append(stocksData, dto.BookStockData{
			Code:   v.Code,
			Status: v.Status,
		})
	}

	var coverUrl string
	if data.CoverId.Valid {
		cover, _ := b.mediaRepository.FindById(ctx, data.CoverId.String)
		if cover.Path != "" {
			coverUrl = path.Join(b.cnf.Server.Asset, cover.Path)
		}
	}

	return dto.BookShowData{
		BookData: dto.BookData{
			Id:          data.Id,
			Isbn:        data.Isbn,
			Title:       data.Title,
			CoverUrl:    coverUrl,
			Description: data.Description,
		},
		Stock: stocksData,
	}, nil
}

func (b *BookService) Create(ctx context.Context, req dto.CreateBookRequest) error {
	coverId := sql.NullString{Valid: false, String: req.CoverId}
	if req.CoverId != "" {
		coverId.Valid = true
	}

	book := domain.Book{
		Id:          uuid.NewString(),
		Isbn:        req.Isbn,
		Title:       req.Title,
		Description: req.Description,
		CoverId:     coverId,
		CreatedAt:   sql.NullTime{Valid: true, Time: time.Now()},
	}
	return b.bookRepository.Save(ctx, &book)
}

// Update implements domain.BookService.
func (b *BookService) Update(ctx context.Context, req dto.UpdateBookRequest) error {
	exist, err := b.bookRepository.FindById(ctx, req.Id)
	if err != nil {
		return err
	}
	if exist.Id == "" {
		return errors.New("Book not Found")
	}
	coverId := sql.NullString{Valid: false, String: req.CoverId}
	if req.CoverId != "" {
		coverId.Valid = true
	}
	exist.Isbn = req.Isbn
	exist.Title = req.Title
	exist.Description = req.Description
	exist.CoverId = coverId
	exist.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}
	return b.bookRepository.Update(ctx, &exist)

}

func (b *BookService) Delete(ctx context.Context, id string) error {
	exist, err := b.bookRepository.FindById(ctx, id)
	if err != nil {
		return err
	}
	if exist.Id == "" {
		return errors.New("Book not Found")
	}
	err = b.bookRepository.Delete(ctx, exist.Id)
	if err != nil {
		return err
	}
	return b.bookStockRepository.DeleteByBookId(ctx, exist.Id)
}
