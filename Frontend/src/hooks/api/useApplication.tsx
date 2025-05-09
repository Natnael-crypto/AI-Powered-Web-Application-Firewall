import {useMutation, useQuery} from '@tanstack/react-query'
import {
  createApplication,
  getApplication,
  getApplications,
  updateApplication, // renamed
} from '../../services/applicationApi'

export function useGetApplications() {
  return useQuery({
    queryKey: ['applications'],
    queryFn: getApplications,
  })
}

export function useGetApplication(application_id: string) {
  return useQuery({
    queryKey: ['application', application_id],
    queryFn: () => getApplication(application_id),
  })
}

export function useAddApplication() {
  return useMutation({
    mutationKey: ['addApplication'],
    mutationFn: createApplication,
  })
}

export function useUpdateApplication() {
  return useMutation({
    mutationKey: ['updateApplication'],
    mutationFn: updateApplication,
  })
}
