package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/NikhilSharma03/noteman-go-grpc/api/models"
	"github.com/NikhilSharma03/noteman-go-grpc/services/notes/notespb"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/status"
)

type NoteHandler struct {
	l *log.Logger
	n notespb.NotesClient
}

func NewNoteHandler(l *log.Logger, n notespb.NotesClient) *NoteHandler {
	return &NoteHandler{l, n}
}

func (n *NoteHandler) GetNotes(rw http.ResponseWriter, rq *http.Request) {
	n.l.Println("Add GetNotes Called..")
	token := rq.Header.Get("Authorization")
	if token == "" {
		http.Error(rw, "no token found in auth header", http.StatusUnauthorized)
		return
	}

	res, err := n.n.GetNote(context.Background(), &notespb.NoteRequest{Token: token})
	e, _ := status.FromError(err)
	if err != nil {
		http.Error(rw, e.Message(), http.StatusInternalServerError)
		return
	}

	notes := models.NewNotes()
	for _, item := range res.GetNoteDetails() {
		newNote := models.NewNote()
		id, _ := primitive.ObjectIDFromHex(item.Id)
		newNote.ID = id
		newNote.Creator = item.Creator
		newNote.Title = item.Title
		newNote.Description = item.Description
		notes.AddNote(newNote)
	}

	err = notes.ToJson(rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (n *NoteHandler) GetNotesByEmail(rw http.ResponseWriter, rq *http.Request) {
	n.l.Println("Add GetNotesByEmail Called..")
	email := mux.Vars(rq)["email"]
	if email == "" {
		http.Error(rw, "no email found", http.StatusUnauthorized)
		return
	}

	token := rq.Header.Get("Authorization")
	if token == "" {
		http.Error(rw, "no token found in auth header", http.StatusUnauthorized)
		return
	}

	res, err := n.n.GetNoteByEmail(context.Background(), &notespb.GetNoteRequest{Token: token, Email: email})
	e, _ := status.FromError(err)
	if err != nil {
		http.Error(rw, e.Message(), http.StatusInternalServerError)
		return
	}

	notes := models.NewNotes()
	for _, item := range res.GetNoteDetails() {
		newNote := models.NewNote()
		id, _ := primitive.ObjectIDFromHex(item.Id)
		newNote.ID = id
		newNote.Creator = item.Creator
		newNote.Title = item.Title
		newNote.Description = item.Description
		notes.AddNote(newNote)
	}

	err = notes.ToJson(rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (n *NoteHandler) AddNote(rw http.ResponseWriter, rq *http.Request) {
	n.l.Println("Add Note Called..")
	note := models.NewNote()
	err := note.FromJson(rq.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := n.n.CreateNote(context.Background(), &notespb.CreateNoteRequest{Token: note.Token, NoteDetails: &notespb.Note{Title: note.Title, Creator: note.Creator, Description: note.Description}})
	e, _ := status.FromError(err)
	if err != nil {
		http.Error(rw, e.Message(), http.StatusInternalServerError)
		return
	}

	rid, err := primitive.ObjectIDFromHex(res.GetNoteDetails().GetId())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	note.ID = rid

	err = note.ToJson(rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
