import axios from '../axios'
//Todo: handle errors

const prefix = "/api/userservice"

export function addUser(req: UserSaveDto): Promise<UserReadDto> {
  return axios.post<UserSaveDto, UserReadDto>(`${prefix}/v1/user`, req)
}

export function updateUser(req: UserSaveDto): Promise<UserReadDto> {
  return axios.put<UserSaveDto, UserReadDto>(`${prefix}/v1/user`, req)
}

export function deleteUser(): Promise<UserReadDto> {
  return axios.delete<any, UserReadDto>(`${prefix}/v1/user`)
}

export function getIdentity(): Promise<UserIdentityDto> {
  return axios.get<any, UserIdentityDto>(`${prefix}/v1/user`)
}