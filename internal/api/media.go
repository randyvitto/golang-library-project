package api

import (
	"belajar-golang-rest-api/lat/domain"
	"belajar-golang-rest-api/lat/dto"
	"belajar-golang-rest-api/lat/internal/config"
	"context"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type mediaApi struct {
	cnf *config.Config
	mediaService domain.MediaService
}

func NewMedia(app *fiber.App,
	cnf *config.Config,
	mediaService domain.MediaService,
	authzMidd fiber.Handler) {
	
	ma:= mediaApi{
		cnf: cnf,
		mediaService: mediaService,
	}

	app.Post("/media",authzMidd, ma.Create)
	app.Static("/media", cnf.Storage.BasePath)
}

func (ma mediaApi) Create(ctx *fiber.Ctx) error{
	c, cancel := context.WithTimeout(ctx.Context(), 10 *time.Second)
	defer cancel()
	file, err := ctx.FormFile("media")
	if err != nil{
		return ctx.SendStatus(http.StatusBadRequest)
	}
	filename := uuid.NewString() + filepath.Ext(file.Filename)
	path := filepath.Join(ma.cnf.Storage.BasePath, filename)
	err = ctx.SaveFile(file, path)

	if err != nil{
		return ctx.Status(http.StatusInternalServerError).
		JSON(dto.CreateResponseError(err.Error()))
	}
	res, err := ma.mediaService.Create(c, dto.CreateMediaRequest{
		Path: filename,
	})

	if err != nil{
	return ctx.Status(http.StatusInternalServerError).
	JSON(dto.CreateResponseError(err.Error()))
	}
	return ctx.Status(http.StatusCreated).
	JSON(res)
}

	