package model

import (
	"context"
	"fmt"

	"github.com/NikhilSharma03/noteman-go-grpc/services/users/db"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	Email    string `json:"string"`
	Password string `json:"string"`
}

func NewUser(email, pass string) *User {
	return &User{Email: email, Password: pass}
}

func (u *User) SignUpUser() error {
	userColl := db.GetCollection("users")

	//Handle user already exists
	count, err := userColl.CountDocuments(context.Background(), bson.M{"email": u.Email})
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("User already exists")
	}

	// Store user in DB
	_, err = userColl.InsertOne(context.Background(), u)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) LoginUser() (bool, error) {
	userColl := db.GetCollection("users")

	isUser := &User{}

	// Check for user in DB
	err := userColl.FindOne(context.Background(), bson.M{"email": u.Email}).Decode(isUser)
	if err != nil {
		return false, fmt.Errorf("no user exits with provided email")
	}

	if isUser.Password == u.Password {
		return true, nil
	}

	return false, fmt.Errorf("invalid password")
}
