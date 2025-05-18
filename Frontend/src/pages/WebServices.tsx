import {useState} from 'react'
import Card from '../components/Card'
import Button from '../components/atoms/Button'
import WebServiceModal from '../components/WebServiceModal'
import WebserviceTable from '../components/WebserviceTable'
import {
  useAddApplication,
  useGetApplications,
  useUpdateApplication,
} from '../hooks/api/useApplication'
import {toast} from 'react-toastify'
import LoadingSpinner from '../components/LoadingSpinner'

interface Config {
  id: string
  application_id: string
  rate_limit: number
  window_size: number
  block_time: number
  detect_bot: boolean
  hostname: string
  max_post_data_size: number
  tls: boolean
}

export interface Application {
  application_id: string
  application_name: string
  description: string
  hostname: string
  ip_address: string
  port: string
  status: boolean
  tls: boolean
  config: Config
  created_at: string
  updated_at: string
}

function WebService() {
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [selectedApp, setSelectedApp] = useState<Application | undefined>()
  const [isSubmitting, setIsSubmitting] = useState(false)

  const {
    data: applications = [],
    isLoading: isLoadingApplications,
    error: applicationsError,
    refetch: refetchApplications,
  } = useGetApplications()

  const {mutate: createApplication} = useAddApplication()
  const {mutate: updateApplication} = useUpdateApplication()

  const toggleModal = () => setIsModalOpen(!isModalOpen)

  const handleOpenCreateModal = () => {
    setSelectedApp(undefined)
    setIsModalOpen(true)
  }

  const handleFormSubmit = async (formData: Partial<Application>) => {
    setIsSubmitting(true)
    const isUpdate = !!formData.application_id

    try {
      const mutationFn = isUpdate ? updateApplication : createApplication

      await mutationFn(formData, {
        onSuccess: () => {
          toast.success(`Application ${isUpdate ? 'updated' : 'created'} successfully!`)
          refetchApplications()
          toggleModal()
        },
        onError: error => {
          toast.error(
            `Failed to ${isUpdate ? 'update' : 'create'} application: ${error.message}`,
          )
        },
        onSettled: () => {
          setIsSubmitting(false)
        },
      })
    } catch (error) {
      toast.error('An unexpected error occurred')
      setIsSubmitting(false)
    }
  }

  if (applicationsError) {
    return (
      <div className="p-4 text-red-500">
        Error loading applications: {applicationsError.message}
        <Button onClick={() => refetchApplications()} classname="ml-4">
          Retry
        </Button>
      </div>
    )
  }

  if (isLoadingApplications) {
    return <LoadingSpinner />
  }

  return (
    <div className="space-y-4">
      <WebServiceModal
        application={selectedApp}
        isOpen={isModalOpen}
        onClose={toggleModal}
        onSubmit={handleFormSubmit}
        isSubmitting={isSubmitting}
      />

      <Card className="flex justify-between items-center py-4 px-6 bg-white">
        <h2 className="text-lg font-semibold">Web Services</h2>
        <Button
          classname="text-white uppercase"
          size="l"
          variant="primary"
          onClick={handleOpenCreateModal}
          disabled={isLoadingApplications}
        >
          Add Service
        </Button>
      </Card>

      <Card className="shadow-md p-4 bg-white">
        <WebserviceTable
          data={applications}
          openModal={() => setIsModalOpen(true)}
          setSelectedApp={setSelectedApp}
          selectedApp={selectedApp}
        />
      </Card>
    </div>
  )
}

export default WebService
