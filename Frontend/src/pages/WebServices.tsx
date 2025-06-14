import {useState} from 'react'
import Card from '../components/Card'
import Button from '../components/atoms/Button'
import WebServiceModal from '../components/WebServiceModal'
import WebserviceTable from '../components/WebserviceTable'
import {
  useAddApplication,
  useDeleteApplication,
  useGetApplications,
  useGetCertification,
  useUpdateApplication,
  useUploadCertificate,
} from '../hooks/api/useApplication'
import {toast} from 'react-toastify'
import LoadingSpinner from '../components/LoadingSpinner'
import CertificateUploadModal from '../components/CertificateModal'

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
  const [isCertModalOpen, setIsCertModalOpen] = useState(false)
  const [selectedApp, setSelectedApp] = useState<Application | undefined>()
  const [isSubmitting, setIsSubmitting] = useState(false)
  const {data: cert} = useGetCertification(selectedApp?.application_id || '', 'cert')
  const {data: key} = useGetCertification(selectedApp?.application_id || '', 'key')
  console.log('certs', cert, key)

  const {
    data: applications = [],
    isLoading: isLoadingApplications,
    error: applicationsError,
    refetch: refetchApplications,
  } = useGetApplications()

  const {mutate: createApplication} = useAddApplication()
  const {mutate: updateApplication} = useUpdateApplication()
  const {mutate: deleteApplication} = useDeleteApplication()
  const {mutate: uploadCertificate} = useUploadCertificate()

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

      mutationFn(formData, {
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

  const handleCertificateUpload = async (certificate: File | null, key: File | null) => {
    if (!selectedApp?.application_id || !certificate || !key) {
      toast.error('Please select both certificate and key files.')
      return
    }

    setIsSubmitting(true)
    try {
      uploadCertificate(
        {
          application_id: selectedApp.application_id,
          certificate,
          key,
        } as {
          application_id: string
          certificate: File
          key: File
        },
        {
          onSuccess: (): void => {
            toast.success('Certificates uploaded successfully!')
            refetchApplications()
            setIsCertModalOpen(false)
          },
          onError: (error: {message: string}): void => {
            toast.error(`Failed to upload certificates: ${error.message}`)
          },
        } as {
          onSuccess: () => void
          onError: (error: {message: string}) => void
        },
      )
    } catch (error) {
      toast.error('An unexpected error occurred')
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleDeleteApplication = (application_id: string) => {
    deleteApplication(application_id, {
      onSuccess: () => toast('Deleted successfully'),
    })
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

      <CertificateUploadModal
        isOpen={isCertModalOpen}
        onClose={() => setIsCertModalOpen(false)}
        onSubmit={handleCertificateUpload}
        isSubmitting={isSubmitting}
        existingCert={cert}
        existingKey={key}
      />

      <Card className="flex items-center justify-between py-4 px-6 bg-white">
        <h2 className="font-semibold text-lg">Web Services</h2>
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

      <Card className="bg-white shadow-md p-4">
        <WebserviceTable
          data={applications ?? []}
          openModal={() => setIsModalOpen(true)}
          setSelectedApp={setSelectedApp}
          selectedApp={selectedApp}
          handleDelete={handleDeleteApplication}
          setIsCertModalOpen={setIsCertModalOpen}
          hasCert={!!cert}
          hasKey={!!key}
        />
      </Card>
    </div>
  )
}

export default WebService
