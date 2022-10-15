package services

import "github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/repositories"

type NoteService interface {
	// Todo: define methods
}

type NoteServiceImpl struct {
	noteRepository repositories.NoteRepository
}

func NewNoteServiceImpl(noteRepository repositories.NoteRepository) *NoteServiceImpl {
	return &NoteServiceImpl{noteRepository: noteRepository}
}
