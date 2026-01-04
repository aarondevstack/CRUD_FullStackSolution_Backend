package handlers

import (
	"context"
	"strconv"

	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/database/ent"
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/database/ent/blog"
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/database/ent/comment"
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/database/ent/user"
	"github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/services/api/dto"
	"github.com/gofiber/fiber/v2"
)

type CommentHandler struct {
	client *ent.Client
}

func NewCommentHandler(client *ent.Client) *CommentHandler {
	return &CommentHandler{client: client}
}

// GetComments returns list of comments
// @Security Bearer
// @Summary Get all comments
// @Tags comments
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {array} dto.CommentResponse
// @Router /comments [get]
func (h *CommentHandler) GetComments(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	comments, err := h.client.Comment.Query().
		WithAuthor().
		WithBlog().
		Limit(limit).
		Offset(offset).
		All(context.Background())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var response []dto.CommentResponse
	for _, cm := range comments {
		resp := dto.CommentResponse{
			ID:        int64(cm.ID),
			Content:   cm.Content,
			BlogID:    int64(cm.Edges.Blog.ID),
			UserID:    int64(cm.Edges.Author.ID),
			Username:  cm.Edges.Author.Username,
			CreatedAt: cm.CreatedAt,
			UpdatedAt: cm.UpdatedAt,
		}
		response = append(response, resp)
	}

	return c.JSON(response)
}

// CreateComment creates a new comment
// @Security Bearer
// @Summary Create comment
// @Tags comments
// @Accept json
// @Produce json
// @Param request body dto.CreateCommentRequest true "Comment details"
// @Success 201 {object} dto.CommentResponse
// @Router /comments [post]
func (h *CommentHandler) CreateComment(c *fiber.Ctx) error {
	var req dto.CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Verify user and blog exist
	userExists, _ := h.client.User.Query().Where(user.ID(int(req.UserID))).Exist(context.Background())
	blogExists, _ := h.client.Blog.Query().Where(blog.ID(int(req.BlogID))).Exist(context.Background())

	if !userExists || !blogExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID or blog ID"})
	}

	cm, err := h.client.Comment.Create().
		SetContent(req.Content).
		SetAuthorID(int(req.UserID)).
		SetBlogID(int(req.BlogID)).
		Save(context.Background())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Re-query to get edges
	cm, _ = h.client.Comment.Query().
		Where(comment.ID(cm.ID)).
		WithAuthor().
		WithBlog().
		Only(context.Background())

	return c.Status(fiber.StatusCreated).JSON(dto.CommentResponse{
		ID:        int64(cm.ID),
		Content:   cm.Content,
		BlogID:    int64(cm.Edges.Blog.ID),
		UserID:    int64(cm.Edges.Author.ID),
		Username:  cm.Edges.Author.Username,
		CreatedAt: cm.CreatedAt,
		UpdatedAt: cm.UpdatedAt,
	})
}
