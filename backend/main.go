package main

import (
	"log"

	"chess-backend/config"
	"chess-backend/models"
	"chess-backend/routes"
	"chess-backend/websocket"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize database
	config.InitDatabase()

	// Auto-migrate database schema
	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Lobby{},
		&models.Game{},
		&models.Move{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return ctx.Status(code).JSON(fiber.Map{
				"error": err.Error(),
				"code":  code,
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: false,
	}))

	// Start WebSocket hub
	go websocket.GameHub.Run()

	// Setup routes
	routes.SetupRoutes(app)

	// Start server
	port := ":" + config.AppConfig.Port
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(port))
}