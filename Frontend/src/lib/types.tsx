import {z} from 'zod'

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

export const SenderEmailSchema = z.object({
  sender_email: z
    .string()
    .email('Invalid email address')
    .nonempty('Sender email is required'),
  app_password: z.string().nonempty('App password is required'),
})

export type SenderEmail = z.infer<typeof SenderEmailSchema>

export type Filter = {
  client_ip: string
  request_id: string
  request_method: string
  request_url: string
  threat_type: string
  user_agent: string
  geo_location: string
  threat_detected: string
  bot_detected: string
  rate_limited: string
  start_date: string
  timestamp: string
  end_date: string
  last_hours: string
  body: string
  response_code: string
  rule_detected: string
  ai_result: string
  ai_threat_type: string
  search: string
  page: string
  application_name: string
  [key: string]: string | undefined
}

export const logFilterType = {
  application_name: 'Application Name',
  request_id: 'Request ID',
  search: 'Search',
  client_ip: 'Client IP',
  request_method: 'Request Method',
  request_url: 'Request URL',
  threat_type: 'Threat Type',
  user_agent: 'User Agent',
  geo_location: 'Geo Location',
  threat_detected: 'Threat Detected',
  bot_detected: 'Bot Detected',
  rate_limited: 'Rate Limited',
  body: 'Body',
  response_code: 'Response Code',
  rule_detected: 'Rule Detected',
  ai_result: 'AI Result',
  ai_threat_type: 'AI Threat Type',
  last_hours: 'Last Hours',
  start_date: 'Start Date',
  end_date: 'End Date',
} as const

export type FilterKey = keyof typeof logFilterType

export type AiModelSetting = {
  id: string
  expected_accuracy: number
  expected_precision: number
  expected_recall: number
  expected_f1: number
  train_every: number
}

export type Condition = {
  rule_type: string
  rule_method: string
  rule_definition: string
}

export type RuleInput = {
  ruleID: string
  action: string
  category: string
  conditions: Condition[]
  applications: string[]
}

export type AppOption = {
  application_id: string
  application_name: string
}

export type RuleResponse = {
  rule_id: string
  rule_type: string
  rule_method: string
  rule_definition: Condition[]
  action: string
  rule_string: string
  created_at: string
  updated_at: string
  is_active: boolean
  applications: AppOption[]
  category: string
}

export const validRuleTypes = [
  'REQUEST_HEADERS',
  'REQUEST_URI',
  'ARGS',
  'ARGS_GET',
  'ARGS_POST',
  'REQUEST_COOKIES',
  'REQUEST_BODY',
  'XML',
  'JSON',
  'REQUEST_METHOD',
  'REQUEST_PROTOCOL',
  'REMOTE_ADDR',
]

export const validRuleMethods = [
  'regex',
  'streq',
  'contains',
  'ipMatch',
  'rx',
  'beginsWith',
  'endsWith',
  'eq',
  'pm',
]

export const validActions = [
  'deny',
  'log',
  'nolog',
  'pass',
  'drop',
  'redirect',
  'capture',
  't:none',
  't:lowercase',
  't:normalizePath',
  't:urlDecode',
  't:compressWhitespace',
  'severity:2',
  'severity:3',
  'status:403',
]

export interface Notification {
  notification_id: string
  user_id: string
  message: string
  timestamp: string
  status: boolean
}

export interface NotificationUpdate {
  ids: string[]
}

export type rateLimitInputtype = {
  rate_limit: number
  window_size?: number
  block_time?: number
}

export interface RequestLog {
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
  ai_result: boolean
  rule_detected: boolean
  ai_threat_type: string
}

export interface NotificationRule {
  id: string
  name: string
  threat_type: string
  threshold: number
  time_window: number
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface NotificationRuleInput {
  threshold: number
  time_window: number
  is_active: boolean
}

export interface UpdateNotificationRuleInput {
  rule_id: string
  data: NotificationRuleInput
}
