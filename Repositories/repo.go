package repositories

import (
	"blog-api/Domain"
	usecases "blog-api/Usecases"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)
type UserRepo struct {
    Collection *mongo.Collection
	Context context.Context
}
func  NewUserRepo(c *mongo.Collection) usecases.UserRepoI {
	ctx := context.Background() 
    return &UserRepo{
		Collection: c, 
		Context : ctx,
	} 
}
func (ur *UserRepo) Register(u *domain.User) error{ 
	_, err := ur.Collection.InsertOne(ur.Context, u)
	return err
}
