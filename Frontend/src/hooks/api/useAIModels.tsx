// hooks/useAIModelHooks.ts
import { useMutation } from '@tanstack/react-query'
import { getAIModels, createAIModel, deleteAIModel, selectAIModel } from '../../services/aiModelsApi'
import { useQuery } from '@tanstack/react-query'


export function useGetAIModels() {
  return useQuery({
    queryKey: ['aiModels'],
    queryFn: getAIModels,
  })
}

export function useCreateModel() {
  return useMutation({
    mutationFn: createAIModel,
  })
}

export function useDeleteModel() {
  return useMutation({
    mutationFn: deleteAIModel,
  })
}

export function useSelectModel() {
  return useMutation({
    mutationFn: selectAIModel,
  })
}
