import axios from '../lib/axios'
import {
  Application,
  ApplicationsResponse,
  AssignmentsResponse,
  rateLimitInputtype,
} from '../lib/types'

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

export async function updateListeningPort(data: {listening_port: string}) {
  const response = await axios.put(`/api/config/update/listening-port`, data)

  if (!response) throw new Error('Something went wrong')

  return await response.data
}

export async function updateRateLimit(data: {
  application_id: string
  data: rateLimitInputtype
}) {
  const response = await axios.put(
    `/api/config/update/rate-limit/${data.application_id}`,
    data.data,
  )

  if (!response) throw new Error('Something went wrong')

  return await response.data
}

export async function updateRemoteLogServer(data: {remote_logServer: string}) {
  const response = await axios.put(`/api/config/update/remote-log-server`, data)

  if (!response) throw new Error('Something went wrong')

  return await response.data
}

export async function getconfig() {
  const response = await axios.get('/api/config/')

  if (!response) throw new Error('Something went wrong')

  return await response.data.data
}

export async function getApplicationConfig(application_id: string) {
  const response = await axios.get(`/api/config/${application_id}`)

  if (!response) throw new Error('Something went wrong')

  return await response.data.data
}

export async function deleteApplication(application_id: string) {
  const response = await axios.delete(`/api/application/${application_id}`)

  if (!response) throw new Error('Something went wrong')

  return await response.data.data
}

export async function updateDetectBOT(data: {
  application_id: string
  data: {detect_bot: boolean}
}) {
  const response = await axios.put(
    `/api/config/update/detect-bot/${data.application_id}`,
    data.data,
  )

  if (!response) throw new Error('Something went wrong')

  return await response.data.data
}

export async function updateMaxDataSize(data: {
  application_id: string
  data: {max_post_data_size: number}
}) {
  const response = await axios.put(
    `/api/config/update/post-data-size/${data.application_id}`,
    data.data,
  )

  if (!response) throw new Error('Something went wrong')

  return await response.data.data
}

export async function uploadCertificate(data: {
  application_id: string
  certificate: File
  key: File
}) {
  const formData = new FormData()
  formData.append('cert', data.certificate)
  formData.append('key', data.key)

  const response = await axios.post(`/api/certs/${data.application_id}`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })

  if (!response) throw new Error('Something went wrong')

  return await response.data
}

export async function getCertificates(application_id: string) {
  const response = await axios.get(`/api/certs`, {
    params: {application_id: application_id, type: 'cert'},
  })

  if (!response) throw new Error('Something went wrong')

  return await response.data
}
