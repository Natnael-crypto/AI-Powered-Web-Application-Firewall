// hooks/useAIModelHooks.ts
import { useMutation } from '@tanstack/react-query'
import { useQuery } from '@tanstack/react-query'
import { createSecurityHeader, deleteSecurityHeader, getSecurityHeaders, updateSecurityHeader } from '../../services/securityHeaders'


export function useGetSecurityHeaders() {
  return useQuery({
    queryKey: ['securityHeaders'],
    queryFn: getSecurityHeaders,
  })
}

export function useCreateSecurityHeader() {
  return useMutation({
    mutationFn: createSecurityHeader,
  })
}

export function useDeleteSecurityHeader() {
  return useMutation({
    mutationFn: deleteSecurityHeader,
  })
}

export function useUpdateSecurityHeader() {
  return useMutation({
    mutationFn: ({ headerId, data }: { headerId: string; data: any }) => 
      updateSecurityHeader(headerId, data),
  })
}

