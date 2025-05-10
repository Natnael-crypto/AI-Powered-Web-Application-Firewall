import axios from '../lib/axios'

export async function getRules() {
  const response = await axios.get('/api/rule', {
    withCredentials: true,
  })
  if (!response) throw new Error('Something went wrong!')

  return response.data.rules
}

export async function activateRule(ruleId: string) {
  const response = await axios.get(`/api/rule/activate/${ruleId}`)

  if (!response) throw new Error('Something went wrong!')

  return response.status
}
export async function deactivateRule(ruleId: string) {
  const response = await axios.get(`/api/rule/deactivate/${ruleId}`)

  if (!response) throw new Error('Something went wrong!')

  return response.status
}
export async function deleteRule(ruleId: string) {
  const response = await axios.delete(`/api/rule/delete/${ruleId}`)

  if (!response) throw new Error('Something went wrong!')

  return response.status
}
