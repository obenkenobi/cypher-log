package notedtos

import "github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/embedded"

type CoreNoteDto struct {
	Title string `json:"title" binding:"required,min=4,max=1000"`
}

func NewCoreNoteDto(title string) CoreNoteDto {
	return CoreNoteDto{Title: title}
}

type CoreNoteDetailsDto struct {
	CoreNoteDto
	Text string `json:"text" binding:"max=40000"`
}

func NewCoreNoteDetailsDto(title string, text string) CoreNoteDetailsDto {
	return CoreNoteDetailsDto{
		CoreNoteDto: NewCoreNoteDto(title),
		Text:        text,
	}
}

type NoteCreateDto struct {
	CoreNoteDetailsDto
}

type NoteUpdateDto struct {
	embedded.BaseId
	CoreNoteDetailsDto
}

type NotePreviewDto struct {
	embedded.BaseCRUDObject
	CoreNoteDto
	TextPreview string `json:"textPreview"`
}

type NoteReadDto struct {
	embedded.BaseCRUDObject
	CoreNoteDetailsDto
}

type NoteIdDto struct {
	embedded.BaseRequiredId
}
