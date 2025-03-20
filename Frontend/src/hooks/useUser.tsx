import {useMutation, useQuery} from '@tanstack/react-query'
import {getuser, getUsers, loginUser} from '../services/userApi'

export const useLogin = () => {
  return useMutation({
    mutationKey: ['login'],
    mutationFn: loginUser,
  })
}

export const useGetUsers = async () => {
  const {data, isLoading, error} = useQuery({
    queryKey: ['getUsers'],
    queryFn: getUsers,
  })

  return {data, error, isLoading}
}

export const useGetUser = async (username: string) => {
  const {} = useQuery({
    queryKey: ['getUser'],
    queryFn: () => getuser(username),
  })
}
