import React, {useState} from 'react'
import Modal from './Modal'
import Button from './atoms/Button'
import LoadingSpinner from './LoadingSpinner'

interface CertificateUploadModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (certificate: File | null, key: File | null) => Promise<void>
  isSubmitting?: boolean
  existingCert?: File
  existingKey?: File
  mode?: 'upload' | 'update'
}

const CertificateUploadModal: React.FC<CertificateUploadModalProps> = ({
  isOpen,
  onClose,
  onSubmit,
  isSubmitting = false,
  existingCert,
  existingKey,
  mode = 'upload',
}) => {
  const [certificate, setCertificate] = useState<File | null>(null)
  const [key, setKey] = useState<File | null>(null)
  const [errors, setErrors] = useState<{certificate?: string; key?: string}>({})

  const handleFileChange = (
    e: React.ChangeEvent<HTMLInputElement>,
    type: 'certificate' | 'key',
  ) => {
    const file = e.target.files?.[0] || null
    if (type === 'certificate') {
      setCertificate(file)
    } else {
      setKey(file)
    }
    setErrors(prev => ({...prev, [type]: undefined}))
  }

  const validateFiles = (): boolean => {
    const newErrors: {certificate?: string; key?: string} = {}

    if (mode === 'update' && existingCert && existingKey && !certificate && !key) {
      newErrors.certificate = 'At least one file must be updated'
      newErrors.key = 'At least one file must be updated'
    } else if (mode === 'upload') {
      if (!certificate) newErrors.certificate = 'Certificate file is required'
      if (!key) newErrors.key = 'Key file is required'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!validateFiles()) return

    const certToSubmit = certificate || (mode === 'update' ? existingCert || null : null)
    const keyToSubmit = key || (mode === 'update' ? existingKey || null : null)

    try {
      await onSubmit(certToSubmit, keyToSubmit)
      setCertificate(null)
      setKey(null)
    } catch (err) {
      console.error('Certificate upload failed:', err)
    }
  }

  if (!isOpen) return null

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={`${mode === 'upload' ? 'Upload' : 'Update'} TLS Certificates`}
    >
      <form onSubmit={handleSubmit}>
        <div className="space-y-6 px-4 py-4">
          <div>
            <label className="block text-sm font-medium text-gray-700">
              Certificate File{' '}
              {mode === 'upload' && <span className="text-red-500">*</span>}
            </label>
            {existingCert && mode === 'update' && (
              <p className="text-xs text-gray-500 italic">
                A certificate is already present
              </p>
            )}
            <input
              type="file"
              accept=".pem,.crt,.cer"
              onChange={e => handleFileChange(e, 'certificate')}
              className="mt-1 block w-full text-sm text-gray-500 file:py-2 file:px-4 file:rounded-md file:border-0 file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
            />
            {errors.certificate && (
              <p className="text-sm text-red-500 mt-1">{errors.certificate}</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700">
              Key File {mode === 'upload' && <span className="text-red-500">*</span>}
            </label>
            {existingKey && mode === 'update' && (
              <p className="text-xs text-gray-500 italic">
                A key file is already present
              </p>
            )}
            <input
              type="file"
              accept=".pem,.key"
              onChange={e => handleFileChange(e, 'key')}
              className="mt-1 block w-full text-sm text-gray-500 file:py-2 file:px-4 file:rounded-md file:border-0 file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
            />
            {errors.key && <p className="text-sm text-red-500 mt-1">{errors.key}</p>}
          </div>

          <div className="flex justify-end pt-6 border-t border-gray-200 mt-6 space-x-3">
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
              classname="px-4 py-2 text-sm text-white min-w-24 flex justify-center items-center"
              disabled={
                isSubmitting ||
                (mode === 'upload' && (!certificate || !key)) ||
                (mode === 'update' && !certificate && !key)
              }
            >
              {isSubmitting ? (
                <LoadingSpinner />
              ) : mode === 'upload' ? (
                'Upload Certificates'
              ) : (
                'Update Certificates'
              )}
            </Button>
          </div>
        </div>
      </form>
    </Modal>
  )
}

export default CertificateUploadModal
