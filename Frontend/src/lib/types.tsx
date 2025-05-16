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

export interface SecurityRule {
  rule_id: string
  rule_type: string
  rule_method: string
  rule_definition: RuleDefinitionItem[] | string // Can be stringified JSON or parsed array
  action: string
  application_id: string
  rule_string: string
  created_by: string
  created_at: string // ISO date string
  updated_at: string // ISO date string
  is_active: boolean
  category: string
}

export type UserRole = 'super_admin' | 'admin'
export type UserStatus = 'active' | 'inactive' | 'suspended'

export interface User {
  user_id: string
  username: string
  password_hash: string
  role: UserRole
  status: UserStatus
  created_at: string
  updated_at: string
  last_login: string
  profile_image_url: string | null
}

export interface AuthResponse {
  user: User
}

export enum Roles {
  SUPER_ADMIN = 'super_admin',
  ADMIN = 'admin',
}

export type DashboardOverAllStats = {
  ai_based_detections: number
  blocked_requests: number
  malicious_ips_blocked: number
  rule_based_detections: number
  total_requests: number
}

export type AdminUser = {
  user_id: string
  username: string
  password_hash: string
  role: string
  status: string
  created_at: string // ISO date string
  updated_at: string // ISO date string
  last_login: string // ISO date string
  profile_image_url: string
  notifications_enabled: boolean
}

export type AdminsResponse = {
  admins: AdminUser[]
}

export type Assignment = {
  id: string
  user_id: string
  application_name: string
  application_id: string
}

export type AssignmentsResponse = {
  assignments: Assignment[]
}

interface RuleDefinitionItem {
  rule_type: string
  rule_method: string
  rule_definition: string
}

export type Filter ={
  client_ip: string
  request_method: string
  request_url: string
  threat_type: string
  user_agent: string
  geo_location:string
  threat_detected: boolean
  bot_detected:boolean
  rate_limited: boolean
  start_date: string
  timestamp: string
  end_date: string
  last_hours: string
  body:string
  response_code: string
  rule_detected: string
  ai_result:boolean
  ai_threat_type: string
  search:string
}