import React from 'react'
import Modal from './Modal'
import {RuleResponse} from '../lib/types'

interface RuleDetailsModalProps {
  isOpen: boolean
  onClose: () => void
  rule: RuleResponse | null
}

const RuleDetailsModal: React.FC<RuleDetailsModalProps> = ({isOpen, onClose, rule}) => {
  if (!rule) return null

  const getStatusBadge = (isActive: boolean) => {
    const baseClasses = 'inline-block px-2 py-1 text-xs rounded'
    return isActive
      ? `${baseClasses} bg-green-100 text-green-800`
      : `${baseClasses} bg-gray-100 text-gray-800`
  }

  const formatTimestamp = (timestamp: string) => {
    const date = new Date(timestamp)
    return date.toLocaleString()
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Rule Details">
      <div className="space-y-6">
        {/* Basic Information */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-xs text-gray-500 mb-1">Rule ID</label>
            <div className="text-sm text-gray-800 font-mono">{rule.rule_id}</div>
          </div>

          <div>
            <label className="block text-xs text-gray-500 mb-1">Status</label>
            <div className={getStatusBadge(rule.is_active)}>
              {rule.is_active ? 'Active' : 'Inactive'}
            </div>
          </div>

          <div>
            <label className="block text-xs text-gray-500 mb-1">Category</label>
            <div className="text-sm text-gray-800 capitalize">{rule.category}</div>
          </div>

          <div>
            <label className="block text-xs text-gray-500 mb-1">Action</label>
            <div className="text-sm text-gray-800 capitalize">{rule.action}</div>
          </div>

          <div>
            <label className="block text-xs text-gray-500 mb-1">Created At</label>
            <div className="text-sm text-gray-800">
              {formatTimestamp(rule.created_at)}
            </div>
          </div>

          <div>
            <label className="block text-xs text-gray-500 mb-1">Applications</label>
            <div className="text-sm text-gray-800">
              {rule.applications?.length > 0
                ? rule.applications.map(app => app.application_name).join(', ')
                : 'All applications'}
            </div>
          </div>
        </div>

        {/* Rule Definition */}
        <div>
          <div className="flex items-center mb-2">
            <div className="flex-grow border-t border-gray-200"></div>
            <span className="mx-4 text-sm font-medium text-gray-500">
              RULE DEFINITION
            </span>
            <div className="flex-grow border-t border-gray-200"></div>
          </div>

          <pre className="bg-gray-50 p-4 rounded-md text-sm font-mono whitespace-pre-wrap overflow-x-auto">
            {rule.rule_string}
          </pre>
        </div>

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

export default RuleDetailsModal
