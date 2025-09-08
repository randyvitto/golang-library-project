package api

import (
	"belajar-golang-rest-api/lat/domain"
	"belajar-golang-rest-api/lat/dto"
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthApi struct {
	authservice domain.AuthService
}

func NewAuth( app *fiber.App, authService domain.AuthService){
	aa := AuthApi{
		authservice : authService,
	}
	app.Post("/auth", aa.Login)
}

func (aa AuthApi) Login(ctx *fiber.Ctx) error{
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.AuthRequest
	if err := ctx.BodyParser(&req); err != nil{
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}

	res, err := aa.authservice.Login(c, req)
	if err != nil{
		return ctx.Status(http.StatusInternalServerError).
		JSON(dto.CreateResponseError(err.Error()))
	}
	return ctx.Status(http.StatusOK).
	JSON(dto.CreateResponseSuccess(res))
}