import axios from '../lib/axios'

export const loginUser = async (userData: {username: string; password: string}) => {
  console.log('here')
  const response = await axios.post('/api/login', userData)
  if (!response) throw new Error('Failed to get users')

  return await response.data
}

export const getUsers = async () => {
  const response = await axios.get('/api/users')

  if (!response) throw new Error('Failed to get users')
  return await response.data
}

export const getuser = async (username: string) => {
  const response = await axios.get(`/users/${username}`)

  if (!response) throw new Error('Failed to get user')

  return await response.data
}

export const addUser = async (userData: {username: string; password: string}) => {
  const response = await axios.post('/users/add', userData)

  if (!response) throw new Error('Failed to add user')

  return await response.data
}

export const updateUsers = async (userId: string) => {
  // const response = axios.post()
}
