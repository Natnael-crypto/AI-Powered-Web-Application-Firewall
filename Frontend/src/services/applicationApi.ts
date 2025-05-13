import axios from '../lib/axios'
import {Application, ApplicationsResponse, AssignmentsResponse} from '../lib/types'

export async function getApplications() {
  const response: {data: ApplicationsResponse} = await axios.get('/api/application', {
    withCredentials: true,
  })
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
  const response = await axios.put(`/api/application/${data.application_id}`, data)

  if (!response) throw new Error('Something went wrong!')

  return response.status
}

export async function getAssignments(): Promise<AssignmentsResponse> {
  const response = await axios.get('/api/application/assignments')

  if (!response) throw new Error('Something went wrong')

  return await response.data
}
export async function assignApplication(data: {
  user_id: string
  application_name: string
}): Promise<AssignmentsResponse> {
  const response = await axios.post('/api/application/assign', data)

  if (!response) throw new Error('Something went wrong')

  return await response.data
}
export async function deleteAssignment(
  assignment_id: string,
): Promise<AssignmentsResponse> {
  const response = await axios.delete(`/api/application/assign/${assignment_id}`)

  if (!response) throw new Error('Something went wrong')

  return await response.data
}
