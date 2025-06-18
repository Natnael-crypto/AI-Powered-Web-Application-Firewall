import {useState} from 'react'
import {useGetAIModels} from '../hooks/api/useAIModels'
import {PencilIcon} from 'lucide-react'
import UpdateAIModelSetting from './UpdateAIModelSetting' // Adjust path if needed
import LoadingSpinner from './LoadingSpinner'

interface AIModel {
  id: string
  models_name: string
  number_requests_used: number
  train_every: number
  accuracy: number
  precision: number
  recall: number
  f1: number
  selected: boolean
  modeled: boolean
  created_at: string
  updated_at: string
  expected_accuracy: number
  expected_precision: number
  expected_recall: number
  expected_f1: number
  percent_of_training_data: number
}

const AIModelTable = () => {
  const {data: aiModels, isLoading, refetch} = useGetAIModels()

  const [isEditModalOpen, setIsEditModalOpen] = useState(false)
  const [selectedModel, setSelectedModel] = useState<AIModel | null>(null)


  const handleEdit = (model: AIModel) => {
    setSelectedModel(model)
    setIsEditModalOpen(true)
  }

  const closeModal = () => {
    setIsEditModalOpen(false)
    setSelectedModel(null)
    refetch()
  }

  if (isLoading) return <LoadingSpinner />
  return (
    <>
      <div className="overflow-x-auto">
        <table className="min-w-full table-auto">
          <thead className="border-b bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">
                Model Name
              </th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">
                Accuracy
              </th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">
                Precision
              </th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">
                Recall
              </th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">
                F1 Score
              </th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">
                Expected Accuracy
              </th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">
                Expected Precision
              </th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">
                Expected Recall
              </th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">
                Expected F1 Score
              </th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">
                Train Every
              </th>
              <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">
                Actions
              </th>
            </tr>
          </thead>
          <tbody>
            {aiModels?.length === 0 ? (
              <tr>
                <td colSpan={9} className="px-6 py-4 text-sm text-gray-900 text-center">
                  No models available
                </td>
              </tr>
            ) : (
              aiModels?.map((model: AIModel) => (
                <tr key={model.id} className="border-b">
                  <td className="px-6 py-4 text-sm text-gray-900 ">
                    {model.models_name}
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-900 text-center">
                    {model.accuracy}
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-900 text-center">
                    {model.precision}
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-900 text-center">
                    {model.recall}
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-900 text-center">
                    {model.f1}
                  </td>
                  <td className="px-2 py-4 text-sm text-gray-900 text-center">
                    {model.expected_accuracy}
                  </td>
                  <td className="px-2 py-4 text-sm text-gray-900 text-center">
                    {model.expected_precision}
                  </td>
                  <td className="px-2 py-4 text-sm text-gray-900 text-center">
                    {model.expected_recall}
                  </td>
                  <td className="px-2 py-4 text-sm text-gray-900 text-center">
                    {model.expected_f1}
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-900 text-center">
                    {model.train_every /86400000}d
                  </td>

                  <td className="px-6 py-4 text-sm text-gray-900 text-center">
                    <button
                      onClick={() => handleEdit(model)}
                      className="text-blue-600 hover:underline px-3"
                    >
                      <PencilIcon size={16} />
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {/* Update Modal */}
      {selectedModel && (
        <UpdateAIModelSetting
          isOpen={isEditModalOpen}
          onClose={closeModal}
          aiModelSetting={selectedModel}
        />
      )}
    </>
  )
}

export default AIModelTable
