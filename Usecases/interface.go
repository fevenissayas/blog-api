package usecases
import (
	"blog-api/Domain"
)
type UserRepoI interface{
    Register(* domain.User) error
	FindByEmail(email string) (*domain.User, error)
	CheckUserExists(UsernameOrEmail, password string) (*domain.User, error)
}
type TokenI interface{
	GenerateToken(* domain.User) (string,error)
	// VerifyToken(* domain.User) (bool, error)
}
