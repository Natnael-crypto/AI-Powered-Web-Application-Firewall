import { useState, useEffect } from 'react'
import { NotificationRule, NotificationRuleInput, UpdateNotificationRuleInput } from '../lib/types'
import { useUpdateNotificationRule } from '../hooks/api/useNotificationRules'
import Modal from './Modal'

interface NotificationRuleModalProps {
  rule: NotificationRule | undefined
  isOpen: boolean
  onClose: () => void
}

export default function NotificationRuleModal({
  rule,
  isOpen,
  onClose,
}: NotificationRuleModalProps) {
  const [formData, setFormData] = useState({
    threshold: 0,
    time_window: 0,
    is_active: false,
  })

  const { mutate: updateRule, isPending } = useUpdateNotificationRule()

  useEffect(() => {
    if (rule) {
      setFormData({
        threshold: rule.threshold,
        time_window: rule.time_window,
        is_active: rule.is_active,
      })
    }
  }, [rule])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value, type, checked } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : Number(value),
    }))
  }

  const handleSave = () => {
    if (!rule) return

    const data: NotificationRuleInput = {
      threshold: formData.threshold,
      time_window: formData.time_window,
      is_active: formData.is_active,
    }

    const updateRuleData: UpdateNotificationRuleInput={
      rule_id: rule.id,
      data: data
    }

    updateRule(
        updateRuleData
      ,
      {
        onSuccess: () => {
          onClose()
        },
      }
    )
  }

  if (!isOpen) return null

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Edit Notification Rule">
      <p className='text-center font-lg'>{rule?.name}</p>
      <div className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-600">Threshold</label>
          <input
            type="number"
            name="threshold"
            value={formData.threshold}
            onChange={handleChange}
            className="w-full px-3 py-2 border rounded-md"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-600">Time Window (min)</label>
          <input
            type="number"
            name="time_window"
            value={formData.time_window}
            onChange={handleChange}
            className="w-full px-3 py-2 border rounded-md"
          />
        </div>

        <div className="flex items-center gap-2">
          <input
            type="checkbox"
            name="is_active"
            checked={formData.is_active}
            onChange={handleChange}
            className="w-4 h-4"
          />
          <label className="text-sm text-gray-600">Active</label>
        </div>
      </div>

      <div className="pt-6 mt-6 border-t border-gray-100 flex justify-end">
        <button
          onClick={handleSave}
          disabled={isPending}
          className="px-4 py-2 text-sm bg-black text-white rounded hover:bg-gray-800 transition-colors"
        >
          {isPending ? 'Saving...' : 'Save Rule'}
        </button>
        <button
          onClick={onClose}
          className="ml-2 px-4 py-2 text-sm text-gray-600 bg-gray-100 rounded hover:bg-gray-200 transition-colors"
        >
          Cancel
        </button>
      </div>
    </Modal>
  )
}
