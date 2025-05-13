import axios from '../lib/axios'

export async function getAIModels() {
  const response = await axios.get('/api/models', {
    withCredentials: true,
  })
  if (!response) throw new Error('Something went wrong!')

  return response.data.models
}

export async function createAIModel(data: any) {
  const response = await axios.post('/api/model/train', data)
  if (!response) throw new Error('Something went wrong!')
  return response.data
}

export async function selectAIModel(modelId: string) {
  const response = await axios.post(`/api/model/select/${modelId}`, {}, {
    withCredentials: true,
  })
  if (!response) throw new Error('Something went wrong!')

  return response.data
}
