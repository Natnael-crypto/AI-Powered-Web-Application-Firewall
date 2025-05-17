import { useState } from 'react'

import AIModelTable from '../components/AIModelTable'
import CreateModelModal from '../components/CreateModelModal'
import { useGetAIModels } from '../hooks/api/useAIModels'
import Card from '../components/Card'

function AIAnalysis() {
  const [isCreateModel, setCreateModel] = useState(false)
  const toggleCreateModel = () => setCreateModel((prev) => !prev)

  const { data: aiModels, isLoading } = useGetAIModels()

  return (
    <div className="flex flex-col gap-8 space-y-4 bg-gradient-to-br from-slate-100 to-white min-h-screen">
      <div className="w-full mx-auto space-y-4">
        <Card className="flex justify-between items-center py-4 px-6 bg-white">
          <h2 className="text-lg font-semibold">AI Models</h2>
        </Card>

        <section className="bg-white border border-slate-200 xl shadow-lg p-8">
          {isLoading ? (
            <div>Loading...</div>
          ) : (
            <AIModelTable aiModels={aiModels} />
          )}
        </section>
      </div>
    </div>
  )
}

export default AIAnalysis
