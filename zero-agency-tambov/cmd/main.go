package main

import (
	"log"

	"zero-agency-tambov/config"
	"zero-agency-tambov/internal/database"
	"zero-agency-tambov/internal/handlers"
	"zero-agency-tambov/internal/repository"
	"zero-agency-tambov/internal/service"
	"zero-agency-tambov/logger"

	"github.com/gofiber/fiber/v2"
)

func main() {
	logger.InitLogger()

	cfg := config.LoadConfig()

	db := database.ConnectDB(cfg)

	newsRepo := repository.NewNewsRepository(db)
	newsService := service.NewNewsService(newsRepo)
	newsHandler := handlers.NewNewsHandler(newsService)

	app := fiber.New()

	app.Post("/news/:id", newsHandler.UpdateNews)
	app.Get("/news", newsHandler.GetAllNews)
	app.Get("/news/:id", newsHandler.GetNewsById)

	log.Fatal(app.Listen(":8080"))
}
