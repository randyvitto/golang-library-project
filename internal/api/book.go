package api

import (
	"belajar-golang-rest-api/lat/domain"
	"belajar-golang-rest-api/lat/dto"
	"belajar-golang-rest-api/lat/internal/util"
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type bookApi struct {
	bookService domain.BookService
}

func NewBook(app *fiber.App, 
	bookService domain.BookService, 
	authzMidd fiber.Handler){

		ba := bookApi{
			bookService: bookService,
		}
		book := app.Group("/books" ,authzMidd)

		book.Get("",  ba.Index)
		book.Post("" , ba.Create)
		book.Get(":id",  ba.Show)
		book.Put(":id",  ba.Update)
		book.Delete(":id",  ba.Delete )
}

func (ba bookApi) Index(ctx *fiber.Ctx) error{
	c, cancel := context.WithTimeout(ctx.Context(), 10 * time.Second)
	defer cancel()

	res, err := ba.bookService.Index(c)
	if err != nil{
		return ctx.Status(http.StatusInternalServerError).
		JSON(dto.CreateResponseError(err.Error()))
	}
	return ctx.JSON(res)
}

func (ba bookApi) Create(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.CreateBookRequest
	if err := ctx.BodyParser(&req) 
	err != nil{
		return ctx.SendStatus(http.StatusUnprocessableEntity)

	}
	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(http.StatusBadRequest).
		JSON(dto.CreateResponseErrorData("Validation failed", fails))
	}
	err := ba.bookService.Create(c, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
		JSON(dto.CreateResponseError(err.Error()))
		
	}
	return ctx.SendStatus(http.StatusOK)
}

func (ba bookApi) Show(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10 * time.Second)
	defer cancel()

	id := ctx.Params("id")
	res, err := ba.bookService.Show(c, id)

	if err != nil{
		return ctx.Status(http.StatusInternalServerError).
		JSON(dto.CreateResponseError(err.Error()))
	}
	return ctx.JSON(res)
}

func (ba bookApi) Update(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10 * time.Second)
	defer cancel()

	var req dto.UpdateBookRequest
	if err := ctx.BodyParser(&req); err != nil{
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}

	id := ctx.Params("id")
	req.Id = id
	err := ba.bookService.Update(c, req)

	if err != nil{
		return ctx.Status(http.StatusInternalServerError).
		JSON(dto.CreateResponseError(err.Error()))
	}
	return ctx.SendStatus(http.StatusOK)

}

func (ba bookApi) Delete(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10 * time.Second)
	defer cancel()

	id := ctx.Params("id")
	err := ba.bookService.Delete(c, id)

	if err != nil{
		return ctx.Status(http.StatusInternalServerError).
		JSON(dto.CreateResponseError(err.Error()))
	}
	return ctx.SendStatus(http.StatusNoContent)

}