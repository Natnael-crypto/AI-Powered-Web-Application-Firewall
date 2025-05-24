import { UpdateNotificationRuleInput } from '../../lib/types'
import { getNotificationRules, updateNotificationRule } from '../../services/notificationRuleService'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

export function useNotificationRules() {
  return useQuery({
    queryKey: ['notificationRule'],
    queryFn: () => getNotificationRules(),
  })
}



export function useUpdateNotificationRule() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['updateNotificationRule'],
    mutationFn: (data:UpdateNotificationRuleInput) => updateNotificationRule(data.rule_id,data.data),
    onSuccess: () => queryClient.invalidateQueries({queryKey: ['notificationRule']}),

  })
}