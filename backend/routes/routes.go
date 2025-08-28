package routes

import (
	"chess-backend/controllers"
	"chess-backend/middleware"
	"chess-backend/websocket"

	"github.com/gofiber/fiber/v2"
	websocketHandler "github.com/gofiber/websocket/v2"
)

func SetupRoutes(app *fiber.App) {
	// Initialize controllers
	authController := &controllers.AuthController{}
	lobbyController := &controllers.LobbyController{}
	gameController := &controllers.GameController{}

	// API routes
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)

	// Protected routes
	protected := api.Group("/", middleware.AuthMiddleware)

	// User routes
	users := protected.Group("/users")
	users.Get("/me", authController.GetMe)
	users.Patch("/me", authController.UpdateMe)

	// Lobby routes
	lobbies := protected.Group("/lobbies")
	lobbies.Post("/", lobbyController.CreateLobby)
	lobbies.Get("/public", lobbyController.GetPublicLobbies)
	lobbies.Post("/:id/join", lobbyController.JoinLobby)
	lobbies.Delete("/:id", lobbyController.DeleteLobby)

	// Game routes
	games := protected.Group("/games")
	games.Post("/", gameController.CreateGame)
	games.Get("/:id", gameController.GetGame)
	games.Get("/:id/moves", gameController.GetGameMoves)
	games.Post("/:id/draw-offer", gameController.OfferDraw)
	games.Post("/:id/draw-respond", gameController.RespondToDraw)
	games.Post("/:id/resign", gameController.Resign)

	// WebSocket upgrade middleware
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocketHandler.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket route
	app.Get("/ws", websocketHandler.New(websocket.HandleWebSocket))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"service": "chess-backend",
		})
	})
}