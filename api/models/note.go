package models

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID          primitive.ObjectID `bson:"_id",omitempty`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Creator     string             `json:"creator"`
	Token       string             `json:"token"`
}

type Notes struct {
	NotesDetails []*Note `json:"notes_details"`
}

func NewNote() *Note {
	return &Note{}
}

func (n *Note) FromJson(r io.Reader) error {
	jd := json.NewDecoder(r)
	return jd.Decode(n)
}

func (n *Note) ToJson(r io.Writer) error {
	je := json.NewEncoder(r)
	return je.Encode(n)
}

func NewNotes() *Notes {
	return &Notes{}
}

func (n *Notes) AddNote(note *Note) {
	n.NotesDetails = append(n.NotesDetails, note)
}

func (n *Notes) ToJson(r io.Writer) error {
	je := json.NewEncoder(r)
	return je.Encode(n)
}
