import axios from '../lib/axios'
export const getSysConf = async () => {
  const response = await axios.get('/api/requests')
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
