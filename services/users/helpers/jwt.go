package helpers

import (
	"time"

	"github.com/golang-jwt/jwt"
)

// Token creation methods
func CreateToken(email string) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["email"] = email
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	// The secret should be in an env file
	token, err := at.SignedString([]byte("ACCESS_SECRET"))
	if err != nil {
		return "", err
	}
	return token, nil
}
