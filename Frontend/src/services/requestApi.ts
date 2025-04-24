import axios from '../lib/axios'

export const getRequests = async (application_name: string) => {
  const response = await axios.get(`/api/requests/?application_name=${application_name}`)
  console.log('response: ', response)
  if (!response) throw new Error(`Something went wrong!`)

  return await response.data.requests
}
