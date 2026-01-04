package handlers

import (
	"context"
	"strconv"

	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/database/ent"
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/database/ent/user"
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/services/api/dto"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	client *ent.Client
}

func NewUserHandler(client *ent.Client) *UserHandler {
	return &UserHandler{client: client}
}

// GetUsers returns list of users
// @Security Bearer
// @Summary Get all users
// @Description Get list of users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {array} dto.UserResponse
// @Router /users [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	users, err := h.client.User.Query().
		Limit(limit).
		Offset(offset).
		All(context.Background())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var response []dto.UserResponse
	for _, u := range users {
		response = append(response, dto.UserResponse{
			ID:        int64(u.ID),
			Username:  u.Username,
			Email:     u.Email,
			Role:      u.Role.String(),
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	return c.JSON(response)
}

// CreateUser creates a new user
// @Security Bearer
// @Summary Create user
// @Description Create a new user (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "User details"
// @Success 201 {object} dto.UserResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	u, err := h.client.User.Create().
		SetUsername(req.Username).
		SetPassword(string(hashedPassword)).
		SetEmail(req.Email).
		SetRole(user.Role(req.Role)).
		Save(context.Background())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.UserResponse{
		ID:        int64(u.ID),
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role.String(),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	})
}
