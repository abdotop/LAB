package controllers

import (
	"chess-backend/chess"
	"chess-backend/config"
	"chess-backend/middleware"
	"chess-backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GameController struct{}

type CreateGameRequest struct {
	WhiteID     uuid.UUID            `json:"whiteId"`
	BlackID     uuid.UUID            `json:"blackId"`
	TimeControl models.TimeControl   `json:"timeControl"`
}

type GameSnapshot struct {
	ID          uuid.UUID      `json:"id"`
	FEN         string         `json:"fen"`
	PGN         string         `json:"pgn"`
	Turn        string         `json:"turn"`
	Clocks      models.GameClocks `json:"clocks"`
	Players     GamePlayers    `json:"players"`
	LegalMoves  []string       `json:"legalMoves,omitempty"`
	Result      *models.GameResult `json:"result,omitempty"`
	Termination *models.Termination `json:"termination,omitempty"`
}

type GamePlayers struct {
	White models.User `json:"white"`
	Black models.User `json:"black"`
}

func (gc *GameController) CreateGame(c *fiber.Ctx) error {
	var req CreateGameRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  "E_BODY_400",
		})
	}

	game := models.Game{
		WhiteID: req.WhiteID,
		BlackID: req.BlackID,
		FEN:     chess.StartingFEN,
		PGN:     "",
		Turn:    "w",
		Clocks: models.GameClocks{
			WhiteMs: req.TimeControl.Initial * 1000,
			BlackMs: req.TimeControl.Initial * 1000,
		},
	}

	if err := config.DB.Create(&game).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create game",
			"code":  "E_DB_500",
		})
	}

	// Load game with players
	config.DB.Preload("White").Preload("Black").First(&game, game.ID)

	return c.Status(201).JSON(fiber.Map{
		"game": game,
	})
}

func (gc *GameController) GetGame(c *fiber.Ctx) error {
	gameIDParam := c.Params("id")
	gameID, err := uuid.Parse(gameIDParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid game ID",
			"code":  "E_ID_400",
		})
	}

	var game models.Game
	if err := config.DB.Preload("White").Preload("Black").First(&game, gameID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Game not found",
			"code":  "E_GAME_404",
		})
	}

	// Generate legal moves for current position
	pos, err := chess.ParseFEN(game.FEN)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Invalid game position",
			"code":  "E_POS_500",
		})
	}

	legalMoves := pos.GenerateLegalMoves()
	legalMoveStrings := make([]string, len(legalMoves))
	for i, move := range legalMoves {
		legalMoveStrings[i] = move.From + move.To + move.Promotion
	}

	snapshot := GameSnapshot{
		ID:     game.ID,
		FEN:    game.FEN,
		PGN:    game.PGN,
		Turn:   game.Turn,
		Clocks: game.Clocks,
		Players: GamePlayers{
			White: game.White,
			Black: game.Black,
		},
		LegalMoves:  legalMoveStrings,
		Result:      game.Result,
		Termination: game.Termination,
	}

	return c.JSON(snapshot)
}

func (gc *GameController) GetGameMoves(c *fiber.Ctx) error {
	gameIDParam := c.Params("id")
	gameID, err := uuid.Parse(gameIDParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid game ID",
			"code":  "E_ID_400",
		})
	}

	var moves []models.Move
	if err := config.DB.
		Preload("Mover").
		Where("game_id = ?", gameID).
		Order("created_at ASC").
		Find(&moves).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch moves",
			"code":  "E_DB_500",
		})
	}

	return c.JSON(moves)
}

func (gc *GameController) OfferDraw(c *fiber.Ctx) error {
	userID, _ := middleware.GetUserFromContext(c)
	
	gameIDParam := c.Params("id")
	gameID, err := uuid.Parse(gameIDParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid game ID",
			"code":  "E_ID_400",
		})
	}

	var game models.Game
	if err := config.DB.First(&game, gameID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Game not found",
			"code":  "E_GAME_404",
		})
	}

	// Check if user is a player in the game
	if game.WhiteID != userID && game.BlackID != userID {
		return c.Status(403).JSON(fiber.Map{
			"error": "Not a player in this game",
			"code":  "E_GAME_403",
		})
	}

	// Check if game is still ongoing
	if game.Result != nil {
		return c.Status(409).JSON(fiber.Map{
			"error": "Game already finished",
			"code":  "E_GAME_409",
		})
	}

	// TODO: Implement draw offer logic (store in database or memory)
	// For now, just return success
	return c.JSON(fiber.Map{
		"status": "OFFERED",
		"message": "Draw offer sent",
	})
}

func (gc *GameController) RespondToDraw(c *fiber.Ctx) error {
	userID, _ := middleware.GetUserFromContext(c)
	
	gameIDParam := c.Params("id")
	gameID, err := uuid.Parse(gameIDParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid game ID",
			"code":  "E_ID_400",
		})
	}

	var req struct {
		Accept bool `json:"accept"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  "E_BODY_400",
		})
	}

	var game models.Game
	if err := config.DB.First(&game, gameID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Game not found",
			"code":  "E_GAME_404",
		})
	}

	// Check if user is a player in the game
	if game.WhiteID != userID && game.BlackID != userID {
		return c.Status(403).JSON(fiber.Map{
			"error": "Not a player in this game",
			"code":  "E_GAME_403",
		})
	}

	if req.Accept {
		// Accept draw
		result := models.GameResultDraw
		termination := models.TerminationAgreement
		game.Result = &result
		game.Termination = &termination

		if err := config.DB.Save(&game).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to update game",
				"code":  "E_DB_500",
			})
		}

		return c.JSON(fiber.Map{
			"result":      result,
			"termination": termination,
			"message":     "Draw accepted",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Draw declined",
	})
}

func (gc *GameController) Resign(c *fiber.Ctx) error {
	userID, _ := middleware.GetUserFromContext(c)
	
	gameIDParam := c.Params("id")
	gameID, err := uuid.Parse(gameIDParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid game ID",
			"code":  "E_ID_400",
		})
	}

	var game models.Game
	if err := config.DB.First(&game, gameID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Game not found",
			"code":  "E_GAME_404",
		})
	}

	// Check if user is a player in the game
	if game.WhiteID != userID && game.BlackID != userID {
		return c.Status(403).JSON(fiber.Map{
			"error": "Not a player in this game",
			"code":  "E_GAME_403",
		})
	}

	// Check if game is still ongoing
	if game.Result != nil {
		return c.Status(409).JSON(fiber.Map{
			"error": "Game already finished",
			"code":  "E_GAME_409",
		})
	}

	// Determine winner
	var result models.GameResult
	if game.WhiteID == userID {
		result = models.GameResultBlackWins
	} else {
		result = models.GameResultWhiteWins
	}

	termination := models.TerminationResign
	game.Result = &result
	game.Termination = &termination

	if err := config.DB.Save(&game).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update game",
			"code":  "E_DB_500",
		})
	}

	return c.JSON(fiber.Map{
		"result":      result,
		"termination": termination,
		"message":     "Game resigned",
	})
}