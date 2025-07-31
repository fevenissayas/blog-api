package usecases
import (
	"blog-api/Domain"
)
type UserRepoI interface{
    Register(* domain.User) error
}

