import axios from '../axios'

const prefix = "/api/userservice"

export function addUser(payload: UserSaveDto): Promise<UserReadDto> {
  return axios.post<any, UserReadDto, UserSaveDto>(`${prefix}/v1/user`, payload)
}

export function updateUser(payload: UserSaveDto): Promise<UserReadDto> {
  return axios.put<any, UserReadDto, UserSaveDto>(`${prefix}/v1/user`, payload)
}

export function deleteUser(): Promise<UserReadDto> {
  return axios.delete<any, UserReadDto, any>(`${prefix}/v1/user`)
}

export function getIdentity(): Promise<UserIdentityDto> {
  return axios.get<any, UserIdentityDto, any>(`${prefix}/v1/user`)
}