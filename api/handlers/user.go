package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/NikhilSharma03/noteman-go-grpc/api/models"
	"github.com/NikhilSharma03/noteman-go-grpc/services/users/userspb"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	l  *log.Logger
	uc userspb.UsersClient
}

func NewUserHandler(l *log.Logger, uc userspb.UsersClient) *UserHandler {
	return &UserHandler{l, uc}
}

func (u *UserHandler) Login(rw http.ResponseWriter, rq *http.Request) {
	u.l.Println("Login called...")

	user := models.NewUser()
	err := user.FromJson(rq.Body)
	if err != nil {
		http.Error(rw, "Something went wrong while unmarshal", http.StatusInternalServerError)
		return
	}

	res, err := u.uc.Login(context.Background(), &userspb.UserRequest{Email: user.Email, Password: user.Password})
	e, _ := status.FromError(err)

	if err != nil {
		http.Error(rw, e.Message(), http.StatusInternalServerError)
		return
	}

	ures := models.NewUserResponse(res.GetToken(), res.GetEmail())
	err = ures.ToJson(rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (u *UserHandler) SignUp(rw http.ResponseWriter, rq *http.Request) {
	u.l.Println("Signup called...")

	user := models.NewUser()
	err := user.FromJson(rq.Body)
	if err != nil {
		http.Error(rw, "Something went wrong while unmarshal", http.StatusInternalServerError)
		return
	}

	res, err := u.uc.SignUp(context.Background(), &userspb.UserRequest{Email: user.Email, Password: user.Password})
	e, _ := status.FromError(err)

	if err != nil {
		http.Error(rw, e.Message(), http.StatusInternalServerError)
		return
	}

	ures := models.NewUserResponse(res.GetToken(), res.GetEmail())
	err = ures.ToJson(rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
