package services

type NoteService interface {
	// Todo: define methods
}

type NoteServiceImpl struct {
}

func NewNoteServiceImpl() *NoteServiceImpl {
	return &NoteServiceImpl{}
}
