package api

import (
	"belajar-golang-rest-api/lat/domain"
	"belajar-golang-rest-api/lat/dto"
	"belajar-golang-rest-api/lat/internal/util"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type bookStockApi struct {
	bookStockService domain.BookStockService
}

func NewBookStock(app *fiber.App,
bookStockService domain.BookStockService, 
authzMidd fiber.Handler) {
	bsa := bookStockApi{
		bookStockService: bookStockService,
	}
	bookStock := app.Group("/book-stocks", authzMidd)
	
	bookStock.Post("",  bsa.Create)
	bookStock.Delete("",  bsa.Delete)
}

func (ba bookStockApi) Create(ctx *fiber.Ctx) error{
	c, cancel := context.WithTimeout(ctx.Context(), 10* time.Second)
	defer cancel()

	var req dto.CreateBookStockRequest
	if err := ctx.BodyParser(&req); err != nil{
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}
	fails :=util.Validate(req)
	if len(fails) > 0{
		return ctx.Status(http.StatusBadRequest).
		JSON(dto.CreateResponseErrorData("Failed validation", fails))
	}
	err := ba.bookStockService.Create(c, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
		JSON(dto.CreateResponseError(err.Error()))
	}
	return ctx.SendStatus(http.StatusCreated)
}

func (ba bookStockApi) Delete(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()
	codeStr :=ctx.Query("code")
	if codeStr == ""{
		return ctx.Status(http.StatusBadRequest).
		JSON(dto.CreateResponseError("Parameters required"))
	}
	codes:= strings.Split( codeStr,";")

	err:= ba.bookStockService.Delete(c, dto.DeleteBookStockRequest{Codes: codes})
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
		JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.Status(http.StatusOK).
		JSON(dto.CreateResponseSuccess("Book stock deleted successfully"))
}