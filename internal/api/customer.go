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

type customerApi struct {
	customerService domain.CustomerService
}

func NewCustomer(app *fiber.App, 
	customerService domain.CustomerService,
	authzMidd fiber.Handler ) {
	ca := customerApi{
		customerService: customerService,
	}

	app.Get("/customers/", authzMidd, ca.index)
	app.Get("/customers/:id", authzMidd, ca.Show)
	app.Post("/customers", authzMidd, ca.Create)
	app.Put("/customers/:id", authzMidd, ca.Update)
	app.Delete("/customers/:id", authzMidd, ca.Delete)
}	

func (ca customerApi) index(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	res, err := ca.customerService.Index(c)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(res)
}


func (ca customerApi) Create(ctx *fiber.Ctx) error{
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.CreateCustomerRequest
	if err := ctx.BodyParser(&req); err != nil{
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}
	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(http.StatusBadRequest).
		JSON(dto.CreateResponseErrorData("validation failed" , fails))
	}
	err := ca.customerService.Create(c, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
		JSON(dto.CreateResponseError(err.Error()))
	}
	return ctx.SendStatus(http.StatusCreated)
}

func (ca customerApi) Update(ctx *fiber.Ctx) error{
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.UpdateCustomerRequest
	if err := ctx.BodyParser(&req); err != nil{
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}
	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(http.StatusBadRequest).
		JSON(dto.CreateResponseErrorData("validation error", fails))
	}
	req.ID = ctx.Params("id")
	err := ca.customerService.Update(c, req)
	if err != nil{
		return ctx.Status(http.StatusInternalServerError).
		JSON(dto.CreateResponseError(err.Error()))
	}
	return ctx.SendStatus(http.StatusOK)
}

func(ca customerApi) Delete(ctx *fiber.Ctx) error{
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id := ctx.Params("id")
	err := ca.customerService.Delete(c, id)
	if err != nil{
		return ctx.Status(http.StatusInternalServerError).
		JSON(dto.CreateResponseError(err.Error()))
	}
	return ctx.SendStatus(http.StatusNoContent)
}

func(ca customerApi) Show(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id := ctx.Params("id")
	data, err:= ca.customerService.Show(c, id)
	if err != nil{
		return ctx.Status(http.StatusInternalServerError).
		JSON(dto.CreateResponseError(err.Error()))
	}
	return ctx.Status(http.StatusOK).
	JSON(data)
}