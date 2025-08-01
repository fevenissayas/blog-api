package infrastructure
import (
    "github.com/golang-jwt/jwt/v5"
	"blog-api/Usecases"
	"blog-api/Domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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
