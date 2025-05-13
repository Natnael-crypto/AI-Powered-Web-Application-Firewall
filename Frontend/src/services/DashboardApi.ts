import axios from '../lib/axios'
import {DashboardOverAllStats} from '../lib/types'

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
export const getMostTargetedEndpoint = async () => {
  const response = await axios.get('/api/requests/most-targeted-endpoints')
  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.most_targeted_endpoints
}
export const getTopAttackTypes = async () => {
  const response = await axios.get('/api/requests/top-attack-types')
  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.top_threat_type
}

export const getOverAllStat = async (appId: string): Promise<DashboardOverAllStats> => {
  const response = await axios.get(`/api/requests/overall-stat/${appId}`)
  if (!response) throw new Error(`Something went wrong!`)

  return await response.data
}
