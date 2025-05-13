import { useState } from 'react'

import AIModelTable from '../components/AIModelTable'
import CreateModelModal from '../components/CreateModelModal'
import { useGetAIModels } from '../hooks/api/useAIModels'

function AIAnalysis() {
  const [isCreateModel, setCreateModel] = useState(false)
  const toggleCreateModel = () => setCreateModel((prev) => !prev)

  const { data: aiModels, isLoading } = useGetAIModels()

  return (
    <div className="flex flex-col gap-8 px-6 py-10 bg-gradient-to-br from-slate-100 to-white min-h-screen">
      <div className="max-w-7xl w-full mx-auto space-y-8">
        <section className="bg-white border border-slate-200 xl shadow-lg p-8">
          <h2 className="text-xl font-semibold text-slate-800 mb-4">AI Models</h2>
          {isLoading ? (
            <div>Loading...</div>
          ) : (
            <AIModelTable aiModels={aiModels} />
          )}
          <button
            className="mt-4 bg-blue-500 text-white px-4 py-2 rounded"
            onClick={toggleCreateModel}
          >
            Create New Model
          </button>
        </section>
      </div>

      <CreateModelModal
        isOpen={isCreateModel}
        onClose={toggleCreateModel}
      />
    </div>
  )
}

export default AIAnalysis
