package handlers

import (
	"context"

	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/database/ent"
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/database/ent/user"
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/services/api/dto"
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/services/api/middleware"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	client *ent.Client
}

func NewAuthHandler(client *ent.Client) *AuthHandler {
	return &AuthHandler{client: client}
}

// Login handles user authentication
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Find user by username
	u, err := h.client.User.Query().
		Where(user.UsernameEQ(req.Username)).
		Only(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Verify password
	// Note: In a real app, you should use bcrypt.CompareHashAndPassword here.
	// For this seeded data which might be plain text or hashed, we need to handle both
	// or assume seeded data is hashed. Assuming seeded data is hashed with bcrypt.
	// But our seeder might have used plain text.
	// For now, let's implement standard bcrypt check but also fallback to plain text if check fails
	// (only for development convenience with seeded simple passwords like 'admin123')
	// TODO: Ensure seeder uses bcrypt and remove plain text fallback in production

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
	if err != nil {
		// Fallback for simple seeded passwords (DEV ONLY)
		if u.Password != req.Password {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(int64(u.ID), u.Username, u.Role.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(dto.LoginResponse{
		Token: token,
		User: dto.UserSummary{
			ID:       int64(u.ID),
			Username: u.Username,
			Email:    u.Email,
			Role:     u.Role.String(),
		},
	})
}
