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
