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

  const {data: applications = []} = useGetApplications()
  const {mutate: createApplication} = useAddApplication()
  const {mutate: updateApplication} = useUpdateApplication()

  const toggleModal = () => setIsModalOpen(!isModalOpen)

  const handleOpenCreateModal = () => {
    setSelectedApp(undefined)
    setIsModalOpen(true)
  }

  const handleFormSubmit = (formData: Partial<Application>) => {
    const isUpdate = !!formData.application_id

    const mutationFn = isUpdate ? updateApplication : createApplication

    mutationFn(formData, {
      onSuccess: () => {
        toggleModal()
      },
      onError: () => {
        console.error('Something went wrong while saving the application.')
      },
    })
  }

  return (
    <div className="space-y-4">
      <WebServiceModal
        application={selectedApp}
        isOpen={isModalOpen}
        onClose={toggleModal}
        onSubmit={handleFormSubmit}
      />

      <Card className="flex justify-between items-center py-4 px-6 bg-white">
        <h2 className="text-lg font-semibold">Web Services</h2>
        <Button
          classname="text-white uppercase"
          size="l"
          variant="primary"
          onClick={handleOpenCreateModal}
        >
          Add Service
        </Button>
      </Card>

      <Card className="shadow-md p-4 bg-white ">
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
