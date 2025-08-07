package domain

type TokenResponse struct {
	AccessToken  string
	RefreshToken string
}

type UpdateBlogInput struct {
	BlogID  string
	UserID  string
	Title   string
	Content string
	Tags    []string
}
type RequestPasswordResetInput struct {
	Email string
}

type ResetPasswordInput struct {
	Token       string
	NewPassword string
}

type PaginatedBlogsResponse struct {
    Blogs      []Blog `json:"blogs"`
    Page       int    `json:"page"`
    Limit      int    `json:"limit"`
    Total      int64  `json:"total"`
    TotalPages int    `json:"total_pages"`
    HasNext    bool   `json:"has_next"`
    HasPrev    bool   `json:"has_prev"`
}
