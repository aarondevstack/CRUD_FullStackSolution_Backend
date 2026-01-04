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
