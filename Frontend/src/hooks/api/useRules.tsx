import {useMutation, useQuery} from '@tanstack/react-query'
import {activateRule, deactivateRule, deleteRule, getRules} from '../../services/rulesApi'

export function useGetRules() {
  return useQuery({
    queryKey: ['applications'],
    queryFn: getRules,
  })
}
export function useActivateRule(ruleId: string) {
  return useQuery({
    queryKey: ['applications', ruleId],
    queryFn: () => activateRule(ruleId),
  })
}
export function useDeactivateRule(ruleId: string) {
  return useQuery({
    queryKey: ['applications', ruleId],
    queryFn: () => deactivateRule(ruleId),
  })
}

export function useUpdateApplication() {
  return useMutation({
    mutationKey: ['updateApplication'],
    mutationFn: deleteRule,
  })
}
