package repositories

import (
	"blog-api/Domain"
	"blog-api/Usecases"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

func (ur *UserRepo) CheckUserExists(UsernameOrEmail,password string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user domain.User

	filter := bson.M{

		"$or" : []bson.M{
			{"username": UsernameOrEmail},
			{"email": UsernameOrEmail},
		},
	}

	err := ur.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil{
		return nil, errors.New("invalid credentials")
	}

	return &user, nil

}

func (ur *UserRepo) FindByEmail(email string) (*domain.User, error) {
	var user domain.User

	err := ur.Collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}