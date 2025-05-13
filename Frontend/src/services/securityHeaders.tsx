import axios from '../lib/axios'

export async function getSecurityHeaders() {
  const response = await axios.get('/api/security-headers', {
    withCredentials: true,
  })
  if (!response) throw new Error('Something went wrong!')

  return response.data.model
}

export async function createSecurityHeader(data: any) {
  const response = await axios.post('/api/security-headers', data)
  if (!response) throw new Error('Something went wrong!')
  return response.data
}

export async function updateSecurityHeader(headerId: string,data:any) {
  const response = await axios.put(`/api/security-headers/${headerId}`,data)
  if (!response) throw new Error('Something went wrong!')
  return response.data
}

export async function deleteSecurityHeader(headerId: string) {
  const response = await axios.delete(`/api/security-headers/${headerId}`)
  if (!response) throw new Error('Something went wrong!')
  return response.data
}
