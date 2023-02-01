import axios from '../axios'

const prefix = "/api/noteservice"

export function addNote(payload: NoteCreateDto): Promise<SuccessDto> {
  return axios.post<any, SuccessDto, NoteCreateDto>(`${prefix}/v1/notes`, payload)
}

export function updateNote(payload: NoteUpdateDto): Promise<SuccessDto> {
  return axios.put<any, SuccessDto, NoteUpdateDto>(`${prefix}/v1/notes`, payload)
}
export function deleteNote(payload: NoteIdDto): Promise<SuccessDto> {
  return axios.delete<any, SuccessDto, NoteIdDto>(`${prefix}/v1/notes`, {data: payload})
}

export function getNoteById(payload: NoteIdDto): Promise<NoteReadDto> {
  return axios.post<any, NoteReadDto, NoteIdDto>(`${prefix}/v1/notes/getById`, payload)
}

export function getNotesPage(payload: PageRequest): Promise<Page<NotePreviewDto>> {
  return axios.post<any, Page<NotePreviewDto>, PageRequest>(`${prefix}/v1/notes/getById`, payload)
}