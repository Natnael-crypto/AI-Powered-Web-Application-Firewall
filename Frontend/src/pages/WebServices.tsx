import {useState} from 'react'
import Card from '../components/Card'
import Button from '../components/atoms/Button'
import WebServiceModal from '../components/WebServiceModal'
import WebserviceTable from '../components/WebserviceTable'
import {useUpdateApplication} from '../hooks/api/useApplication'

export interface Application {
  application_id: string
  application_name: string
  description: string
  hostname: string
  ip_address: string
  port: string
  status: boolean
  tls: boolean
  created_at: string
  updated_at: string
}

function WebService() {
  const [isModalOpen, setIsModalOpen] = useState(false)
  const toggleModal = () => setIsModalOpen(!isModalOpen)
  const [selectedApp, setSelectedApp] = useState<Application>()
  const {mutate} = useUpdateApplication()

  const handleUpdateApplication = async (
    e: {preventDefault: () => void},
    formData: any,
  ) => {
    e.preventDefault()
    try {
      mutate(formData, {
        onSuccess: data => {
          toggleModal()
        },
        onError: () => {
          console.log('Invalid username or password')
        },
      })
    } catch (error) {
      console.log('Invalid username or password')
    }
  }

  return (
    <div className="space-y-4">
      <WebServiceModal
        application={selectedApp}
        isOpen={isModalOpen}
        onClose={toggleModal}
        onSubmit={data => {
          handleUpdateApplication({preventDefault: () => {}}, data)
        }}
      />

      <Card className="flex justify-between items-center py-4 px-6 bg-white">
        <h2 className="text-lg font-semibold">Web Services</h2>
        <Button
          classname="text-white uppercase"
          size="l"
          variant="primary"
          onClick={toggleModal}
        >
          Add Service
        </Button>
      </Card>
      <Card className="shadow-md p-4 bg-white ">
        <WebserviceTable
          data={[]}
          openModal={toggleModal}
          setSelectedApp={setSelectedApp}
        />
      </Card>
    </div>
  )
}

export default WebService
