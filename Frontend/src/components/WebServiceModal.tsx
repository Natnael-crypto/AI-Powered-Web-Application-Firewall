import React, {useEffect, useState} from 'react'
import Modal from './Modal'
import Button from './atoms/Button'
import {validateIPAddress, validateHostname, validatePort} from '../lib/utils'
import LoadingSpinner from './LoadingSpinner'

export interface WebServiceData {
  application_id?: string
  application_name: string
  description: string
  hostname: string
  ip_address: string
  port: string
  status: boolean
  tls: boolean
  // Add certificate fields if needed
  // certificate?: File
  // key?: File
}

interface WebServiceModalProps {
  application?: WebServiceData
  isOpen: boolean
  onClose: () => void
  onSubmit: (data: WebServiceData) => void
  isSubmitting?: boolean
}

const WebServiceModal: React.FC<WebServiceModalProps> = ({
  isOpen,
  onClose,
  onSubmit,
  application,
  isSubmitting = false,
}) => {
  const [form, setForm] = useState<WebServiceData>({
    application_name: '',
    description: '',
    hostname: '',
    ip_address: '',
    port: '',
    status: true,
    tls: false,
  })

  const [errors, setErrors] = useState<Record<string, string>>({})
  const [touched, setTouched] = useState<Record<string, boolean>>({})

  useEffect(() => {
    if (application) {
      setForm({
        application_name: application.application_name,
        description: application.description,
        hostname: application.hostname,
        ip_address: application.ip_address,
        port: application.port,
        status: application.status,
        tls: application.tls,
      })
    } else {
      // Reset form when creating new application
      setForm({
        application_name: '',
        description: '',
        hostname: '',
        ip_address: '',
        port: '',
        status: true,
        tls: false,
      })
    }
    setErrors({})
    setTouched({})
  }, [application, isOpen])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const {name, value, type, checked} = e.target as HTMLInputElement
    setForm(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value,
    }))

    // Mark field as touched
    setTouched(prev => ({...prev, [name]: true}))

    // Clear error when user types
    if (errors[name]) {
      setErrors(prev => {
        const newErrors = {...prev}
        delete newErrors[name]
        return newErrors
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
    } else if (!validateIPAddress(form.ip_address)) {
      newErrors.ip_address = 'Invalid IP address format'
    }

    if (!form.port.trim()) {
      newErrors.port = 'Port is required'
    } else if (!validatePort(form.port)) {
      newErrors.port = 'Port must be between 1 and 65535'
    }

    // Add TLS certificate validation if needed
    if (form.tls) {
      // Validate certificate files here if implementing file upload
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validateForm()) {
      return
    }

    try {
      const dataToSubmit = application?.application_id
        ? {...form, application_id: application.application_id}
        : form

      onSubmit(dataToSubmit)
    } catch (error) {
      console.error('Submission error:', error)
    }
  }

  const getFieldError = (fieldName: string) => {
    return touched[fieldName] && errors[fieldName] ? errors[fieldName] : null
  }

  if (!isOpen) return null

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={application ? 'Edit Web Service' : 'Add Web Service'}
    >
      <form onSubmit={handleSubmit}>
        <div className="space-y-6 px-4 py-4">
          {/* Main Form Fields */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-1">
              <label className="block text-sm font-medium text-gray-700">
                Application Name <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                name="application_name"
                value={form.application_name}
                onChange={handleChange}
                onBlur={() => setTouched(prev => ({...prev, application_name: true}))}
                className={`w-full px-3 py-2 border ${getFieldError('application_name') ? 'border-red-500' : 'border-gray-300'} rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition`}
                placeholder="Enter application name"
                required
              />
              {getFieldError('application_name') && (
                <p className="mt-1 text-sm text-red-500">{errors.application_name}</p>
              )}
            </div>

            <div className="space-y-1">
              <label className="block text-sm font-medium text-gray-700">
                Hostname <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                name="hostname"
                value={form.hostname}
                onChange={handleChange}
                onBlur={() => setTouched(prev => ({...prev, hostname: true}))}
                className={`w-full px-3 py-2 border ${getFieldError('hostname') ? 'border-red-500' : 'border-gray-300'} rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition`}
                placeholder="example.com"
                required
              />
              {getFieldError('hostname') && (
                <p className="mt-1 text-sm text-red-500">{errors.hostname}</p>
              )}
            </div>

            <div className="space-y-1">
              <label className="block text-sm font-medium text-gray-700">
                IP Address <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                name="ip_address"
                value={form.ip_address}
                onChange={handleChange}
                onBlur={() => setTouched(prev => ({...prev, ip_address: true}))}
                className={`w-full px-3 py-2 border ${getFieldError('ip_address') ? 'border-red-500' : 'border-gray-300'} rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition`}
                placeholder="192.168.1.1"
                required
              />
              {getFieldError('ip_address') && (
                <p className="mt-1 text-sm text-red-500">{errors.ip_address}</p>
              )}
            </div>

            <div className="space-y-1">
              <label className="block text-sm font-medium text-gray-700">
                Port <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                name="port"
                value={form.port}
                onChange={handleChange}
                onBlur={() => setTouched(prev => ({...prev, port: true}))}
                className={`w-full px-3 py-2 border ${getFieldError('port') ? 'border-red-500' : 'border-gray-300'} rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition`}
                placeholder="8080"
                required
              />
              {getFieldError('port') && (
                <p className="mt-1 text-sm text-red-500">{errors.port}</p>
              )}
            </div>
          </div>

          {/* Description Field */}
          <div className="space-y-1">
            <label className="block text-sm font-medium text-gray-700">Description</label>
            <textarea
              name="description"
              value={form.description}
              onChange={handleChange}
              rows={3}
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition resize-none"
              placeholder="Enter description"
            />
          </div>

          {/* Toggle Options */}
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

          {/* TLS Certificate Fields (Conditional) */}
          {form.tls && (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6 pt-4">
              <div className="space-y-1">
                <label className="block text-sm font-medium text-gray-700">
                  Certificate File <span className="text-red-500">*</span>
                </label>
                <input
                  type="file"
                  className="block w-full text-sm text-gray-500
                    file:mr-4 file:py-2 file:px-4
                    file:rounded-md file:border-0
                    file:text-sm file:font-semibold
                    file:bg-blue-50 file:text-blue-700
                    hover:file:bg-blue-100"
                  required={form.tls}
                />
              </div>
              <div className="space-y-1">
                <label className="block text-sm font-medium text-gray-700">
                  Key File <span className="text-red-500">*</span>
                </label>
                <input
                  type="file"
                  className="block w-full text-sm text-gray-500
                    file:mr-4 file:py-2 file:px-4
                    file:rounded-md file:border-0
                    file:text-sm file:font-semibold
                    file:bg-blue-50 file:text-blue-700
                    hover:file:bg-blue-100"
                  required={form.tls}
                />
              </div>
            </div>
          )}

          {/* Form Actions */}
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
