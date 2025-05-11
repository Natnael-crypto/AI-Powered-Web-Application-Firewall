import axios from '../lib/axios'
export const getRequests = async (application_name: string) => {
  const response = await axios.get('/api/requests', {
    withCredentials: true,
    params: {application_name},
  })
  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.requests
}

export const getDeviceStat = async () => {
  const response = await axios.get('/api/requests/os-stats', {
    withCredentials: true,
  })

  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.os_statistics
}
