// hooks/useAIModelHooks.ts
import { useMutation } from '@tanstack/react-query'
import { getAIModels, selectAIModel, updateAIModelSetting } from '../../services/aiModelsApi'
import { useQuery } from '@tanstack/react-query'


export function useGetAIModels() {
  return useQuery({
    queryKey: ['aiModels'],
    queryFn: getAIModels,
  })
}

export function useUpdateAiModelSetting() {
  return useMutation({
    mutationFn: updateAIModelSetting,
  })
}

export function useSelectModel() {
  return useMutation({
    mutationFn: selectAIModel,
  })
}
