package dto

type BookData struct {
	Id          string `json:"id"`
	Isbn        string `json:"isbn"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type BookStockData struct {
	Code   string `json:"code"`
	Status string `json:"status"`
}

type BookShowData struct {
	BookData
	Stock []BookStockData
}

type CreateBookRequest struct {
	Isbn        string `json:"isbn" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateBookRequest struct {
	Id          string `json:"-"`
	Isbn        string `json:"isbn" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}
