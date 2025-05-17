import {Rule} from '../components/RuleDetailModal'
import axios from '../lib/axios'

export async function getRules() {
  const response = await axios.get('/api/rule', {
    withCredentials: true,
  })
  if (!response) throw new Error('Something went wrong!')

  return response.data.rules
}

export async function createRule(data: any) {
  const response = await axios.post('/api/rule/add', data)

  if (!response) throw new Error('Something went wrong!')

  return response.data
}

export async function updateRule(data: Partial<Rule>) {
  const response = await axios.put(`/api/rule/update/${data.rule_id}`, data)

  if (!response) throw new Error('Something went wrong!')

  return response.data
}

export async function activateRule(ruleId: string) {
  const response = await axios.get(`/api/rule/activate/${ruleId}`)

  if (!response) throw new Error('Something went wrong!')

  return response.data
}
export async function deactivateRule(ruleId: string) {
  const response = await axios.get(`/api/rule/deactivate/${ruleId}`)

  if (!response) throw new Error('Something went wrong!')

  return response.data
}
export async function deleteRule(ruleId: string) {
  const response = await axios.delete(`/api/rule/delete/${ruleId}`)

  if (!response) throw new Error('Something went wrong!')

  return response.data
}
