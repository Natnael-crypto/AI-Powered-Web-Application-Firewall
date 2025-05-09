import axios from '../lib/axios'
import {Application, ApplicationsResponse} from '../lib/types'

export async function getApplications() {
  const response: {data: ApplicationsResponse} = await axios.get('/api/application')
  if (!response) throw new Error('Something went wrong!')

  return response.data.applications
}

export async function getApplication(applicationId: string) {
  const response = await axios.get(`/api/application/${applicationId}`)
  if (!response) throw new Error('Something went wrong!')

  return response.data.application
}

export async function createApplication(data: Partial<Application>) {
  const response = await axios.post('/api/application/add', data)

  if (!response) throw new Error('Something went wrong!')

  return response.status
}
export async function updateApplication(data: Partial<Application>) {
  const response = await axios.put(`/application/${data.application_id}`, data)

  if (!response) throw new Error('Something went wrong!')

  return response.status
}
