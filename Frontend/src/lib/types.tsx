export type FilterValues = {
  ipAddress: string
  port: string
  domain: string
  startAt: string
  endAt: string
}

export type AttackLog = {
  ipAddress: string
  application: string
  attackCount: number
  duration: string
  startAt: string
}

export type LogTable = {
  action: string
  url: string
  attackType: string
  ipAddress: string
  time: string
}

export enum logFilterType {
  CLIENTIP,
  COUNTRY,
  METHOD,
  PROTOCOL,
}

export enum filterOperations {
  EQUALS_TO,
  GREATER_THAN,
  LESS_THAN,
  NOT_EQUAL,
}

// types/toast.ts
export type ToastType = 'success' | 'error' | 'info' | 'warning' | 'loading'

export interface Toast {
  id: string
  message: string
  type: ToastType
  createdAt: number
  duration?: number
}

export interface ToastContextType {
  toasts: Toast[]
  addToast: (message: string, type?: ToastType, duration?: number) => void
  removeToast: (id: string) => void
  clearToasts: () => void
}

export type ApplicationConfig = {
  id: string
  application_id: string
  rate_limit: number
  window_size: number
  block_time: number
  detect_bot: boolean
  hostname: string
  max_post_data_size: number
  tls: boolean
}

export type Application = {
  application_id: string
  application_name: string
  config: ApplicationConfig
  created_at: string // or Date if you'll parse it
  description: string
  hostname: string
  ip_address: string
  port: string
  status: boolean
  tls: boolean
  updated_at: string // or Date if you'll parse it
}

export type ApplicationsResponse = {
  applications: Application[]
}
