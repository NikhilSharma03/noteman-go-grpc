package models

import (
	"encoding/json"
	"io"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func NewUser() *User {
	return &User{}
}

func NewUserResponse(token, email string) *UserResponse {
	return &UserResponse{Token: token, Email: email}
}

func (u *UserResponse) ToJson(r io.Writer) error {
	je := json.NewEncoder(r)
	return je.Encode(u)
}

func (u *User) FromJson(r io.Reader) error {
	je := json.NewDecoder(r)
	return je.Decode(u)
}
