package mappers

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/models"
	nDTOs "github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/notedtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmappers"
)

func MapCoreNoteDetailsAndNoteToNoteDetailsDto(
	coreNoteDetailsDto *nDTOs.CoreNoteDetailsDto,
	note *models.Note,
	noteDetailsDto *nDTOs.NoteDetailsDto,
) {
	sharedmappers.MapMongoModelToBaseCrudObject(note, &(noteDetailsDto.BaseCRUDObject))
	noteDetailsDto.CoreNoteDetailsDto = *coreNoteDetailsDto
}

func MapTextPreviewAndCoreNoteAndNoteToNoteReadDto(
	textPreview string,
	coreNoteDto *nDTOs.CoreNoteDto,
	note *models.Note,
	noteReadDto *nDTOs.NoteReadDto,
) {
	sharedmappers.MapMongoModelToBaseCrudObject(note, &(noteReadDto.BaseCRUDObject))
	noteReadDto.CoreNoteDto = *coreNoteDto
	noteReadDto.TextPreview = textPreview
}
