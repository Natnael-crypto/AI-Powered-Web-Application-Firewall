import React, {useState} from 'react'
import Modal from './Modal'
import Button from './atoms/Button'

export interface WebServiceData {
  application_name: string
  description: string
  hostname: string
  ip_address: string
  port: string
  status: boolean
  tls: boolean
}

interface WebServiceModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (data: WebServiceData) => void
}

const WebServiceModal: React.FC<WebServiceModalProps> = ({isOpen, onClose, onSubmit}) => {
  const [form, setForm] = useState<WebServiceData>({
    application_name: '',
    description: '',
    hostname: '',
    ip_address: '',
    port: '',
    status: true,
    tls: false,
  })

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const {name, value, type, checked} = e.target as HTMLInputElement
    setForm(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value,
    }))
  }

  const handleSubmit = () => {
    onSubmit(form)
    onClose()
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Add Web Service">
      <div className="space-y-6">
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Application Name
            </label>
            <input
              type="text"
              name="application_name"
              value={form.application_name}
              onChange={handleChange}
              className="w-full  border border-gray-300 px-4 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
              placeholder="Enter application name"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Description
            </label>
            <textarea
              name="description"
              value={form.description}
              onChange={handleChange}
              rows={3}
              className="w-full  border border-gray-300 px-4 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none resize-none"
              placeholder="Enter description"
            />
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Hostname
              </label>
              <input
                type="text"
                name="hostname"
                value={form.hostname}
                onChange={handleChange}
                className="w-full  border border-gray-300 px-4 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
                placeholder="example.com"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                IP Address
              </label>
              <input
                type="text"
                name="ip_address"
                value={form.ip_address}
                onChange={handleChange}
                className="w-full  border border-gray-300 px-4 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
                placeholder="192.168.1.1"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Port</label>
            <input
              type="text"
              name="port"
              value={form.port}
              onChange={handleChange}
              className="w-full  border border-gray-300 px-4 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none"
              placeholder="8080"
            />
          </div>

          <div className="flex flex-col sm:flex-row gap-6">
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                name="status"
                checked={form.status}
                onChange={handleChange}
                className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              />
              <span className="text-sm font-medium text-gray-900">Active</span>
            </label>

            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                name="tls"
                checked={form.tls}
                onChange={handleChange}
                className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              />
              <span className="text-sm font-medium text-gray-900">TLS Enabled</span>
            </label>
          </div>
        </div>

        <div className="flex justify-end gap-3 pt-4 border-t text-white">
          <Button variant="secondary" onClick={onClose}>
            Cancel
          </Button>
          <Button variant="primary" onClick={handleSubmit}>
            Submit
          </Button>
        </div>
      </div>
    </Modal>
  )
}

export default WebServiceModal
