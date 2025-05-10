import axios from '../lib/axios'
import {Application, ApplicationsResponse} from '../lib/types'

export async function getApplications() {
  const response: {data: ApplicationsResponse} = await axios.get('/api/application', {
    withCredentials: true,
  })
  if (!response) throw new Error('Something went wrong!')

  return response.data.applications
}

export async function getApplication(applicationId: string) {
  const response = await axios.get(`/application/${applicationId}`)
  if (!response) throw new Error('Something went wrong!')

  return response.data.application
}

export async function createApplication(data: Partial<Application>) {
  const response = await axios.post('/application/add', data)

  if (!response) throw new Error('Something went wrong!')

  return response.status
}
export async function updateApplication(data: Partial<Application>) {
  const response = await axios.put(`/api/application/${data.application_id}`, data)

  if (!response) throw new Error('Something went wrong!')

  return response.status
}
