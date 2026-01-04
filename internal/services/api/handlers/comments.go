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
// @Success 200 {object} dto.PaginatedCommentResponse
// @Router /comments [get]
func (h *CommentHandler) GetComments(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	// Get total count
	total, err := h.client.Comment.Query().Count(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	comments, err := h.client.Comment.Query().
		WithAuthor().
		WithBlog().
		Limit(limit).
		Offset(offset).
		All(context.Background())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var data []dto.CommentResponse
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
		data = append(data, resp)
	}

	return c.JSON(dto.PaginatedCommentResponse{
		Data:  data,
		Total: total,
	})
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request: " + err.Error()})
	}

	// Validate required fields
	if req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Content is required"})
	}

	blogID := parseID(req.BlogID)
	if blogID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Blog ID is required"})
	}

	// If UserID is not provided in request, get it from JWT
	userID := parseID(req.UserID)
	if userID == 0 {
		if uid, ok := c.Locals("user_id").(int64); ok {
			userID = uid
		}
	}

	// Verify user and blog exist
	userExists, _ := h.client.User.Query().Where(user.ID(int(userID))).Exist(context.Background())
	blogExists, _ := h.client.Blog.Query().Where(blog.ID(int(blogID))).Exist(context.Background())

	if !userExists || !blogExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID or blog ID"})
	}

	cm, err := h.client.Comment.Create().
		SetContent(req.Content).
		SetAuthorID(int(userID)).
		SetBlogID(int(blogID)).
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

// GetComment returns a single comment
// @Security Bearer
// @Summary Get comment by ID
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} dto.CommentResponse
// @Router /comments/{id} [get]
func (h *CommentHandler) GetComment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	cm, err := h.client.Comment.Query().
		Where(comment.ID(id)).
		WithAuthor().
		WithBlog().
		Only(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Comment not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(dto.CommentResponse{
		ID:        int64(cm.ID),
		Content:   cm.Content,
		BlogID:    int64(cm.Edges.Blog.ID),
		UserID:    int64(cm.Edges.Author.ID),
		Username:  cm.Edges.Author.Username,
		CreatedAt: cm.CreatedAt,
		UpdatedAt: cm.UpdatedAt,
	})
}

// UpdateComment updates a comment
// @Security Bearer
// @Summary Update comment
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Param request body dto.UpdateCommentRequest true "Comment details"
// @Success 200 {object} dto.CommentResponse
// @Router /comments/{id} [patch]
func (h *CommentHandler) UpdateComment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var req dto.UpdateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	update := h.client.Comment.UpdateOneID(id)
	if req.Content != nil {
		update.SetContent(*req.Content)
	}

	cm, err := update.Save(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Comment not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Re-query to get edges
	cm, _ = h.client.Comment.Query().
		Where(comment.ID(cm.ID)).
		WithAuthor().
		WithBlog().
		Only(context.Background())

	return c.JSON(dto.CommentResponse{
		ID:        int64(cm.ID),
		Content:   cm.Content,
		BlogID:    int64(cm.Edges.Blog.ID),
		UserID:    int64(cm.Edges.Author.ID),
		Username:  cm.Edges.Author.Username,
		CreatedAt: cm.CreatedAt,
		UpdatedAt: cm.UpdatedAt,
	})
}

// DeleteComment deletes a comment
// @Security Bearer
// @Summary Delete comment
// @Tags comments
// @Param id path int true "Comment ID"
// @Success 204 "No Content"
// @Router /comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	err = h.client.Comment.DeleteOneID(id).Exec(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Comment not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
