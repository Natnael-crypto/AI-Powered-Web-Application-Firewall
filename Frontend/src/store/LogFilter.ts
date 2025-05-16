import { create } from 'zustand'
import { Filter } from '../lib/types'

interface FilterState {
  filters: Filter  // Updated to use `Filter` instead of `Partial<Filter>`
  tempFilter: { key: string; value: string }
  setTempFilter: (key: string, value: string) => void
  addFilter: () => void
  removeFilter: (key: string) => void
  clearFilters: () => void
  applyFilters: () => void
  appliedFilters: Filter  // Updated to use `Filter` instead of `Partial<Filter>`
}

export const useLogFilter = create<FilterState>((set, get) => ({
  filters: {
    search: '',
    page: '',
    client_ip: '',
    request_method: '',
    request_url: '',
    threat_type: '',
    user_agent: '',
    geo_location: '',
    threat_detected: '',
    bot_detected: '',
    rate_limited: '',
    start_date: '',
    timestamp: '',
    end_date: '',
    last_hours: '',
    body: '',
    response_code: '',
    rule_detected: '',
    ai_result: '',
    ai_threat_type: '',
    application_name: ''
  },
  tempFilter: { key: '', value: '' },
  setTempFilter: (key, value) =>
    set(state => ({ tempFilter: { ...state.tempFilter, [key]: value } })),
  addFilter: () => {
    const { tempFilter, filters } = get()

    // Skip adding filter if tempFilter is empty
    if (!tempFilter.key || !tempFilter.value) return

    // Skip adding filter if it already exists
    if (filters[tempFilter.key as keyof Filter]) return

    set(state => ({
      filters: {
        ...state.filters,
        [state.tempFilter.key]: state.tempFilter.value,
      },
      tempFilter: { key: '', value: '' },
    }))
  },
  removeFilter: key =>
    set(state => {
      const updated = { ...state.filters }
      delete updated[key]
      return { filters: updated }
    }),
  clearFilters: () => set({ filters: get().filters, appliedFilters: {
    client_ip: '',
    request_method: '',
    request_url: '',
    threat_type: '',
    user_agent: '',
    geo_location: '',
    threat_detected: '',
    bot_detected: '',
    rate_limited: '',
    start_date: '',
    timestamp: '',
    end_date: '',
    last_hours: '',
    body: '',
    response_code: '',
    rule_detected: '',
    ai_result: '',
    ai_threat_type: '',
    search: '',
    page: '',
    application_name: ''
  } }),  // Reset filters, retain appliedFilters
  applyFilters: () => set(state => ({ appliedFilters: { ...state.filters } })),
  appliedFilters: {
    application_name:'',
    client_ip: '',
    request_method: '',
    request_url: '',
    threat_type: '',
    user_agent: '',
    geo_location: '',
    threat_detected: '',
    bot_detected: '',
    rate_limited: '',
    start_date: '',
    timestamp: '',
    end_date: '',
    last_hours: '',
    body: '',
    response_code: '',
    rule_detected: '',
    ai_result: '',
    ai_threat_type: '',
    search: '',
    page: ''
  },
}))
