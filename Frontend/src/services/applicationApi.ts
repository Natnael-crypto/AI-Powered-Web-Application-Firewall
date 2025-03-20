import axios from '../lib/axios'

type applicationBody = {
  application_name: string
  description: string
  hostname: string
  ip_address: string
  port: string
  status: boolean
  tls: boolean
}

export async function getApplications() {
  const response = await axios.get('/api/application')
  if (!response) throw new Error('Something went wrong!')

  return response.data.applications
}

export async function getApplication(applicationId: string) {
  const response = await axios.get(`/api/application/${applicationId}`)
  if (!response) throw new Error('Something went wrong!')

  return response.data.application
}

export async function createApplication(data: applicationBody) {
  const response = await axios.post('/api/application/add', data)

  if (!response) throw new Error('Something went wrong!')

  return response.status
}
