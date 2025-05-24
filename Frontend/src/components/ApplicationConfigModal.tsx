import {useState, useEffect} from 'react'
import {
  useGetConfig,
  useUpdateListeningPort,
  useUpdateMaxDataSize,
  useUpdateRateLimit,
  useUpdateRemoteLogServer,
} from '../hooks/api/useApplication'
import Modal from './Modal'

interface ConfigModalProps {
  appId: string
  isOpen: boolean
  onClose: () => void
  data:
    | {
        id?: string
        rate_limit?: number
        window_size?: number
        block_time?: number
        detect_bot?: boolean
        max_post_data_size?: number
        host_name?: string
        tls?: boolean
      }
    | undefined
}

interface ConfigForm {
  rate_limit: number
  window_size: number
  block_time: number
  detect_bot: boolean
  hostname: string
  max_post_data_size: number
  tls: boolean
  listening_port: string
  remote_logServer: string
}

export default function ApplicationConfigModal({
  appId,
  isOpen,
  onClose,
  data,
}: ConfigModalProps) {
  const {mutate: updateListeningPort, isPending: isListeningPortUpdating} =
    useUpdateListeningPort()
  const {mutate: updateRateLimit, isPending: isRateLimitUpdating} = useUpdateRateLimit()
  const {mutate: updateRemoteLogServer, isPending: isRemoteLogUpdating} =
    useUpdateRemoteLogServer()
  const {data: serverConfig} = useGetConfig()

  const [formData, setFormData] = useState<ConfigForm>({
    rate_limit: 50,
    window_size: 10,
    block_time: 0,
    detect_bot: false,
    hostname: '',
    max_post_data_size: 5,
    tls: true,
    listening_port: '',
    remote_logServer: '',
  })

  useEffect(() => {
    if (data) {
      setFormData(prev => ({
        ...prev,
        ...data,
        ...serverConfig,
      }))
    }
  }, [data])

  const {mutate: updateMaxDataSize, isPending: isMaxDataLoading} = useUpdateMaxDataSize()

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const {name, value, type, checked} = e.target
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value,
    }))
  }

  const handleMaxDataSize = () => {
    updateMaxDataSize({
      application_id: appId,
      data: {
        max_post_data_size: Number(formData.max_post_data_size),
      },
    })
  }
  const handleUpdateRateLimit = () => {
    updateRateLimit({
      application_id: appId,
      data: {
        rate_limit: Number(formData.rate_limit),
        window_size: Number(formData.window_size),
        block_time: Number(formData.block_time),
      },
    })
  }

  const handleUpdateListeningPort = () => {
    updateListeningPort({
      listening_port: formData.listening_port,
    })
  }

  const handleUpdateRemoteLogServer = () => {
    updateRemoteLogServer({
      remote_logServer: formData.remote_logServer,
    })
  }

  if (!isOpen) return null

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Application Configuration">
      {/* Application-specific Configuration */}
      <div className="mb-8">
        <div className="flex items-center mb-4">
          <div className="flex-grow border-t border-gray-200"></div>
          <span className="mx-4 text-sm font-medium text-gray-500">
            APPLICATION SETTINGS
          </span>
          <div className="flex-grow border-t border-gray-200"></div>
        </div>

        {/* Rate Limit Section */}
        <div className="space-y-4 mb-6">
          <h4 className="text-sm font-semibold text-gray-700">Rate Limiting</h4>

          <div className="grid grid-cols-3 gap-4">
            <div>
              <label className="block text-xs text-gray-500 mb-1">Rate Limit</label>
              <input
                type="number"
                name="rate_limit"
                value={formData.rate_limit}
                onChange={handleChange}
                className="w-full px-3 py-2 border rounded-md"
              />
            </div>

            <div>
              <label className="block text-xs text-gray-500 mb-1">Window (sec)</label>
              <input
                type="number"
                name="window_size"
                value={formData.window_size}
                onChange={handleChange}
                className="w-full px-3 py-2 border rounded-md"
              />
            </div>

            <div>
              <label className="block text-xs text-gray-500 mb-1">Block Time (sec)</label>
              <input
                type="number"
                name="block_time"
                value={formData.block_time}
                onChange={handleChange}
                className="w-full px-3 py-2 border rounded-md"
              />
            </div>
          </div>

          <div className="flex justify-end">
            <button
              onClick={handleUpdateRateLimit}
              disabled={isRateLimitUpdating}
              className="px-4 py-2 text-sm bg-black text-white rounded hover:bg-gray-800 transition-colors"
            >
              {isRateLimitUpdating ? 'Saving...' : 'Save Rate Limits'}
            </button>
          </div>
        </div>
        <div className="space-y-3 mb-6">
          <div className="flex items-center gap-4">
            <div className="flex-1">
              <label className="block text-xs text-gray-500 mb-1">Maximum Post Data Size</label>
              <input
                type="text"
                name="max_post_data_size"
                value={formData.max_post_data_size}
                onChange={handleChange}
                className="w-full px-3 py-2 border rounded-md"
              />
            </div>
            <button
              onClick={handleMaxDataSize}
              disabled={isMaxDataLoading}
              className="self-end px-4 py-2 text-sm bg-black text-white rounded hover:bg-gray-800 transition-colors whitespace-nowrap"
            >
              {isMaxDataLoading ? 'Saving...' : 'Save Port'}
            </button>
          </div>
        </div>
      </div>

      {/* Footer */}
      <div className="pt-6 mt-6 border-t border-gray-100 flex justify-end">
        <button
          onClick={onClose}
          className="px-4 py-2 text-sm text-gray-600 bg-gray-100 rounded hover:bg-gray-200 transition-colors"
        >
          Close Configuration
        </button>
      </div>
    </Modal>
  )
}
