import axios from '../lib/axios'

export const getRequests = async (application_name: string) => {
  const response = await axios.get(`/requests/?application_name=${application_name}`)
  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.requests
}
