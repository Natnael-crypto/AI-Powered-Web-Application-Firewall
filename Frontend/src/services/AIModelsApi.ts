import axios from '../lib/axios'

export async function getAIModels() {
  const response = await axios.get('/api/models', {
    withCredentials: true,
  })
  if (!response) throw new Error('Something went wrong!')

  return response.data.model
}

export async function updateAIModelSetting(data: any) {
  const response = await axios.post('/api/model/update/setting', data)
  if (!response) throw new Error('Something went wrong!')
  return response.data
}

export async function updateAIModelTrainTime(data: any) {
  const response = await axios.post('/api/model/update/time', data)
  if (!response) throw new Error('Something went wrong!')
  return response.data
}

export async function selectAIModel(modelId: string) {
  const response = await axios.get(`/api/model/select/${modelId}`)
  if (!response) throw new Error('Something went wrong!')

  return response.data
}

