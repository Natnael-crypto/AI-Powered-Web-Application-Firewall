import axios from '../lib/axios'
import {AdminsResponse} from '../lib/types'

export const loginUser = async (userData: {username: string; password: string}) => {
  const response = await axios.post('/api/login', userData, {
    withCredentials: true, // This is crucial
  })
  if (!response) throw new Error('Failed to get users')
  return response.data
}
export const getUsers = async (): Promise<AdminsResponse> => {
  const response = await axios.get<AdminsResponse>('/api/users/', {withCredentials: true})

  if (!response) throw new Error('Failed to get users')
  return response.data
}

export const getuser = async (username: string) => {
  const response = await axios.get(`/users/${username}`)

  if (!response) throw new Error('Failed to get user')

  return await response.data.admins
}

export const addUser = async (userData: {username: string; password: string}) => {
  const response = await axios.post('/api/users/add', userData)

  if (!response) throw new Error('Failed to add user')

  return await response.data
}

export const updateUsers = async (_userId: string) => {
  // Todo: Implement
}

export const isLoggedIn = async () => {
  const response = await axios.get('/api/is-logged-in')

  if (!response) throw new Error('Failed to add user')

  return await response.data.user
}
