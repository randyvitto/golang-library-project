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

type CustomerService struct {
	customerRepository domain.CustomerRepository
}



// Create implements domain.CustomerService.

func NewCustomer(customerRepository domain.CustomerRepository) domain.CustomerService {
	return CustomerService{
		customerRepository: customerRepository,
	}
}

// Create implements domain.CustomerService.
func (c CustomerService) Create(ctx context.Context, req dto.CreateCustomerRequest) error {
	customer := domain.Customer{
		ID:         uuid.NewString(),
		Code:       req.Code,
		Name:       req.Name,
		Created_at: sql.NullTime{Valid: true, Time: time.Now()},
	}
	return c.customerRepository.Save(ctx, &customer)
}

// Index implements domain.CustomerService.
func (c CustomerService) Index(ctx context.Context) ([]dto.CustomerData, error) {
	customer, err := c.customerRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	var customerData []dto.CustomerData
	for _, v := range customer {
		customerData = append(customerData, dto.CustomerData{
			ID:   v.ID,
			Code: v.Code,
			Name: v.Name,
		})
	}
	return customerData, nil
}

// Update implements domain.CustomerService.
func (c CustomerService) Update(ctx context.Context, req dto.UpdateCustomerRequest) error {
	exist, err := c.customerRepository.FindById(ctx, req.ID)
	if err != nil {
		return err
	}
	if exist.ID == "" {
		return errors.New("Customer not found")
	}
	exist.Code = req.Code
	exist.Name = req.Name
	exist.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}

	return c.customerRepository.Update(ctx, &exist)
}

// Delete implements domain.CustomerService.
func (c CustomerService) Delete(ctx context.Context, id string) error {
	exist, err := c.customerRepository.FindById(ctx, id)
	if err != nil {
		return err
	}
	if exist.ID == "" {
		return errors.New("Customer not found")
	}
	return c.customerRepository.Delete(ctx, id)
}

// Show implements domain.CustomerService.
func (c CustomerService) Show(ctx context.Context, id string) (dto.CustomerData, error) {
	exist, err := c.customerRepository.FindById(ctx, id)
	if err != nil{
		return dto.CustomerData{} , err
	}
	if exist.ID == ""{
		return dto.CustomerData{}, errors.New("Customer data not found")
	}
	return dto.CustomerData{
		ID: exist.Code,
		Code: exist.Code,
		Name: exist.Name,
	}, nil
}