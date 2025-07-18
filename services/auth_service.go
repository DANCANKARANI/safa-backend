package services


import (

	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type Claims struct {
	UserID *uuid.UUID `json:"user_id"`
	FullName	string `json:"full_name"`
	Role string `json:"role"`
	jwt.StandardClaims
}
//the function loads .env and return the secretKey
func LoadSecretKey()string{
	my_secret_key := os.Getenv("MY_SECRET_KEY")
	if my_secret_key==""{
		err := godotenv.Load(".env")
	if err != nil {
		return err.Error()
	}
	}
	
	return my_secret_key
}

/*
Generates a JWT token
@params claims *Claims 
@params expiration time
*/
func GenerateToken(claims Claims,expiration_time time.Duration) (string, error) {
	my_secret_key := LoadSecretKey()
	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: time.Now().Add(expiration_time).Unix(),
		Issuer:    "dancan",
		
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(my_secret_key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

/*
Validates the token string
@params tokenString
*/ 
func ValidateToken(tokenString string)(*Claims,error){
	token,err := jwt.ParseWithClaims(tokenString, &Claims{},func(token *jwt.Token) (interface{}, error) {
		return []byte("MY_SECRET_KEY"), nil
	})
	if err != nil {
		return nil, err
	}
	claims,ok := token.Claims.(*Claims); 
	if ! ok{
		return nil, errors.New("invalid user token")
	}
	

	return claims,nil
}

/*
Invalidates token when logged out
@params tokenString
*/

/*
gets the users id from the token
@params claims *Claims
*/
func GetAuthUserID(c *fiber.Ctx,claims *Claims)(*uuid.UUID,error){
	if claims == nil{
		return nil, errors.New("unauthorized user denied. user details not found")
	}
	//extract the user ID from the claims
	userClaimID := claims.UserID
	if userClaimID ==nil{
		return nil, errors.New("unauthorized user denied. user details not found")
	}
	return userClaimID, nil
}