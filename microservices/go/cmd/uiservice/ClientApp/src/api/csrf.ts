import axios from '../axios'
import {AxiosResponse} from "axios";

export function getCsrfToken(): Promise<AxiosResponse<SuccessDto>> {
  return axios.get<SuccessDto>(`/csrf`)
}