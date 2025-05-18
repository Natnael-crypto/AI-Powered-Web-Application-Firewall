import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'
import {
  assignApplication,
  createApplication,
  deleteAssignment,
  getApplication,
  getApplicationConfig,
  getApplications,
  getAssignments,
  getconfig,
  updateApplication,
  updateDetectBOT,
  updateListeningPort,
  updateMaxDataSize,
  updateRateLimit,
  updateRemoteLogServer,
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

export function useUpdateListeningPort() {
  return useMutation({
    mutationKey: ['updateConfig'],
    mutationFn: updateListeningPort,
  })
}
export function useUpdateRateLimit() {
  return useMutation({
    mutationKey: ['updateConfig'],
    mutationFn: updateRateLimit,
  })
}
export function useUpdateRemoteLogServer() {
  return useMutation({
    mutationKey: ['updateConfig'],
    mutationFn: updateRemoteLogServer,
  })
}

export function useGetApplicationConfig(application_id: string) {
  return useQuery({
    queryKey: ['getApplicationConfig'],
    queryFn: () => getApplicationConfig(application_id),
  })
}

export function useGetConfig() {
  return useQuery({
    queryKey: ['getApplicationConfig'],
    queryFn: getconfig,
  })
}

export function useUpdateDetectBot() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['detect_bot'],
    mutationFn: updateDetectBOT,
    onSuccess: () => queryClient.invalidateQueries({queryKey: ['applications']}),
  })
}
export function useUpdateMaxDataSize() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['detect_bot'],
    mutationFn: updateMaxDataSize,
    onSuccess: () => queryClient.invalidateQueries({queryKey: ['applications']}),
  })
}
