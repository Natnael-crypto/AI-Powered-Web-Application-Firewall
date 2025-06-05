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
      <div key={index} className="text-sm py-1 px-2 hover:bg-gray-100 rounded">
        <span className="text-gray-500 font-medium">{header.split(':')[0]}:</span>{' '}
        <span className="text-gray-800 break-all">
          {header.split(':').slice(1).join(':').trim()}
        </span>
      </div>
    ))
  }

  console.log(request)
  const formatTimestamp = (timestamp: string) => {
    const date = new Date(timestamp)
    return date.toLocaleString()
  }

  const getStatusBadge = (status: string) => {
    const baseClasses =
      'px-3 py-1 rounded-full text-xs font-medium inline-flex items-center'
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
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 p-1">
          <div className="space-y-1">
            <label className="block text-xs font-medium text-gray-500 uppercase tracking-wider">
              Application
            </label>
            <div className="text-sm text-gray-800 font-medium">
              {request?.application_name}
            </div>
          </div>

          <div className="space-y-1">
            <label className="block text-xs font-medium text-gray-500 uppercase tracking-wider">
              Status
            </label>
            <div className={getStatusBadge(request?.status ?? '')}>{request?.status}</div>
          </div>

          <div className="space-y-1">
            <label className="block text-xs font-medium text-gray-500 uppercase tracking-wider">
              Client IP
            </label>
            <div className="text-sm text-gray-800 font-mono">{request?.client_ip}</div>
          </div>

          <div className="space-y-1">
            <label className="block text-xs font-medium text-gray-500 uppercase tracking-wider">
              Geo Location
            </label>
            <div className="text-sm text-gray-800">{request?.geo_location}</div>
          </div>

          <div className="space-y-1">
            <label className="block text-xs font-medium text-gray-500 uppercase tracking-wider">
              Method
            </label>
            <div className="text-sm text-gray-800 font-medium">
              {request?.request_method}
            </div>
          </div>

          <div className="space-y-1">
            <label className="block text-xs font-medium text-gray-500 uppercase tracking-wider">
              Response Code
            </label>
            <div className="text-sm font-semibold text-gray-800">
              {request?.response_code}
            </div>
          </div>

          <div className="space-y-1">
            <label className="block text-xs font-medium text-gray-500 uppercase tracking-wider">
              Timestamp
            </label>
            <div className="text-sm text-gray-800">
              {formatTimestamp(request?.timestamp ?? '')}
            </div>
          </div>

          <div className="space-y-1">
            <label className="block text-xs font-medium text-gray-500 uppercase tracking-wider">
              Request ID
            </label>
            <div className="text-sm text-gray-800 font-mono break-all">
              {request?.request_id}
            </div>
          </div>
        </div>

        {/* Divider */}
        <div className="border-t border-gray-200 my-2"></div>

        {/* Request URL */}
        <div className="space-y-2">
          <label className="block text-xs font-medium text-gray-500 uppercase tracking-wider">
            Request URL
          </label>
          <div className="text-sm text-gray-800 break-all p-3 bg-gray-50 rounded border border-gray-200">
            {request?.request_url}
          </div>
        </div>

        {/* Threat Detection Section */}
        {request?.threat_detected && (
          <div className="border border-red-200 bg-red-50 rounded-lg p-4 space-y-3">
            <div className="flex items-center">
              <div className="flex-grow border-t border-red-200"></div>
              <span className="mx-4 text-xs font-semibold uppercase tracking-wider text-red-600">
                Threat Detected
              </span>
              <div className="flex-grow border-t border-red-200"></div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-1">
                <label className="block text-xs text-red-600 font-medium">
                  Threat Type
                </label>
                <div className="text-sm text-red-800 font-medium">
                  {request.threat_type}
                </div>
              </div>

              <div className="space-y-1">
                <label className="block text-xs text-red-600 font-medium">
                  Matched Rules
                </label>
                <div className="text-sm text-red-800">{request.matched_rules}</div>
              </div>

              <div className="space-y-1">
                <label className="block text-xs text-red-600 font-medium">
                  Bot Detected
                </label>
                <div className="text-sm text-red-800 font-medium">
                  {request.bot_detected ? 'Yes' : 'No'}
                </div>
              </div>

              <div className="space-y-1">
                <label className="block text-xs text-red-600 font-medium">
                  Rate Limited
                </label>
                <div className="text-sm text-red-800 font-medium">
                  {request.rate_limited ? 'Yes' : 'No'}
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Headers Section */}
        <div className="space-y-2">
          <div className="flex items-center">
            <div className="flex-grow border-t border-gray-200"></div>
            <span className="mx-4 text-xs font-semibold uppercase tracking-wider text-gray-500">
              Request Headers
            </span>
            <div className="flex-grow border-t border-gray-200"></div>
          </div>

          <div className="bg-gray-50 p-2 rounded border border-gray-200 space-y-1 max-h-60 overflow-y-auto">
            {formatHeaders(request?.headers ?? '')}
          </div>
        </div>

        {/* User Agent */}
        {request?.user_agent && (
          <div className="space-y-2">
            <label className="block text-xs font-medium text-gray-500 uppercase tracking-wider">
              User Agent
            </label>
            <div className="text-sm text-gray-800 bg-gray-50 p-3 rounded border border-gray-200 break-all">
              {request.user_agent}
            </div>
          </div>
        )}

        {/* Body */}
        {request?.body && (
          <div className="space-y-2">
            <label className="block text-xs font-medium text-gray-500 uppercase tracking-wider">
              Body
            </label>
            <div className="text-sm text-gray-800 bg-gray-50 p-3 rounded border border-gray-200 break-all max-h-30 overflow-y-auto">
              {request.body}
            </div>
          </div>
        )}

        {/* AI Analysis */}
        {request && (
          <div>
            <div className="flex items-center mb-2">
              <div className="flex-grow border-t border-gray-200"></div>
              <span className="mx-4 text-xs font-semibold uppercase tracking-wider text-gray-500">
                AI Analysis
              </span>
              <div className="flex-grow border-t border-gray-200"></div>
            </div>
            <div className="text-sm text-gray-800 bg-gray-50 p-3 rounded space-y-1">
              <div>
                {request.ai_result === false ? (
                  <span className="text-green-600 font-medium">Allowed by AI</span>
                ) : (
                  <span className="text-red-600 font-medium">Blocked by AI</span>
                )}
              </div>
              {request.ai_threat_type && (
                <div>
                  <span className="font-semibold">AI Threat Type:</span>{' '}
                  {request.ai_threat_type}
                </div>
              )}
            </div>
          </div>
        )}

        {/* Footer */}
        <div className="pt-4 border-t border-gray-200 flex justify-end">
          <button
            onClick={onClose}
            className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
          >
            Close
          </button>
        </div>
      </div>
    </Modal>
  )
}

export default RequestDetailsModal
