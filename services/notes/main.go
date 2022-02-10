package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/NikhilSharma03/noteman-go-grpc/services/notes/db"
	"github.com/NikhilSharma03/noteman-go-grpc/services/notes/helpers"
	"github.com/NikhilSharma03/noteman-go-grpc/services/notes/models"
	"github.com/NikhilSharma03/noteman-go-grpc/services/notes/notespb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// Server struct
type server struct {
	notespb.UnimplementedNotesServer
}

func (*server) CreateNote(ctx context.Context, rq *notespb.CreateNoteRequest) (*notespb.CreateNoteResponse, error) {
	token, creator, title, description := rq.GetToken(), rq.GetNoteDetails().GetCreator(), rq.GetNoteDetails().GetTitle(), rq.GetNoteDetails().GetDescription()

	// Token validation
	tokenEmail, err := helpers.ExtractTokenMetadata(token)
	if err != nil || tokenEmail.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid token")
	}

	// User validation
	if tokenEmail.Email != creator {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated user")
	}

	//  Data validation
	if creator == "" || title == "" || description == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid token")
	}

	note := models.NewNote(title, description, creator)
	err = note.CreateNote()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while creating note")
	}

	return &notespb.CreateNoteResponse{
		Message: "Note created successfully",
		NoteDetails: &notespb.Note{
			Id:          note.ID.Hex(),
			Title:       note.Title,
			Creator:     note.Creator,
			Description: note.Description,
		},
	}, nil
}

func (*server) GetNote(ctx context.Context, rq *notespb.NoteRequest) (*notespb.NoteResponse, error) {
	token := rq.GetToken()

	// Token validation
	tokenEmail, err := helpers.ExtractTokenMetadata(token)
	if err != nil || tokenEmail.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid token")
	}

	notes := models.NewNotes()
	err = notes.GetNotes()
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var noteResult = []*notespb.Note{}

	for _, item := range notes.NotesList {
		newNote := &notespb.Note{
			Id:          item.ID.Hex(),
			Title:       item.Title,
			Description: item.Description,
			Creator:     item.Creator,
		}
		noteResult = append(noteResult, newNote)
	}

	return &notespb.NoteResponse{
		Message:     "Fetched notes successfully",
		NoteDetails: noteResult,
	}, nil
}

func (*server) GetNoteByEmail(ctx context.Context, rq *notespb.GetNoteRequest) (*notespb.GetNoteResponse, error) {
	token, creator := rq.GetToken(), rq.GetEmail()

	// Token validation
	tokenEmail, err := helpers.ExtractTokenMetadata(token)
	if err != nil || tokenEmail.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid token")
	}

	notes := models.NewNotes()
	err = notes.GetNotesByEmail(creator)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var noteResult = []*notespb.Note{}

	for _, item := range notes.NotesList {
		newNote := &notespb.Note{
			Id:          item.ID.Hex(),
			Title:       item.Title,
			Description: item.Description,
			Creator:     item.Creator,
		}
		noteResult = append(noteResult, newNote)
	}

	return &notespb.GetNoteResponse{
		Message:     "Fetched notes successfully",
		NoteDetails: noteResult,
	}, nil
}

func main() {
	// Listener
	lis, err := net.Listen("tcp", ":3001")
	if err != nil {
		log.Fatalln("Unable to create listener!", err)
	}

	// Database Connection
	dbClient, err := db.ConnectDB()
	if err != nil {
		log.Fatalln("unable to connect to Database!", err)
	}
	fmt.Println("Connected to Database...")

	// Note gRPC Server
	noteServer := grpc.NewServer()
	notespb.RegisterNotesServer(noteServer, &server{})
	reflection.Register(noteServer)

	go func() {
		fmt.Println("Running Notes gRPC Server at 3001...")
		if err := noteServer.Serve(lis); err != nil {
			log.Fatalln("unable to serve!", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, os.Kill)
	<-ch
	fmt.Println("Stopping gRPC Server...")
	noteServer.Stop()
	fmt.Println("Stopping Listener...")
	lis.Close()
	fmt.Println("Stopping Database...")
	dbClient.Disconnect(context.Background())
}
