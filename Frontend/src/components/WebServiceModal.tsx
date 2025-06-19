import React, {useEffect, useState} from 'react'
import Modal from './Modal'
import Button from './atoms/Button'
import LoadingSpinner from './LoadingSpinner'
import {validateIPAddressOrDomain, validateHostname, validatePort} from '../lib/utils'
import {useGetCertification} from '../hooks/api/useApplication'

export interface WebServiceData {
  application_id?: string
  application_name: string
  description: string
  hostname: string
  ip_address: string
  port: string
  status: boolean
  tls: boolean
}

interface Certificate {
  id: string
  domain: string
  expiry_date: string
}

interface WebServiceModalProps {
  application?: WebServiceData
  isOpen: boolean
  onClose: () => void
  onSubmit: (data: WebServiceData) => void
  isSubmitting?: boolean
  onCertificateUpload?: (applicationId: string) => void
}

const defaultForm: WebServiceData = {
  application_name: '',
  description: '',
  hostname: '',
  ip_address: '',
  port: '',
  status: true,
  tls: false,
}

const WebServiceModal: React.FC<WebServiceModalProps> = ({
  application,
  isOpen,
  onClose,
  onSubmit,
  isSubmitting = false,
  onCertificateUpload,
}) => {
  const [form, setForm] = useState<WebServiceData>(defaultForm)
  const [errors, setErrors] = useState<Record<string, string>>({})
  const [touched, setTouched] = useState<Record<string, boolean>>({})

  const {data: certifications, isLoading: isLoadingCerts} = useGetCertification(
    application?.application_id ?? '',
    {
      enabled: !!application?.application_id && isOpen,
    },
  )

  useEffect(() => {
    if (isOpen) {
      setForm(application ? {...application} : defaultForm)
      setErrors({})
      setTouched({})
    }
  }, [application, isOpen])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const {name, value, type} = e.target
    const checked =
      type === 'checkbox' && 'checked' in e.target
        ? (e.target as HTMLInputElement).checked
        : undefined
    setForm(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value,
    }))

    setTouched(prev => ({...prev, [name]: true}))

    if (errors[name]) {
      setErrors(prev => {
        const updated = {...prev}
        delete updated[name]
        return updated
      })
    }
  }

  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {}

    if (!form.application_name.trim()) {
      newErrors.application_name = 'Application name is required'
    }
    if (!form.hostname.trim()) {
      newErrors.hostname = 'Hostname is required'
    } else if (!validateHostname(form.hostname)) {
      newErrors.hostname = 'Invalid hostname format'
    }

    if (!form.ip_address.trim()) {
      newErrors.ip_address = 'IP address is required'
    } else if (!validateIPAddressOrDomain(form.ip_address)) {
      newErrors.ip_address = 'Invalid IP address format'
    }

    if (!form.port.trim()) {
      newErrors.port = 'Port is required'
    } else if (!validatePort(form.port)) {
      newErrors.port = 'Port must be between 1 and 65535'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!validateForm()) return

    const payload = application?.application_id
      ? {...form, application_id: application.application_id}
      : form

    onSubmit(payload)
  }

  const getFieldError = (field: string) =>
    touched[field] && errors[field] ? errors[field] : null

  if (!isOpen) return null

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={application ? 'Edit Web Service' : 'Add Web Service'}
    >
      <form onSubmit={handleSubmit}>
        <div className="space-y-6 px-4 py-4">
          {/* Fields Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {[
              {
                label: 'Application Name',
                name: 'application_name',
                placeholder: 'Enter application name',
                required: true,
              },
              {
                label: 'Hostname',
                name: 'hostname',
                placeholder: 'example.com',
                required: true,
              },
              {
                label: 'IP Address or Domain',
                name: 'ip_address',
                placeholder: '192.168.1.1 or www.aait.com',
                required: true,
              },
              {label: 'Port', name: 'port', placeholder: '8080', required: true},
            ].map(({label, name, placeholder, required}) => (
              <div key={name} className="space-y-1">
                <label className="block text-sm font-medium text-gray-700">
                  {label} {required && <span className="text-red-500">*</span>}
                </label>
                <input
                  type="text"
                  name={name}
                  value={(form as any)[name]}
                  onChange={handleChange}
                  onBlur={() => setTouched(prev => ({...prev, [name]: true}))}
                  className={`w-full px-3 py-2 border ${
                    getFieldError(name) ? 'border-red-500' : 'border-gray-300'
                  } rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 transition`}
                  placeholder={placeholder}
                />
                {getFieldError(name) && (
                  <p className="mt-1 text-sm text-red-500">{errors[name]}</p>
                )}
              </div>
            ))}
          </div>

          {/* Description */}
          <div className="space-y-1">
            <label className="block text-sm font-medium text-gray-700">Description</label>
            <textarea
              name="description"
              value={form.description}
              onChange={handleChange}
              rows={3}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 transition resize-none"
              placeholder="Enter description"
            />
          </div>

          {/* Toggles */}
          <div className="flex flex-wrap gap-6 pt-2">
            <label className="inline-flex items-center">
              <input
                type="checkbox"
                name="status"
                checked={form.status}
                onChange={handleChange}
                className="h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              />
              <span className="ml-2 text-sm text-gray-700">Active</span>
            </label>
            <label className="inline-flex items-center">
              <input
                type="checkbox"
                name="tls"
                checked={form.tls}
                onChange={handleChange}
                className="h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              />
              <span className="ml-2 text-sm text-gray-700">TLS Enabled</span>
            </label>
          </div>

          {/* Certifications */}
          {application?.application_id && onCertificateUpload && (
            <div className="mt-6 space-y-2 border-t pt-4">
              <div className="flex justify-between items-center">
                <h3 className="text-md font-semibold text-gray-800">TLS Certificates</h3>
                <Button
                  variant="primary"
                  type="button"
                  onClick={() => onCertificateUpload(application.application_id!)}
                  classname="px-3 py-1 text-sm"
                >
                  Upload Certificate
                </Button>
              </div>

              {isLoadingCerts ? (
                <p className="text-gray-500 text-sm">Loading certificates...</p>
              ) : certifications?.length ? (
                <ul className="list-disc ml-5 text-sm text-gray-700">
                  {certifications.map((cert: Certificate) => (
                    <li key={cert.id}>
                      <span className="font-medium">{cert.domain}</span> — expires on{' '}
                      {cert.expiry_date}
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="text-gray-500 text-sm">No certificates uploaded.</p>
              )}
            </div>
          )}

          {/* Actions */}
          <div className="flex justify-end gap-3 pt-6 mt-6 border-t border-gray-200">
            <Button
              variant="secondary"
              onClick={onClose}
              classname="px-4 py-2 text-white text-sm"
              disabled={isSubmitting}
            >
              Cancel
            </Button>
            <Button
              variant="primary"
              type="submit"
              classname="px-4 py-2 text-white text-sm flex items-center justify-center min-w-24"
              disabled={isSubmitting}
            >
              {isSubmitting ? (
                <LoadingSpinner />
              ) : application ? (
                'Update Service'
              ) : (
                'Create Service'
              )}
            </Button>
          </div>
        </div>
      </form>
    </Modal>
  )
}

export default WebServiceModal
