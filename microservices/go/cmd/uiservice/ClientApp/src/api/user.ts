import axios from '../axios'
import { AxiosResponse } from 'axios'

const prefix = "/api/userservice"

export function addUser(payload: UserSaveDto): Promise<AxiosResponse<UserReadDto>> {
  return axios.post<UserReadDto, AxiosResponse<UserReadDto>, UserSaveDto>(`${prefix}/v1/user`, payload)
}

export function updateUser(payload: UserSaveDto): Promise<AxiosResponse<UserReadDto>> {
  return axios.put<UserReadDto, AxiosResponse<UserReadDto>, UserSaveDto>(`${prefix}/v1/user`, payload)
}

export function deleteUser(): Promise<AxiosResponse<UserReadDto>> {
  return axios.delete<UserReadDto>(`${prefix}/v1/user`)
}

export function getIdentity(): Promise<AxiosResponse<UserIdentityDto>> {
  return axios.get<UserIdentityDto>(`${prefix}/v1/user/me`)
}