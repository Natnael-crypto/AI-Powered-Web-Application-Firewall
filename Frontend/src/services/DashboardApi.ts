import axios from '../lib/axios'
import {DashboardOverAllStats} from '../lib/types'

export const getMapStat = async (appId: string,time:any) => {
  const response = await axios.get(`/api/requests/all-blocked-countries?application_id=${appId}&start_date=${time.start}&end_date=${time.end}`)
  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.requests
}
export const getOSStat = async (appId: string,time:any) => {
  const response = await axios.get(`/api/requests/os-stats?application_id=${appId}&start_date=${time.start}&end_date=${time.end}`)
  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.requests
}

export const getResponseStatusCodeStat = async (appId: string,time:any) => {
  const response = await axios.get(`/api/requests/response-status-stats?application_id=${appId}&start_date=${time.start}&end_date=${time.end}`)
  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.requests
}
export const getMostTargetedEndpoint = async (appId: string,time:any) => {
  const response = await axios.get(`/api/requests/most-targeted-endpoints?application_id=${appId}&start_date=${time.start}&end_date=${time.end}`)
  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.most_targeted_endpoints
}
export const getTopAttackTypes = async (appId: string,time:any) => {
  const response = await axios.get(`/api/requests/top-attack-types?application_id=${appId}&start_date=${time.start}&end_date=${time.end}`)
  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.top_threat_type
}

export const getOverAllStat = async (appId: string,time:any): Promise<DashboardOverAllStats> => {
  const response = await axios.get(`/api/requests/overall-stat?application_id=${appId}&start_date=${time.start}&end_date=${time.end}`)
  if (!response) throw new Error(`Something went wrong!`)
  return await response.data
}

export const getRateStat = async (appId: string,time:any): Promise<DashboardOverAllStats> => {
  const response = await axios.get(`/api/requests/requests-per-minute?application_id=${appId}&start_date=${time.start}&end_date=${time.end}`)
  if (!response) throw new Error(`Something went wrong!`)
  return await response.data
}