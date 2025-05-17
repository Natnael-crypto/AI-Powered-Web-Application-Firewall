import axios from '../lib/axios'
import { Filter } from '../lib/types'
export const getRequests = async (filter: Partial<Filter>) => {
  const response = await axios.get(`/api/requests?page=${filter.page}&client_ip=${filter.client_ip}&request_method=${filter.request_method}&request_url=${filter.request_url}&threat_type=${filter.threat_type}&user_agent=${filter.user_agent}&geo_location=${filter.geo_location}&threat_detected=${filter.threat_detected}&bot_detected=${filter.bot_detected}&rate_limited=${filter.rate_limited}&start_date=${filter.start_date}&end_date=${filter.end_date}&last_hours=${filter.last_hours}&body=${filter.body}&response_code=${filter.response_code}&rule_detected=${filter.rule_detected}&ai_result=${filter.ai_result}&ai_threat_type=${filter.ai_threat_type}&search=${filter.search}`)
  if (!response) throw new Error(`Something went wrong!`)
  return await response.data
}

export const getDeviceStat = async (selectedApp:string,timeRange:any) => {
  const response = await axios.get(`/api/requests/os-stats?application_id=${selectedApp}&start_date=${timeRange.start}&end_date=${timeRange.end}`)
  if (!response) throw new Error(`Something went wrong!`)
  return await response.data.os_statistics
}
