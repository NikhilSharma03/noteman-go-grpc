package handlers

import (
	"log"
	"net/http"

	"github.com/NikhilSharma03/noteman-go-grpc/services/notes/notespb"
)

type NoteHandler struct {
	l *log.Logger
	n notespb.NotesClient
}

func NewNoteHandler(l *log.Logger, n notespb.NotesClient) *NoteHandler {
	return &NoteHandler{l, n}
}

func (n *NoteHandler) GetPost(rw http.ResponseWriter, rq *http.Request) {

}

func (n *NoteHandler) GetPostByEmail(rw http.ResponseWriter, rq *http.Request) {

}

func (n *NoteHandler) AddPost(rw http.ResponseWriter, rq *http.Request) {

}
