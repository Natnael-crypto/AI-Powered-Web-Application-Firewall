import axios from '../lib/axios'
export const getRequests = async (application_name: string) => {
  const response = await axios.get('/api/requests', {
    withCredentials: true,
    params: {application_name},
  })
  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.requests
}

export const getDeviceStat = async (selectedApp:string,timeRange:any) => {
  const response = await axios.get(`/api/requests/os-stats?application_id=${selectedApp}&start_date=${timeRange.start}&end_date=${timeRange.end}`, {
    withCredentials: true,
  })

  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.os_statistics
}
