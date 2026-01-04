package dto

// PaginatedUserResponse represents paginated user list response
type PaginatedUserResponse struct {
	Data  []UserResponse `json:"data"`
	Total int            `json:"total"`
}

// PaginatedBlogResponse represents paginated blog list response
type PaginatedBlogResponse struct {
	Data  []BlogResponse `json:"data"`
	Total int            `json:"total"`
}

// PaginatedCommentResponse represents paginated comment list response
type PaginatedCommentResponse struct {
	Data  []CommentResponse `json:"data"`
	Total int               `json:"total"`
}
