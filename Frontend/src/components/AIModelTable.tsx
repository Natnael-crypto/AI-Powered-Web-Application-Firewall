
interface AIModel {
  id: string
  models_name: string
  number_requests_used: number
  percent_train_data: number
  percent_normal_requests: number
  num_trees: number
  max_depth: number
  min_samples_split: number
  min_samples_leaf: number
  max_features: string
  criterion: string
  accuracy: number
  precision: number
  recall: number
  f1: number
  selected: boolean
  modeled: boolean
  created_at: string
  updated_at: string
}

interface AIModelTableProps {
  aiModels: AIModel[] | undefined
}

const AIModelTable = ({ aiModels = [] }: AIModelTableProps) => {
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
            <th className="px-6 py-3 text-left text-sm font-semibold text-gray-700">Created At</th>
          </tr>
        </thead>
        <tbody>
          {aiModels.length === 0 ? (
            <tr>
              <td colSpan={7} className="px-6 py-4 text-sm text-gray-900 text-center">No models available</td>
            </tr>
          ) : (
            aiModels.map((model) => (
              <tr key={model.id} className="border-b">
                <td className="px-6 py-4 text-sm text-gray-900">{model.models_name}</td>
                <td className="px-6 py-4 text-sm text-gray-900">{model.accuracy}</td>
                <td className="px-6 py-4 text-sm text-gray-900">{model.precision}</td>
                <td className="px-6 py-4 text-sm text-gray-900">{model.recall}</td>
                <td className="px-6 py-4 text-sm text-gray-900">{model.f1}</td>
                <td className="px-6 py-4 text-sm text-gray-900">
                  {model.selected ? 'Yes' : 'No'}
                </td>
                <td className="px-6 py-4 text-sm text-gray-900">
                  {model.modeled ? 'Yes' : 'No'}
                </td>
                <td className="px-6 py-4 text-sm text-gray-900">{model.created_at}</td>
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  )
}

export default AIModelTable
