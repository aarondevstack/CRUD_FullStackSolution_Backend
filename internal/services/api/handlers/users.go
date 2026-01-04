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

// GetUser returns a single user
// @Security Bearer
// @Summary Get user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.UserResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	u, err := h.client.User.Query().Where(user.ID(id)).Only(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(dto.UserResponse{
		ID:        int64(u.ID),
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role.String(),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	})
}

// UpdateUser updates a user
// @Security Bearer
// @Summary Update user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body dto.UpdateUserRequest true "User details"
// @Success 200 {object} dto.UserResponse
// @Router /users/{id} [patch]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	update := h.client.User.UpdateOneID(id)
	if req.Username != nil {
		update.SetUsername(*req.Username)
	}
	if req.Email != nil {
		update.SetEmail(*req.Email)
	}
	if req.Role != nil {
		update.SetRole(user.Role(*req.Role))
	}
	if req.Password != nil {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		update.SetPassword(string(hashed))
	}

	u, err := update.Save(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(dto.UserResponse{
		ID:        int64(u.ID),
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role.String(),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	})
}

// DeleteUser deletes a user
// @Security Bearer
// @Summary Delete user
// @Description Delete user and all related blogs/comments
// @Tags users
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	err = h.client.User.DeleteOneID(id).Exec(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
