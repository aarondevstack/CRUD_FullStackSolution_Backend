package dto

import "time"

// CommentResponse represents comment data in API responses
type CommentResponse struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	BlogID    int64     `json:"blog_id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateCommentRequest represents comment creation request
type CreateCommentRequest struct {
	Content string `json:"content" validate:"required"`
	BlogID  int64  `json:"blog_id" validate:"required"`
	UserID  int64  `json:"user_id" validate:"required"`
}

// UpdateCommentRequest represents comment update request
type UpdateCommentRequest struct {
	Content *string `json:"content,omitempty" validate:"omitempty,min=1"`
}
