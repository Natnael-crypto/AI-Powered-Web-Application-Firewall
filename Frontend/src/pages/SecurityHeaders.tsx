import { useState } from 'react'

import SecurityHeaderTable from '../components/SecurityHeaderTable'
import CreateSecurityHeaderModal from '../components/CreateSecurityHeaderModal'
import { useGetSecurityHeaders } from '../hooks/api/useSecurityHeaders'
import Card from '../components/Card'
import LoadingSpinner from '../components/LoadingSpinner'

function SecurityHeaders() {
  const [isCreateOpen, setCreateOpen] = useState(false)
  const toggleCreateModal = () => setCreateOpen((prev) => !prev)

  const { data: securityHeaders, isLoading } = useGetSecurityHeaders()

    if (isLoading) return <LoadingSpinner />

  return (
    <div className="flex flex-col gap-8 space-y-4 bg-gradient-to-br from-slate-100 to-white min-h-screen">
      <div className="w-full mx-auto space-y-4">
        <Card className="flex justify-between items-center py-4 px-6 bg-white">
          <h2 className="text-lg font-semibold">Security Headers</h2>
          <button
            className="mt-4 bg-blue-500 text-white px-4 py-2 rounded"
            onClick={toggleCreateModal}
            style={{backgroundColor: '#1F263E'}}
          >
            Create New Header
          </button>
        </Card>

        <section className="bg-white border border-slate-200 xl shadow-lg p-8">
          {isLoading ? (
            <div>Loading...</div>
          ) : (
            <SecurityHeaderTable securityHeaders={securityHeaders} />
          )}
        </section>
      </div>

      <CreateSecurityHeaderModal
        isOpen={isCreateOpen}
        onClose={toggleCreateModal}
      />
    </div>
  )
}

export default SecurityHeaders
