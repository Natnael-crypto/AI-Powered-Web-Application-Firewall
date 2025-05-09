import axios from '../lib/axios'

export const getMapStat = async () => {
  const response = await axios.get('/api/requests/all-blocked-countries')
  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.requests
}
export const getOSStat = async () => {
  const response = await axios.get('/api/requests/os-stats')
  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.requests
}
export const getResponseStatusCodeStat = async () => {
  const response = await axios.get('/api/requests/response-status-stats')
  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.requests
}
