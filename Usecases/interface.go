package usecases
import (
	"blog-api/Domain"
)
type UserRepoI interface{
    Register(* domain.User) error
	FindByEmail(email string) (*domain.User, error)
	CheckUserExists(UsernameOrEmail, password string) (*domain.User, error)
}

