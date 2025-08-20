package app

import (
	"fmt"
	"log"

	"trip-planner/internal/config"
	"trip-planner/internal/repository"
	"trip-planner/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Run(cfg *config.Config) {
	repo, err := repository.NewPostgressRepo(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to create repo: %v", err)
	}

	defer repo.Close()

	srv := service.NewService(repo)
	handler := NewHandler(srv)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
            AllowOrigins: "http://localhost:5173", // или "*" если для всего
            AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
            AllowHeaders: "Origin, Content-Type, Accept",
        }))

	handler.RegisterRoutes(app)

	fmt.Printf("Server is running on :3000\n")
	log.Fatal(app.Listen(":3000"))

}
