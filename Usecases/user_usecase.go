package usecases
import (
	"blog-api/Domain"
)
type UserUsecase struct{
	Repo UserRepoI
}
func NewUserUsecase (repo UserRepoI) *UserUsecase{
    return &UserUsecase{
		Repo : repo,
	}
}
func (uc * UserUsecase) Register (user *domain.User) error{
	err := uc.Repo.Register(user)
	return err
}
