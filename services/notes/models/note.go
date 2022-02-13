package models

import (
	"context"

	"github.com/NikhilSharma03/noteman-go-grpc/services/notes/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID          primitive.ObjectID `bson:"_id",omitempty`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Creator     string             `json:"creator"`
}

func NewNote(title, description, creator string) *Note {
	return &Note{Title: title, Description: description, Creator: creator}
}

func (n *Note) CreateNote() error {
	noteColl := db.GetCollection("notes")
	n.ID = primitive.NewObjectID()

	_, err := noteColl.InsertOne(context.Background(), n)
	if err != nil {
		return err
	}

	return nil
}

type Notes struct {
	NotesList []Note `json:"notes"`
}

func NewNotes() *Notes {
	return &Notes{}
}

func (n *Notes) GetNotes() error {
	noteColl := db.GetCollection("notes")

	data, err := noteColl.Find(context.Background(), bson.D{})
	if err != nil {
		return err
	}

	var noteData []Note

	for data.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem Note
		err := data.Decode(&elem)
		if err != nil {
			return err
		}

		noteData = append(noteData, elem)
	}

	n.NotesList = noteData

	return nil
}

func (n *Notes) GetNotesByEmail(email string) error {
	noteColl := db.GetCollection("notes")

	data, err := noteColl.Find(context.Background(), bson.D{{"creator", email}})
	if err != nil {
		return err
	}

	var noteData []Note

	for data.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem Note
		err := data.Decode(&elem)
		if err != nil {
			return err
		}

		noteData = append(noteData, elem)
	}

	n.NotesList = noteData

	return nil
}
