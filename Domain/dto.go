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