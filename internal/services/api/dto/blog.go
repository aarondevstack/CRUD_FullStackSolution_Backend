package dto

import "time"

// BlogResponse represents blog data in API responses
type BlogResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateBlogRequest represents blog creation request
type CreateBlogRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=255"`
	Content string `json:"content" validate:"required"`
	UserID  int64  `json:"user_id" validate:"required"`
}

// UpdateBlogRequest represents blog update request
type UpdateBlogRequest struct {
	Title   *string `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Content *string `json:"content,omitempty"`
}
