package handlers

import (
	"context"
	"strconv"

	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/database/ent"
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/database/ent/blog"
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/database/ent/user"
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/services/api/dto"
	"github.com/gofiber/fiber/v2"
)

type BlogHandler struct {
	client *ent.Client
}

func NewBlogHandler(client *ent.Client) *BlogHandler {
	return &BlogHandler{client: client}
}

// GetBlogs returns list of blogs
// @Security Bearer
// @Summary Get all blogs
// @Tags blogs
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {array} dto.BlogResponse
// @Router /blogs [get]
func (h *BlogHandler) GetBlogs(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	blogs, err := h.client.Blog.Query().
		WithAuthor().
		Limit(limit).
		Offset(offset).
		All(context.Background())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var response []dto.BlogResponse
	for _, b := range blogs {
		resp := dto.BlogResponse{
			ID:        int64(b.ID),
			Title:     b.Title,
			Content:   b.Content,
			UserID:    int64(b.Edges.Author.ID),
			Username:  b.Edges.Author.Username,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		}
		response = append(response, resp)
	}

	return c.JSON(response)
}

// CreateBlog creates a new blog
// @Security Bearer
// @Summary Create blog
// @Tags blogs
// @Accept json
// @Produce json
// @Param request body dto.CreateBlogRequest true "Blog details"
// @Success 201 {object} dto.BlogResponse
// @Router /blogs [post]
func (h *BlogHandler) CreateBlog(c *fiber.Ctx) error {
	var req dto.CreateBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Verify user exists first
	exists, err := h.client.User.Query().Where(user.ID(int(req.UserID))).Exist(context.Background())
	if err != nil || !exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	b, err := h.client.Blog.Create().
		SetTitle(req.Title).
		SetContent(req.Content).
		SetAuthorID(int(req.UserID)).
		Save(context.Background())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Re-query to get author info
	b, _ = h.client.Blog.Query().Where(blog.ID(b.ID)).WithAuthor().Only(context.Background())

	return c.Status(fiber.StatusCreated).JSON(dto.BlogResponse{
		ID:        int64(b.ID),
		Title:     b.Title,
		Content:   b.Content,
		UserID:    int64(b.Edges.Author.ID),
		Username:  b.Edges.Author.Username,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	})
}

// GetBlog returns a single blog
// @Security Bearer
// @Summary Get blog by ID
// @Tags blogs
// @Accept json
// @Produce json
// @Param id path int true "Blog ID"
// @Success 200 {object} dto.BlogResponse
// @Router /blogs/{id} [get]
func (h *BlogHandler) GetBlog(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	b, err := h.client.Blog.Query().
		Where(blog.ID(id)).
		WithAuthor().
		Only(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Blog not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(dto.BlogResponse{
		ID:        int64(b.ID),
		Title:     b.Title,
		Content:   b.Content,
		UserID:    int64(b.Edges.Author.ID),
		Username:  b.Edges.Author.Username,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	})
}

// UpdateBlog updates a blog
// @Security Bearer
// @Summary Update blog
// @Tags blogs
// @Accept json
// @Produce json
// @Param id path int true "Blog ID"
// @Param request body dto.UpdateBlogRequest true "Blog details"
// @Success 200 {object} dto.BlogResponse
// @Router /blogs/{id} [patch]
func (h *BlogHandler) UpdateBlog(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var req dto.UpdateBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	update := h.client.Blog.UpdateOneID(id)
	if req.Title != nil {
		update.SetTitle(*req.Title)
	}
	if req.Content != nil {
		update.SetContent(*req.Content)
	}

	b, err := update.Save(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Blog not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Re-query to get author info
	b, _ = h.client.Blog.Query().Where(blog.ID(b.ID)).WithAuthor().Only(context.Background())

	return c.JSON(dto.BlogResponse{
		ID:        int64(b.ID),
		Title:     b.Title,
		Content:   b.Content,
		UserID:    int64(b.Edges.Author.ID),
		Username:  b.Edges.Author.Username,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	})
}

// DeleteBlog deletes a blog
// @Security Bearer
// @Summary Delete blog
// @Description Delete blog and all related comments
// @Tags blogs
// @Param id path int true "Blog ID"
// @Success 204 "No Content"
// @Router /blogs/{id} [delete]
func (h *BlogHandler) DeleteBlog(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	err = h.client.Blog.DeleteOneID(id).Exec(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Blog not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
