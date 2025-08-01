package usecases

import (
	domain "blog-api/Domain"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	Repo UserRepoI
	TokenService TokenI
}

func NewUserUsecase(repo UserRepoI, tokenservice TokenI) *UserUsecase {
	return &UserUsecase{
		Repo: repo,
		TokenService: tokenservice,
	}
}
func (uc *UserUsecase) Register(user *domain.User) error {
	err := uc.Repo.Register(user)
	return err
}
func (uc *UserUsecase) AuthenticateUser(usernameOrEmail, password string) (string, error) {
	user, err := uc.Repo.CheckUserExists(usernameOrEmail, password)

	if err != nil {
		return "", errors.New("database error:")
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
    token , err := uc.TokenService.GenerateToken(user)    
	if err != nil{
		return "", errors.New("could not generate the token")
	}
	return token,nil
}
