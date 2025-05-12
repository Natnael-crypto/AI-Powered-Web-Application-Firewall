import axios from '../lib/axios'

export async function getAIModels() {
  const response = await axios.get('/models', {
    withCredentials: true,
  })
  if (!response) throw new Error('Something went wrong!')

  return response.data.models
}

export async function createAIModel(data: any) {
  const response = await axios.post('/model/train', data, {
    withCredentials: true,
  })

  if (!response) throw new Error('Something went wrong!')

  return response.data
}

export async function selectAIModel(modelId: string) {
  const response = await axios.post(`/model/select/${modelId}`, {}, {
    withCredentials: true,
  })

  if (!response) throw new Error('Something went wrong!')

  return response.data
}
