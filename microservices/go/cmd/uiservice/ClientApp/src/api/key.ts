import axios from '../axios'
import {AxiosResponse} from "axios";

const prefix = "/api/keyservice"

export function userKeyExists(): Promise<AxiosResponse<ExistsDto>> {
  return axios.get<ExistsDto>(`${prefix}/v1/userKey/exists`)
}

export function createUserKey(payload: PasscodeCreateDto): Promise<AxiosResponse<SuccessDto>> {
  return axios.post<SuccessDto, AxiosResponse<SuccessDto>, PasscodeCreateDto>(
    `${prefix}/v1/userKey/passcode`, payload)
}

export function newKeySession(payload: PasscodeDto): Promise<AxiosResponse<SuccessDto>> {
  return axios.post<SuccessDto, AxiosResponse<SuccessDto>, PasscodeDto>(
    `${prefix}/v1/userKey/newSession`, payload)
}