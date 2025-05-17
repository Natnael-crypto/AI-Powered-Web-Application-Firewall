
import { useDeleteModel, useSelectModel } from '../hooks/api/useAIModels'
import { PencilIcon, TrashIcon } from 'lucide-react'



interface AIModel {
  id: string
  models_name: string
  number_requests_used: number
  train_every:number,
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
  percent_of_training_data: 0
}

interface AIModelTableProps {
  aiModels: AIModel[] | undefined
}

const AIModelTable = ({ aiModels = [] }: AIModelTableProps) => {
  const deleteMutation = useDeleteModel()
  const selectMutation = useSelectModel()

  const handleSelect = (id: string) => {
    selectMutation.mutate(id)
  }

  const handleDelete = (id: string) => {
    if (confirm('Are you sure you want to delete this model?')) {
      deleteMutation.mutate(id)
    }
  }

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full table-auto">
        <thead className="border-b bg-gray-50">
          <tr>
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Model Name</th>
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Accuracy</th>
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Precision</th>
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Recall</th>
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">F1 Score</th>
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Selected</th>
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Modeled</th>
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Train Every</th>
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Actions</th>
          </tr>
        </thead>
        <tbody>
          {aiModels.length === 0 ? (
            <tr>
              <td colSpan={9} className="px-6 py-4 text-sm text-gray-900 text-center">No models available</td>
            </tr>
          ) : (
            aiModels.map((model) => (
              <tr key={model.id} className="border-b">
                <td className="px-6 py-4 text-sm text-gray-900">{model.models_name}</td>
                <td className="px-6 py-4 text-sm text-gray-900">{model.accuracy}</td>
                <td className="px-6 py-4 text-sm text-gray-900">{model.precision}</td>
                <td className="px-6 py-4 text-sm text-gray-900">{model.recall}</td>
                <td className="px-6 py-4 text-sm text-gray-900">{model.f1}</td>
                <td className="px-6 py-4 text-sm text-gray-900">{model.selected ? 'Yes' : 'No'}</td>
                <td className="px-6 py-4 text-sm text-gray-900">{model.modeled ? 'Yes' : 'No'}</td>
                <td className="px-6 py-4 text-sm text-gray-900">{(model.train_every/3600000)}h</td>
                <td className="px-6 py-4 text-sm text-gray-900">
                  <button
                  onClick={() => handleSelect(model.id)}
                  className="text-blue-600 hover:underline px-3"
                >
                  <PencilIcon size={16} />
                </button>
                <button
                  onClick={() => handleDelete(model.id)}
                  className="text-red-600 hover:underline"
                >
                  <TrashIcon size={16} />
                </button>
                </td>
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  )
}

export default AIModelTable
