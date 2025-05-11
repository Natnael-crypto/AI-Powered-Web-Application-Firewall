import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'
import {
  assignApplication,
  createApplication,
  deleteAssignment,
  getApplication,
  getApplications,
  getAssignments,
  updateApplication,
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
export function useAssignApplication() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['assignApplication'],
    mutationFn: assignApplication,
    onSuccess: () => queryClient.invalidateQueries({queryKey: ['GetappAssignments']}),
  })
}
export function useDeleteAssignment() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['deleteAssignment'],
    mutationFn: deleteAssignment,
    onSuccess: () => queryClient.invalidateQueries({queryKey: ['GetappAssignments']}),
  })
}

export function useGetApplicationAssignments() {
  return useQuery({
    queryKey: ['GetappAssignments'],
    queryFn: getAssignments,
  })
}
