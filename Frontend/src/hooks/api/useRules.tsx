import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'
import {
  activateRule,
  createRule,
  deactivateRule,
  deleteRule,
  getRules,
  updateRule,
} from '../../services/rulesApi'

export function useGetRules() {
  return useQuery({
    queryKey: ['getRules'],
    queryFn: getRules,
  })
}
export function useActivateRule() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['activateRules'],
    mutationFn: activateRule,
    onSuccess: () => queryClient.invalidateQueries({queryKey: ['getRules']}),
  })
}
export function useDeactivateRule() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['deactivateRules'],
    mutationFn: deactivateRule,
    onSuccess: () => queryClient.invalidateQueries({queryKey: ['getRules']}),
  })
}

export function useDeleteRules() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['deleteRules'],
    mutationFn: deleteRule,
    onSuccess: () => queryClient.invalidateQueries({queryKey: ['getRules']}),
  })
}

export function useCreateRule() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['createRule'],
    mutationFn: createRule,
    onSuccess: () => queryClient.invalidateQueries({queryKey: ['getRules']}),
  })
}
export function useUpdateRule() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['updateRule'],
    mutationFn: updateRule,
    onSuccess: () => queryClient.invalidateQueries({queryKey: ['getRules']}),
  })
}
