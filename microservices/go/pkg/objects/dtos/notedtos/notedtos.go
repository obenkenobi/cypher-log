package notedtos

import "github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/embedded"

type CoreNoteDto struct {
	Title string `json:"title" binding:"required,alphanumunicode,min=4,max=1000"`
}

type CoreNoteDetailsDto struct {
	CoreNoteDto
	Text *string `json:"text" binding:"max=40000"`
}

func (c CoreNoteDetailsDto) GetText() string {
	if c.Text != nil {
		return *c.Text
	}
	return ""
}

type NoteCreateDto struct {
	CoreNoteDetailsDto
}

type NoteUpdateDto struct {
	embedded.BaseId
	CoreNoteDetailsDto
}

type NoteReadDto struct {
	embedded.BaseCRUDObject
	CoreNoteDto
	TextPreview string
}

type NoteDetailsDto struct {
	embedded.BaseCRUDObject
	CoreNoteDetailsDto
}
