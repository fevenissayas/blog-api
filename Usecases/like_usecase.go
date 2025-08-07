package usecases

import (
	domain "blog-api/Domain"
	"time"
)

type LikeUsecase struct {
	likeRepo domain.IBlogLikeRepository
}

func NewLikeUsecase(likeRepo domain.IBlogLikeRepository) domain.IBlogLikeUsecase {
	return &LikeUsecase{
		likeRepo: likeRepo,
	}
}

func (uc *LikeUsecase) LikeBlog(blogID, userID string) error {
	like := &domain.Like{
		BlogID:    blogID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	return uc.likeRepo.AddLike(like)
}

func (uc *LikeUsecase) UnlikeBlog(blogID, userID string) error {
	return uc.likeRepo.RemoveLike(blogID, userID)
}

func (uc *LikeUsecase) IsBlogLiked(blogID, userID string) (bool, error) {
	return uc.likeRepo.IsLiked(blogID, userID)
}

func (uc *LikeUsecase) GetLikeCount(blogID string) (int, error) {
	return uc.likeRepo.CountLikes(blogID)
}
