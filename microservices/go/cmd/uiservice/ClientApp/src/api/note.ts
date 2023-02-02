import axios from '../axios'
import {AxiosResponse} from "axios";

const prefix = "/api/noteservice"

export function addNote(payload: NoteCreateDto): Promise<AxiosResponse<SuccessDto>> {
  return axios.post<SuccessDto, AxiosResponse<SuccessDto>, NoteCreateDto>(`${prefix}/v1/notes`, payload)
}

export function updateNote(payload: NoteUpdateDto): Promise<AxiosResponse<SuccessDto>> {
  return axios.put<SuccessDto, AxiosResponse<SuccessDto>, NoteUpdateDto>(`${prefix}/v1/notes`, payload)
}
export function deleteNote(payload: NoteIdDto): Promise<AxiosResponse<SuccessDto>> {
  return axios.delete<SuccessDto, AxiosResponse<SuccessDto>, NoteIdDto>(`${prefix}/v1/notes`, {data: payload})
}

export function getNoteById(payload: NoteIdDto): Promise<AxiosResponse<NoteReadDto>> {
  return axios.post<NoteReadDto, AxiosResponse<NoteReadDto>, NoteIdDto>(`${prefix}/v1/notes/getById`, payload)
}

export function getNotesPage(payload: PageRequest): Promise<AxiosResponse<Page<NotePreviewDto>>> {
  return axios.post<Page<NotePreviewDto>, AxiosResponse<Page<NotePreviewDto>>, PageRequest>(
    `${prefix}/v1/notes/getById`, payload)
}