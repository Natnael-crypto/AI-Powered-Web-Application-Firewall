import React, {useState} from 'react'
import Modal from './Modal'
import Button from './atoms/Button'
import LoadingSpinner from './LoadingSpinner'

interface CertificateUploadModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (certificate: File, key: File) => Promise<void>
  isSubmitting?: boolean
}

const CertificateUploadModal: React.FC<CertificateUploadModalProps> = ({
  isOpen,
  onClose,
  onSubmit,
  isSubmitting = false,
}) => {
  const [certificate, setCertificate] = useState<File | null>(null)
  const [key, setKey] = useState<File | null>(null)
  const [errors, setErrors] = useState<{certificate?: string; key?: string}>({})

  const handleFileChange = (
    e: React.ChangeEvent<HTMLInputElement>,
    type: 'certificate' | 'key',
  ) => {
    const file = e.target.files?.[0]
    if (file) {
      if (type === 'certificate') {
        setCertificate(file)
      } else {
        setKey(file)
      }
      setErrors(prev => ({...prev, [type]: undefined}))
    }
  }

  const validateFiles = (): boolean => {
    const newErrors: {certificate?: string; key?: string} = {}

    if (!certificate) {
      newErrors.certificate = 'Certificate file is required'
    }

    if (!key) {
      newErrors.key = 'Key file is required'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validateFiles() || !certificate || !key) {
      return
    }

    try {
      await onSubmit(certificate, key)
      setCertificate(null)
      setKey(null)
    } catch (error) {
      console.error('Certificate upload error:', error)
    }
  }

  if (!isOpen) return null

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Upload TLS Certificates">
      <form onSubmit={handleSubmit}>
        <div className="space-y-6 px-4 py-4">
          <div className="space-y-4">
            <div className="space-y-1">
              <label className="block text-sm font-medium text-gray-700">
                Certificate File <span className="text-red-500">*</span>
              </label>
              <input
                type="file"
                accept=".pem,.crt,.cer"
                onChange={e => handleFileChange(e, 'certificate')}
                className="block w-full text-sm text-gray-500
                  file:mr-4 file:py-2 file:px-4
                  file:rounded-md file:border-0
                  file:text-sm file:font-semibold
                  file:bg-blue-50 file:text-blue-700
                  hover:file:bg-blue-100"
                required
              />
              {errors.certificate && (
                <p className="mt-1 text-sm text-red-500">{errors.certificate}</p>
              )}
            </div>

            <div className="space-y-1">
              <label className="block text-sm font-medium text-gray-700">
                Key File <span className="text-red-500">*</span>
              </label>
              <input
                type="file"
                accept=".pem,.key"
                onChange={e => handleFileChange(e, 'key')}
                className="block w-full text-sm text-gray-500
                  file:mr-4 file:py-2 file:px-4
                  file:rounded-md file:border-0
                  file:text-sm file:font-semibold
                  file:bg-blue-50 file:text-blue-700
                  hover:file:bg-blue-100"
                required
              />
              {errors.key && <p className="mt-1 text-sm text-red-500">{errors.key}</p>}
            </div>
          </div>

          <div className="flex justify-end gap-3 pt-6 mt-6 border-t border-gray-200">
            <Button
              variant="secondary"
              onClick={onClose}
              classname="px-4 text-white py-2 text-sm"
              disabled={isSubmitting}
            >
              Cancel
            </Button>
            <Button
              variant="primary"
              type="submit"
              classname="px-4 py-2 text-white text-sm flex items-center justify-center min-w-24"
              disabled={isSubmitting || !certificate || !key}
            >
              {isSubmitting ? <LoadingSpinner /> : 'Upload Certificates'}
            </Button>
          </div>
        </div>
      </form>
    </Modal>
  )
}

export default CertificateUploadModal
