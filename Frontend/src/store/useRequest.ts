import {create} from 'zustand'

type request = {
  request_id: string
  application_name: string
  client_ip: string
  request_method: string
  request_url: string
  headers: string
  body: string
  timestamp: string
  response_code: number
  status: string
  matched_rules: string
  threat_detected: boolean
  threat_type: string
  bot_detected: boolean
  geo_location: string
  rate_limited: boolean
  user_agent: string
  ai_analysis_result: string
}

interface RequestState {
  requests: request[]
  getRequests: (requests: request[]) => void
}

export const useRequests = create<RequestState>(set => ({
  requests: [],
  getRequests: (requests: request[]) => set({requests: requests}),
}))
