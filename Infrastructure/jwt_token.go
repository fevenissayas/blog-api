package infrastructure

import (
	"blog-api/Domain"
	"blog-api/Usecases"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
type JwtService struct{
}
func NewJwtService() usecases.TokenI{
    return &JwtService{}
}
type Claims struct {
	Username string
	Userid primitive.ObjectID 
	Role domain.Role
	jwt.RegisteredClaims
}
func (js *JwtService) GenerateToken(user *domain.User)(string, error){
	var jwtKey = []byte(Env.Jwt_Secret)
	expirationTime := time.Now().AddDate(0, 0, 1)
	claims := &Claims{
		Username : user.Username,
		Role : user.Role,
		Userid : user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func (js *JwtService) ValidateToken(tokenString string) (*Claims, error) {
	jwtKey := []byte(Env.Jwt_Secret)

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalid
	}

	return claims, nil
}
