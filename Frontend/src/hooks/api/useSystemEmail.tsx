import { useMutation, useQuery } from '@tanstack/react-query'
import { getUserEmail, createUserEmail, updateUserEmail, deleteUserEmail } from '../../services/configApi'

export function useGetUseEmail() {
  return useQuery({
    queryKey: ['getEmail'],
    queryFn: getUserEmail,
  })
}

export function useAddUserEmail() {
  return useMutation({
    mutationFn: ({ email, id }: { email: string; id: string }) =>
      createUserEmail(email, id),
  })
}

export function useUpdateUserEmail() {
  return useMutation({
    mutationFn: ({ email, id }: { email: string; id: string }) =>
      updateUserEmail(email, id),
  })
}


export function useDeleteUserEmail() {
  return useMutation({
    mutationFn: ({ id }: { id: string }) =>
      deleteUserEmail(id),
  })
}