import axios from '../axios'

const prefix = "/api/keyservice"

export function userKeyExists(): Promise<ExistsDto> {
  return axios.get<any, ExistsDto>(`${prefix}/v1/userKey/exists`)
}

export function createUserKey(dto: PasscodeCreateDto): Promise<SuccessDto> {
  return axios.post<any, SuccessDto, PasscodeCreateDto>(`${prefix}/v1/userKey/passcode`, dto)
}

export function newKeySession(dto: PasscodeDto): Promise<SuccessDto> {
  return axios.post<any, SuccessDto, PasscodeDto>(`${prefix}/v1/userKey/newSession`, dto)
}