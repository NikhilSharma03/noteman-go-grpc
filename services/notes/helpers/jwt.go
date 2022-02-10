package helpers

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

type JWTMetaData struct {
	Email string
}

func VerifyToken(rqToken string) (*jwt.Token, error) {
	tokenString := rqToken
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("ACCESS_SECRET"), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// Token extract methods
func ExtractTokenMetadata(rqToken string) (*JWTMetaData, error) {
	token, err := VerifyToken(rqToken)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userEmail := fmt.Sprintf("%v", claims["email"])
		return &JWTMetaData{userEmail}, nil
	}
	return nil, err
}
