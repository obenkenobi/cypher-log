import axios from '../axios'

export function getCsrfToken(): Promise<void> {
  return axios.get(`/csrf`)
}