package controllers

import (
	"chess-backend/config"
	"chess-backend/middleware"
	"chess-backend/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User        models.User `json:"user"`
	AccessToken string      `json:"accessToken"`
}

func (ac *AuthController) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  "E_BODY_400",
		})
	}

	// Check if username already exists
	var existingUser models.User
	if err := config.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{
			"error": "Username already exists",
			"code":  "E_USER_409",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to hash password",
			"code":  "E_HASH_500",
		})
	}

	// Create user
	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Rating:   1200,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create user",
			"code":  "E_DB_500",
		})
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate token",
			"code":  "E_TOKEN_500",
		})
	}

	return c.Status(201).JSON(AuthResponse{
		User:        user,
		AccessToken: token,
	})
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  "E_BODY_400",
		})
	}

	// Find user
	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid credentials",
			"code":  "E_CREDS_401",
		})
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid credentials",
			"code":  "E_CREDS_401",
		})
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(user.ID, user.Username)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate token",
			"code":  "E_TOKEN_500",
		})
	}

	return c.JSON(AuthResponse{
		User:        user,
		AccessToken: token,
	})
}

func (ac *AuthController) GetMe(c *fiber.Ctx) error {
	userID, _ := middleware.GetUserFromContext(c)

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
			"code":  "E_USER_404",
		})
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}

func (ac *AuthController) UpdateMe(c *fiber.Ctx) error {
	userID, _ := middleware.GetUserFromContext(c)

	var req struct {
		AvatarURL *string `json:"avatarUrl"`
		FCMToken  *string `json:"fcmToken"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
			"code":  "E_BODY_400",
		})
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
			"code":  "E_USER_404",
		})
	}

	// Update fields
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}
	if req.FCMToken != nil {
		user.FCMToken = req.FCMToken
	}

	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update user",
			"code":  "E_DB_500",
		})
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}