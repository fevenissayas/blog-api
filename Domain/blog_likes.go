package domain

import (
	"time"
)

type Like struct {
	ID        string
	BlogID    string
	UserID    string
	CreatedAt time.Time
}

type IBlogLikeRepository interface {
	AddLike(like *Like) error
	RemoveLike(blogID, userID string) error
	IsLiked(blogID, userID string) (bool, error)
	CountLikes(blogID string) (int, error)
}

type IBlogLikeUsecase interface {
	LikeBlog(blogID, userID string) error
	RemoveLikeBlog(blogID, userID string) error
	IsBlogLiked(blogID, userID string) (bool, error)
	GetLikeCount(blogID string) (int, error)
}
