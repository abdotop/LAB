package controllers

import (
	"chess-backend/config"
	"chess-backend/middleware"
	"chess-backend/models"
	"math/rand"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LobbyController struct{}

type CreateLobbyRequest struct {
	Type        string                 `json:"type" validate:"required,oneof=PUBLIC PRIVATE"`
	TimeControl models.TimeControl     `json:"timeControl"`
	InviteOnly  bool                  `json:"inviteOnly,omitempty"`
}

func (lc *LobbyController) CreateLobby(c *fiber.Ctx) error {
	userID, _ := middleware.GetUserFromContext(c)

	var req CreateLobbyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  "E_BODY_400",
		})
	}

	lobby := models.Lobby{
		Type:        models.LobbyType(req.Type),
		CreatorID:   userID,
		TimeControl: req.TimeControl,
		Status:      models.LobbyStatusOpen,
	}

	// Generate invite code for private lobbies
	if req.Type == "PRIVATE" || req.InviteOnly {
		inviteCode := generateInviteCode()
		lobby.InviteCode = &inviteCode
	}

	if err := config.DB.Create(&lobby).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create lobby",
			"code":  "E_DB_500",
		})
	}

	// Load creator relationship
	config.DB.Preload("Creator").First(&lobby, lobby.ID)

	return c.Status(201).JSON(fiber.Map{
		"lobby": lobby,
	})
}

func (lc *LobbyController) GetPublicLobbies(c *fiber.Ctx) error {
	var lobbies []models.Lobby
	
	if err := config.DB.
		Preload("Creator").
		Where("type = ? AND status = ?", models.LobbyTypePublic, models.LobbyStatusOpen).
		Order("created_at DESC").
		Find(&lobbies).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch lobbies",
			"code":  "E_DB_500",
		})
	}

	return c.JSON(fiber.Map{
		"lobbies": lobbies,
	})
}

func (lc *LobbyController) JoinLobby(c *fiber.Ctx) error {
	userID, _ := middleware.GetUserFromContext(c)
	
	lobbyIDParam := c.Params("id")
	var lobbyID uuid.UUID
	var err error

	// Check if it's a UUID or invite code
	if lobbyID, err = uuid.Parse(lobbyIDParam); err != nil {
		// Try to find by invite code
		var lobby models.Lobby
		if err := config.DB.Where("invite_code = ?", lobbyIDParam).First(&lobby).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "Lobby not found",
				"code":  "E_LOBBY_404",
			})
		}
		lobbyID = lobby.ID
	}

	var lobby models.Lobby
	if err := config.DB.Preload("Creator").First(&lobby, lobbyID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Lobby not found",
			"code":  "E_LOBBY_404",
		})
	}

	// Check if lobby is still open
	if lobby.Status != models.LobbyStatusOpen {
		return c.Status(409).JSON(fiber.Map{
			"error": "Lobby is not available",
			"code":  "E_LOBBY_409",
		})
	}

	// Check if user is trying to join their own lobby
	if lobby.CreatorID == userID {
		return c.Status(409).JSON(fiber.Map{
			"error": "Cannot join your own lobby",
			"code":  "E_LOBBY_409",
		})
	}

	// Get the joiner user
	var joiner models.User
	if err := config.DB.First(&joiner, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
			"code":  "E_USER_404",
		})
	}

	// Create game
	game := models.Game{
		WhiteID: lobby.CreatorID,
		BlackID: userID,
		FEN:     "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		PGN:     "",
		Turn:    "w",
		Clocks: models.GameClocks{
			WhiteMs: lobby.TimeControl.Initial * 1000,
			BlackMs: lobby.TimeControl.Initial * 1000,
		},
	}

	// Randomly assign colors
	if rand.Intn(2) == 1 {
		game.WhiteID = userID
		game.BlackID = lobby.CreatorID
	}

	if err := config.DB.Create(&game).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create game",
			"code":  "E_DB_500",
		})
	}

	// Update lobby status
	lobby.Status = models.LobbyStatusMatched
	config.DB.Save(&lobby)

	// Load game with players
	config.DB.Preload("White").Preload("Black").First(&game, game.ID)

	return c.JSON(fiber.Map{
		"lobby": lobby,
		"game":  game,
	})
}

func (lc *LobbyController) DeleteLobby(c *fiber.Ctx) error {
	userID, _ := middleware.GetUserFromContext(c)
	
	lobbyIDParam := c.Params("id")
	lobbyID, err := uuid.Parse(lobbyIDParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid lobby ID",
			"code":  "E_ID_400",
		})
	}

	var lobby models.Lobby
	if err := config.DB.First(&lobby, lobbyID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Lobby not found",
			"code":  "E_LOBBY_404",
		})
	}

	// Check if user is the creator
	if lobby.CreatorID != userID {
		return c.Status(403).JSON(fiber.Map{
			"error": "Only lobby creator can delete lobby",
			"code":  "E_LOBBY_403",
		})
	}

	// Update status instead of deleting
	lobby.Status = models.LobbyStatusClosed
	if err := config.DB.Save(&lobby).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to close lobby",
			"code":  "E_DB_500",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Lobby closed successfully",
	})
}

func generateInviteCode() string {
	// Generate a 6-character alphanumeric code
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	
	var result strings.Builder
	for i := 0; i < 6; i++ {
		result.WriteByte(charset[rand.Intn(len(charset))])
	}
	return result.String()
}