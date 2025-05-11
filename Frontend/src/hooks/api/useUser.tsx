import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'
import {
  activateUser,
  addUser,
  deActivateUser,
  deleteUser,
  getuser,
  getUsers,
  isLoggedIn,
  loginUser,
} from '../../services/userApi'

export const useLogin = () => {
  return useMutation({
    mutationKey: ['login'],
    mutationFn: loginUser,
  })
}

export const useGetUsers = () => {
  return useQuery({
    queryKey: ['getUsers'],
    queryFn: getUsers,
  })
}

export const useGetUser = (username: string) => {
  const {} = useQuery({
    queryKey: ['getUser'],
    queryFn: () => getuser(username),
  })
}

export const useIsLoggedIn = () => {
  return useQuery({
    queryKey: ['getMe'],
    queryFn: isLoggedIn,
    retry: 4,
  })
}

export function useAddAdmin() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['updateApplication'],
    mutationFn: addUser,
    onSuccess: () => {
      queryClient.invalidateQueries({queryKey: ['getUsers']})
    },
  })
}

export function useDeleteUser() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['deleteUser'],
    mutationFn: deleteUser,
    onSuccess: () => queryClient.invalidateQueries({queryKey: ['getUsers']}),
  })
}

export function useDeactivateUser() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['deactivateUser'],
    mutationFn: deActivateUser,
    onSuccess: () => queryClient.invalidateQueries({queryKey: ['getUsers']}),
  })
}
export function useActivateUser() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['deactivateUser'],
    mutationFn: activateUser,
    onSuccess: () => queryClient.invalidateQueries({queryKey: ['getUsers']}),
  })
}
export function useUpdateUser() {}
