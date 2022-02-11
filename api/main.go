package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/NikhilSharma03/noteman-go-grpc/api/handlers"
	"github.com/NikhilSharma03/noteman-go-grpc/services/notes/notespb"
	"github.com/NikhilSharma03/noteman-go-grpc/services/users/userspb"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func main() {
	log := log.New(os.Stdout, "noteman-api", log.LstdFlags)

	// gRPC Client Connection
	uClient, err := grpc.Dial(":3000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	nClient, err := grpc.Dial(":3001", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	uc := userspb.NewUsersClient(uClient)
	nc := notespb.NewNotesClient(nClient)

	nh := handlers.NewNoteHandler(log, nc)
	uh := handlers.NewUserHandler(log, uc)

	sm := mux.NewRouter()

	// Get Handlers
	getHandlers := sm.Methods(http.MethodGet).Subrouter()
	getHandlers.HandleFunc("/api/notes", nh.GetPost)
	getHandlers.HandleFunc("/api/notes/{email}", nh.GetPostByEmail)

	// Post Handlers
	postHandlers := sm.Methods(http.MethodPost).Subrouter()
	postHandlers.HandleFunc("/api/users/login", uh.Login)
	postHandlers.HandleFunc("/api/users/signup", uh.SignUp)
	postHandlers.HandleFunc("/api/notes", nh.AddPost)

	s := http.Server{
		Addr:         ":8000",
		Handler:      sm,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	fmt.Println("Listening on port 8000")

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
