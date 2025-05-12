import { useMutation, useQuery } from '@tanstack/react-query'
import { getAIModels, createAIModel } from '../../services/AIModelsApi'

export function useGetAIModels() {
  return useQuery({
    queryKey: ['aiModels'],
    queryFn: getAIModels,
  })
}

export function useCreateAIModel() {
  return useMutation({
    mutationFn: createAIModel,
  })
}