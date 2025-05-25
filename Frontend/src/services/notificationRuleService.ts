import axios from "../lib/axios";
import { NotificationRule } from '../lib/types'

const API_URL = '/api/notification-rule'

export async function getNotificationRules(): Promise<NotificationRule[]> {
  const res = await axios.get(API_URL)
  return res.data.notification_rules
}

export async function updateNotificationRule(ruleId: string, updates: Partial<NotificationRule>) {
  const res = await axios.put(`${API_URL}/${ruleId}`, updates)
  return res.data
}
