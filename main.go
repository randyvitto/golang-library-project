package main

import (
	"belajar-golang-rest-api/lat/dto"
	"belajar-golang-rest-api/lat/internal/api"
	"belajar-golang-rest-api/lat/internal/config"
	"belajar-golang-rest-api/lat/internal/connection"
	"belajar-golang-rest-api/lat/internal/repository"
	"belajar-golang-rest-api/lat/internal/service"
	"net/http"

	jwtMid "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cnf := config.Get()
	dbConnection := connection.GetDatabase(cnf.Database)

	app := fiber.New()

	jwtMid := jwtMid.New(jwtMid.Config{
		SigningKey: jwtMid.SigningKey{Key: []byte(cnf.Jwt.Key)},
		ErrorHandler: func(ctx *fiber.Ctx, err error)error {
			return ctx.Status(http.StatusUnauthorized).
			JSON(dto.CreateResponseError("Unauthorized, Please Login"))
			
		},
	})

	customerRepository := repository.NewCustomer(dbConnection)
	userRepository := repository.NewUser(dbConnection)
	bookRepository := repository.NewBook(dbConnection)
	bookStockRepository := repository.NewBookStock(dbConnection)
	journalRepository := repository.NewJournal(dbConnection)
	mediaRepository := repository.NewMedia(dbConnection)

	customerService := service.NewCustomer(customerRepository)
	authService := service.NewAuth(cnf, 
		userRepository)
	bookService := service.NewBook(cnf, bookRepository, 
		bookStockRepository, mediaRepository)
	bookStockService := service.NewBookStock(bookRepository, 
		bookStockRepository)
	journalservice := service.NewJournal(journalRepository,
	bookRepository,
	bookStockRepository,
	customerRepository)
	mediaService := service.NewMedia(cnf,
	mediaRepository)

	api.NewCustomer(app, customerService, jwtMid)
	api.NewAuth(app, authService)
	api.NewBook(app, bookService, jwtMid)
	api.NewBookStock(app, bookStockService, jwtMid)
	api.NewJournal(app, journalservice, jwtMid)
	api.NewMedia(app, cnf, mediaService, jwtMid )

	app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}

func developers(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON("data")
}
