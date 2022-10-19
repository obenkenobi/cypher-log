package mappers

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/noteservice/models"
	nDTOs "github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/notedtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedmappers"
)

func MapCoreNoteDetailsAndNoteToNoteReadDto(
	coreNoteDetailsDto *nDTOs.CoreNoteDetailsDto,
	note *models.Note,
	noteReadDto *nDTOs.NoteReadDto,
) {
	sharedmappers.MapMongoModelToBaseCrudObject(note, &(noteReadDto.BaseCRUDObject))
	noteReadDto.CoreNoteDetailsDto = *coreNoteDetailsDto
}

func MapTextPreviewAndCoreNoteAndNoteToNotePreviewDto(
	textPreview string,
	coreNoteDto *nDTOs.CoreNoteDto,
	note *models.Note,
	notePreviewDto *nDTOs.NotePreviewDto,
) {
	sharedmappers.MapMongoModelToBaseCrudObject(note, &(notePreviewDto.BaseCRUDObject))
	notePreviewDto.CoreNoteDto = *coreNoteDto
	notePreviewDto.TextPreview = textPreview
}
