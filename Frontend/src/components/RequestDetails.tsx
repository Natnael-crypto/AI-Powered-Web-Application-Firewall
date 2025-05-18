import React from 'react'
import Modal from './Modal'
import {RequestLog} from '../lib/types'

interface RequestDetailsModalProps {
  isOpen: boolean
  onClose: () => void
  request: RequestLog | undefined
}

const RequestDetailsModal: React.FC<RequestDetailsModalProps> = ({
  isOpen,
  onClose,
  request,
}) => {
  const formatHeaders = (headers: string) => {
    return headers.split(', ').map((header, index) => (
      <div key={index} className="text-sm">
        <span className="text-gray-500">{header.split(':')[0]}:</span>{' '}
        <span className="text-gray-800">
          {header.split(':').slice(1).join(':').trim()}
        </span>
      </div>
    ))
  }

  const formatTimestamp = (timestamp: string) => {
    const date = new Date(timestamp)
    return date.toLocaleString()
  }

  const getStatusBadge = (status: string) => {
    const baseClasses = 'px-2 py-1 rounded-full text-xs font-medium'
    if (status === 'blocked') {
      return `${baseClasses} bg-red-100 text-red-800`
    } else if (status === 'allowed') {
      return `${baseClasses} bg-green-100 text-green-800`
    }
    return `${baseClasses} bg-gray-100 text-gray-800`
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Request Details">
      <div className="space-y-6">
        {/* Basic Info Section */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-xs text-gray-500 mb-1">Application</label>
            <div className="text-sm text-gray-800">{request?.application_name}</div>
          </div>

          <div>
            <label className="block text-xs text-gray-500 mb-1">Status</label>
            <div className={getStatusBadge(request?.status ?? '')}>{request?.status}</div>
          </div>

          <div>
            <label className="block text-xs text-gray-500 mb-1">Client IP</label>
            <div className="text-sm text-gray-800">{request?.client_ip}</div>
          </div>

          <div>
            <label className="block text-xs text-gray-500 mb-1">Geo Location</label>
            <div className="text-sm text-gray-800">{request?.geo_location}</div>
          </div>

          <div>
            <label className="block text-xs text-gray-500 mb-1">Method</label>
            <div className="text-sm text-gray-800">{request?.request_method}</div>
          </div>

          <div>
            <label className="block text-xs text-gray-500 mb-1">Response Code</label>
            <div className="text-sm text-gray-800">{request?.response_code}</div>
          </div>

          <div>
            <label className="block text-xs text-gray-500 mb-1">Timestamp</label>
            <div className="text-sm text-gray-800">
              {formatTimestamp(request?.timestamp ?? '')}
            </div>
          </div>

          <div>
            <label className="block text-xs text-gray-500 mb-1">Request ID</label>
            <div className="text-sm text-gray-800 font-mono">{request?.request_id}</div>
          </div>
        </div>

        {/* Request URL */}
        <div>
          <label className="block text-xs text-gray-500 mb-1">Request URL</label>
          <div className="text-sm text-gray-800 break-all p-2 bg-gray-50 rounded">
            {request?.request_url}
          </div>
        </div>

        {/* Threat Detection Section */}
        {request?.threat_detected && (
          <div className="border border-red-200 bg-red-50 rounded p-4">
            <div className="flex items-center mb-2">
              <div className="flex-grow border-t border-red-200"></div>
              <span className="mx-4 text-sm font-medium text-red-500">
                THREAT DETECTED
              </span>
              <div className="flex-grow border-t border-red-200"></div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-xs text-red-500 mb-1">Threat Type</label>
                <div className="text-sm text-red-800">{request.threat_type}</div>
              </div>

              <div>
                <label className="block text-xs text-red-500 mb-1">Matched Rules</label>
                <div className="text-sm text-red-800">{request.matched_rules}</div>
              </div>

              <div>
                <label className="block text-xs text-red-500 mb-1">Bot Detected</label>
                <div className="text-sm text-red-800">
                  {request.bot_detected ? 'Yes' : 'No'}
                </div>
              </div>

              <div>
                <label className="block text-xs text-red-500 mb-1">Rate Limited</label>
                <div className="text-sm text-red-800">
                  {request.rate_limited ? 'Yes' : 'No'}
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Headers Section */}
        <div>
          <div className="flex items-center mb-2">
            <div className="flex-grow border-t border-gray-200"></div>
            <span className="mx-4 text-sm font-medium text-gray-500">
              REQUEST HEADERS
            </span>
            <div className="flex-grow border-t border-gray-200"></div>
          </div>

          <div className="bg-gray-50 p-3 rounded space-y-1 max-h-40 overflow-y-auto">
            {formatHeaders(request?.headers ?? '')}
          </div>
        </div>

        {/* User Agent */}
        {request?.user_agent && (
          <div>
            <label className="block text-xs text-gray-500 mb-1">User Agent</label>
            <div className="text-sm text-gray-800 bg-gray-50 p-2 rounded break-all">
              {request.user_agent}
            </div>
          </div>
        )}

        {/* AI Analysis */}
        {request?.ai_result && (
          <div>
            <div className="flex items-center mb-2">
              <div className="flex-grow border-t border-gray-200"></div>
              <span className="mx-4 text-sm font-medium text-gray-500">AI ANALYSIS</span>
              <div className="flex-grow border-t border-gray-200"></div>
            </div>
            <div className="text-sm text-gray-800 bg-gray-50 p-3 rounded">
              {request.ai_result}
            </div>
          </div>
        )}

        {/* Footer */}
        <div className="pt-4 border-t border-gray-100 flex justify-end">
          <button
            onClick={onClose}
            className="px-4 py-2 text-sm text-gray-600 bg-gray-100 rounded-md hover:bg-gray-200 transition-colors"
          >
            Close
          </button>
        </div>
      </div>
    </Modal>
  )
}

export default RequestDetailsModal
