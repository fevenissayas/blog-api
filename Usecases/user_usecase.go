package usecases

import (
	domain "blog-api/Domain"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	Repo UserRepoI
}

func NewUserUsecase(repo UserRepoI) *UserUsecase {
	return &UserUsecase{
		Repo: repo,
	}
}
func (uc *UserUsecase) Register(user *domain.User) error {
	err := uc.Repo.Register(user)
	return err
}

func (uc *UserUsecase) AuthenticateUser(usernameOrEmail, password string) (*domain.User, error) {
	user, err := uc.Repo.CheckUserExists(usernameOrEmail, password)

	if err != nil {
		return nil, errors.New("database error:")
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	//TODO: token generation logic
	//return token , nil
	return user, nil
}
