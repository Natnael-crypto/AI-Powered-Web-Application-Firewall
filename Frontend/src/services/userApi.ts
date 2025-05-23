import axios from '../lib/axios'

export const loginUser = async (userData: {username: string; password: string}) => {
  const response = await axios.post('/api/login', userData, {
    withCredentials: true, // This is crucial
  })
  if (!response) throw new Error('Failed to get users')
  return response.data
}
export const getUsers = async () => {
  const response = await axios.get('/api/users/')

  if (!response) throw new Error('Failed to get users')
  return response.data
}

export const getAllUser = async () => {
  const response = await axios.get('/api/users/')

  if (!response) throw new Error('Failed to get users')
  return response.data.admins
}

export const getuser = async (username: string) => {
  const response = await axios.get(`/api/users/${username}`)

  if (!response) throw new Error('Failed to get user')

  return await response.data.admins
}


export const getUserById = async (id: string) => {
  const response = await axios.get(`/api/users/id/${id}`)

  if (!response) throw new Error('Failed to get user')

  return await response.data.admin
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

export async function deleteUser(username: string) {
  const response = await axios.delete(`/api/users/delete/${username}`)

  if (!response) throw new Error('Failed to delete user')

  return await response.data
}
export async function deActivateUser(username: string) {
  const response = await axios.put(`/api/users/inactive/${username}`)

  if (!response) throw new Error('Failed to delete user')

  return await response.data
}
export async function activateUser(username: string) {
  const response = await axios.put(`/api/users/active/${username}`)

  if (!response) throw new Error('Failed to delete user')

  return await response.data
}
