package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/NikhilSharma03/noteman-go-grpc/services/users/db"
	"github.com/NikhilSharma03/noteman-go-grpc/services/users/helpers"
	"github.com/NikhilSharma03/noteman-go-grpc/services/users/model"
	"github.com/NikhilSharma03/noteman-go-grpc/services/users/userspb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// Server struct
type server struct {
	userspb.UnimplementedUsersServer
}

// Login
func (*server) Login(ctx context.Context, rq *userspb.UserRequest) (*userspb.UserResponse, error) {
	email, password := rq.GetEmail(), rq.GetPassword()
	if email == "" || password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Please provide an email and password!")
	}

	// User login
	user := model.NewUser(email, password)
	isAuthenticated, err := user.LoginUser()
	if err != nil || !isAuthenticated {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// Token generation
	token, err := helpers.CreateToken(email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &userspb.UserResponse{
		Token: token,
		Email: email,
	}, nil
}

// SignUp
func (*server) SignUp(ctx context.Context, rq *userspb.UserRequest) (*userspb.UserResponse, error) {
	email, password := rq.GetEmail(), rq.GetPassword()
	if email == "" || password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Please provide an email and password!")
	}

	// User registration
	user := model.NewUser(email, password)
	err := user.SignUpUser()
	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, err.Error())
	}

	// Token creation
	token, err := helpers.CreateToken(email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &userspb.UserResponse{
		Token: token,
		Email: user.Email,
	}, nil
}

func main() {
	// Listener
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalln("Unable to create listener!", err)
	}

	// Database Connection
	dbClient, err := db.ConnectDB()
	if err != nil {
		log.Fatalln("unable to connect to Database!", err)
	}
	fmt.Println("Connected to Database...")

	// User gRPC Server
	userServer := grpc.NewServer()
	userspb.RegisterUsersServer(userServer, &server{})
	reflection.Register(userServer)

	go func() {
		fmt.Println("Running Users gRPC Server at 3000...")
		if err := userServer.Serve(lis); err != nil {
			log.Fatalln("unable to serve!", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, os.Kill)
	<-ch
	fmt.Println("Stopping gRPC Server...")
	userServer.Stop()
	fmt.Println("Stopping Listener...")
	lis.Close()
	fmt.Println("Stopping Database...")
	dbClient.Disconnect(context.Background())
}
