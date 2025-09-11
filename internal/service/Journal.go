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

// Journal service handles journal operations
type JournalService struct {
	journalRepository   domain.JournalRepository
	bookRepository      domain.BookRepository
	bookStockRepository domain.BookStockRepository
	customerRepository domain.CustomerRepository
	chargeRepository domain.ChargeRepository
}

// Create implements domain.JournalService.


// NewJournalService creates a new journal service
func NewJournal(journalRepository domain.JournalRepository,
	bookRepository domain.BookRepository,
	bookStockRepository domain.BookStockRepository,
	customerRepository domain.CustomerRepository,
	chargeRepository domain.ChargeRepository) domain.JournalService {
	return &JournalService{
		journalRepository:   journalRepository,
		bookRepository:      bookRepository,
		bookStockRepository: bookStockRepository,
		customerRepository: customerRepository,
		chargeRepository: chargeRepository,
	}
}
// Index implements domain.JournalService.
func (j *JournalService) Index(ctx context.Context, se domain.JournalSearch) ([]dto.JournalData, error) {
	journals, err := j.journalRepository.Find(ctx,se)
	if err != nil{
		return nil ,err
	}
	customerId:= make([]string,0)
	bookId:= make([]string, 0)
	for _, v := range journals{
		customerId = append(customerId, v.CustomerId)
		bookId = append(bookId, v.BookId)
	}
	customers := make(map[string]domain.Customer)
	if len(customerId) > 0 {
		customerDb, _ := j.customerRepository.FindByIds(ctx, customerId)
		for _, v := range customerDb{
			customers[v.ID]= v
		}
	}  
	books := make(map[string]domain.Book)
	if len(bookId) > 0 {
		bookDb, _ := j.bookRepository.FindByIds(ctx, bookId)
		for _, v := range bookDb{
			books[v.Id] =v
		}
	}

	result := make([]dto.JournalData, 0)
	for _, v :=  range journals {
			book := dto.BookData{}
				if v2, e := books[v.BookId]; e{
					book = dto.BookData{
						Id :v2.Id,
						Isbn: v2.Isbn,
						Title :v2.Title,
						Description : v2.Description,
					}
				}
				customer := dto.CustomerData{}
				if v2, e := customers[v.CustomerId]; e {
					customer =dto.CustomerData{
						ID: v2.ID,
						Code: v2.Code,
						Name: v2.Name,
					}
				}
				result = append(result, dto.JournalData{
					Id: v.Id,
					BookStock: v.StockCode,
					Book: book,
					Customer: customer,
					Status: v.Status,
					BorrowedAt : v.BorrowedAt.Time,
					ReturnedAt : v.ReturnedAt.Time,
				})
			}
			return result, nil
}

func (j *JournalService) Create(ctx context.Context, req dto.CreateJournalRequest) error {
	book, err := j.bookRepository.FindById(ctx, req.BookId)
	if err != nil{
		return err
	}
	if book.Id == ""{
		return domain.BookNotFound
	}
	stock, err := j.bookStockRepository.FindByBookAndCode(ctx, book.Id, req.BookStock)
	if err != nil{
		return err
	}

	if stock.Code == "" {
		return domain.BookNotFound
	}
	if stock.Status != domain.BookStockStatusAvailable{
		return errors.New("stock already borrowed")
	}

	journal := domain.Journal{
		Id : uuid.NewString(),
		BookId: req.BookId,
		StockCode: req.BookStock,
		CustomerId: req.CustomerId,
		Status: domain.JournalStatusInProgress,
		DueAt: sql.NullTime{Valid: true, Time: time.Now().Add(7 * 24 *time.Hour)},
		BorrowedAt: sql.NullTime{Valid : true, Time: time.Now()},
	} 
	err = j.journalRepository.Save(ctx, &journal)
	if err != nil{
		return err
	}
	stock.Status = domain.BookStockStatusBorrowed
	stock.BorrowedAt = journal.BorrowedAt
	stock.BorrowerId = sql.NullString{Valid: true, String: journal.CustomerId} 
	return j.bookStockRepository.Update(ctx, &stock)

}



// Return implements domain.JournalService.
func (j *JournalService) Return(ctx context.Context, req dto.ReturnJournalRequest) error {
	journal, err := j.journalRepository.FindById(ctx, req.JournalId)
	if err != nil{
		return err
	}
	if journal.Id == ""{
		return domain.JournalNotFound
	}

	stock, err := j.bookStockRepository.FindByBookAndCode(ctx, journal.BookId, journal.StockCode)
	if err != nil{
		return err
	}
	if stock.Code != ""{
		stock.Status = domain.BookStockStatusAvailable
		stock.BorrowerId = sql.NullString{Valid: false}
		stock.BorrowedAt= sql.NullTime{Valid: false}
		err = j.bookStockRepository.Update(ctx, &stock)
		if err != nil{
			return err
		}
		
	}
	journal.Status = domain.JournalStatusCompleted
	journal.ReturnedAt = sql.NullTime{Valid: true, Time: time.Now()}
	err =  j.journalRepository.Update(ctx, &journal)
	if err != nil{
			return err
		}
	
	hoursLate := time.Now().Sub(journal.DueAt.Time).Hours()
	if hoursLate >= 24 {
		daysLate := int(hoursLate/24)
		dailyLateFee := 5000
		charge := domain.Charge{
			Id: uuid.NewString(),
			JournalId: journal.Id,
			DaysLate: daysLate,
			DailyLateFee: dailyLateFee,
			Total: dailyLateFee * daysLate,
			UserId: req.UserId,
			CreatedAt: sql.NullTime{Valid: true, Time: time.Now()},
		}
		err = j.chargeRepository.Save(ctx, &charge)
	}
	return err
}