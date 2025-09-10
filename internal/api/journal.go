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

type journalApi struct {
	journalService domain.JournalService
}

func NewJournal(app *fiber.App, 
	journalService domain.JournalService, 
	authzMidd fiber.Handler){
	
	ja := journalApi{
		journalService: journalService,
	}

	app.Get("/journals",authzMidd, ja.Index)
	app.Post("/journals",authzMidd, ja.Create)
	app.Put("journals/:id",authzMidd, ja.Update)
}

func (ja journalApi) Index(ctx *fiber.Ctx) error{
	c, cancel := context.WithTimeout(ctx.Context(), 10 * time.Second)
	defer cancel()

	customerId := ctx.Query("customer_id")
	status := ctx.Query("status")
	res , err := ja.journalService.Index(c, domain.JournalSearch{
		CustomerId: customerId,
		Status: status,
	})
	if err != nil{
		return ctx.Status(http.StatusInternalServerError).
		JSON(dto.CreateResponseError(err.Error()))
	}
	return ctx.Status(http.StatusOK).
	JSON(res)
}

func (ja journalApi) Create(ctx *fiber.Ctx) error{
	c, cancel := context.WithTimeout(ctx.Context(), 10 * time.Second)
	defer cancel()

	var req dto.CreateJournalRequest
	if err :=  ctx.BodyParser(&req); err != nil{
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(http.StatusBadRequest).
		JSON(dto.CreateResponseErrorData("validate error", fails))
	}


 	err := ja.journalService.Create(c, req)
	if err != nil{
		return ctx.Status(http.StatusInternalServerError).
		JSON(dto.CreateResponseError(err.Error()))
	}
	return ctx.SendStatus(http.StatusCreated)
}

func (ja journalApi) Update(ctx *fiber.Ctx) error{
	c, cancel := context.WithTimeout(ctx.Context(), 10 * time.Second)
	defer cancel()

	id := ctx.Params("id") 

 	err := ja.journalService.Return(c, dto.ReturnJournalRequest{
		JournalId: id,
	})
	if err != nil{
		return ctx.Status(http.StatusInternalServerError).
		JSON(dto.CreateResponseError(err.Error()))
	}
	return ctx.SendStatus(http.StatusCreated)
}
