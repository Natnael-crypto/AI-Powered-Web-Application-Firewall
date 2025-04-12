import {useMutation, useQuery} from '@tanstack/react-query'
import {
  createApplication,
  getApplication,
  getApplications,
} from '../services/applicationApi'

export function useGetAppliactions() {
  return useQuery({
    queryKey: ['applications'],
    queryFn: getApplications,
  })
}
export function useGetAppliaction(application_id: string) {
  return useQuery({
    queryKey: ['application'],
    queryFn: () => getApplication(application_id),
  })
}

export function useAddApplication() {
  return useMutation({
    mutationKey: ['addApplication'],
    mutationFn: createApplication,
  })
}
