import { useMutation, useQuery } from '@tanstack/react-query'
import { createSysEmail, getSysEmail, updateSysEmail } from '../../services/configApi'

export function useGetSysEmail() {
  return useQuery({
    queryKey: ['getEmail'],
    queryFn: getSysEmail,
  })
}

export function useAddSysEmail() {
  return useMutation({
    mutationFn: ({ email, active }: { email: string; active: boolean }) =>
      createSysEmail(email, active),
  })
}

export function useUpdateSysEmail() {
  return useMutation({
    mutationFn: ({ email, active }: { email: string; active: boolean }) =>
      updateSysEmail(email, active),
  })
}
