interface CoreNoteDto {
  title: string
}

interface CoreNoteDetailsDto extends CoreNoteDto {
  text: string
}

interface NoteCreateDto extends CoreNoteDetailsDto {}

interface NoteUpdateDto extends BaseId, CoreNoteDetailsDto {}

interface NotePreviewDto extends BaseCRUDObject, CoreNoteDto {
  textPreview: string
}
interface NoteReadDto extends BaseCRUDObject, CoreNoteDetailsDto {}

interface NoteIdDto extends BaseRequiredId {}

