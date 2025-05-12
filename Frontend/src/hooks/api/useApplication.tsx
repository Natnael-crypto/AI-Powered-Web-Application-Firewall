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
export function useAssignApplication(p0: {onSuccess: () => void}) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['assignApplication'],
    mutationFn: assignApplication,
    onSuccess: () => {
      queryClient.invalidateQueries({queryKey: ['GetappAssignments']}), p0
    },
  })
}
export function useDeleteAssignment(p0: {onSuccess: () => void}) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['deleteAssignment'],
    mutationFn: deleteAssignment,
    onSuccess: () => {
      queryClient.invalidateQueries({queryKey: ['GetappAssignments']}), p0
    },
  })
}

export function useGetApplicationAssignments() {
  return useQuery({
    queryKey: ['GetappAssignments'],
    queryFn: getAssignments,
  })
}
