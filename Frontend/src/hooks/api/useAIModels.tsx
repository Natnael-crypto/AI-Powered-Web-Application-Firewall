// hooks/useAIModelHooks.ts
import { useMutation } from '@tanstack/react-query'
import { getAIModels, selectAIModel, updateAIModelSetting, updateAIModelTrainTime } from '../../services/aiModelsApi'
import { useQuery } from '@tanstack/react-query'


export function useGetAIModels() {
  return useQuery({
    queryKey: ['aiModels'],
    queryFn: getAIModels,
  })
}

export function useCreateModel() {
  return useMutation({
    mutationFn: updateAIModelSetting,
  })
}

export function useUpdateModelTrainTime() {
  return useMutation({
    mutationFn: updateAIModelTrainTime,
  })
}

export function useSelectModel() {
  return useMutation({
    mutationFn: selectAIModel,
  })
}
