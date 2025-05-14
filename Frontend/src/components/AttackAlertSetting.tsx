import { useEffect, useState } from 'react'
import { Info } from 'lucide-react'
import {
  useGetSysEmail,
  useAddSysEmail,
  useUpdateSysEmail,
} from '../hooks/api/useSystemEmail'

const AttackAlertSettings = () => {
  const [alertType, setAlertType] = useState<'Telegram' | 'Email'>('Email')
  const [webhook, setWebhook] = useState('')
  const [isActive, setIsActive] = useState(false)
  const [emailFetched, setEmailFetched] = useState(false)

  const { data, isSuccess } = useGetSysEmail()

  useEffect(() => {
    console.log(data)
    if (isSuccess && data && !emailFetched) {
      setWebhook(data.email || '')
      setIsActive(data.active || false)
      setEmailFetched(true)
    }
  }, [isSuccess, data, emailFetched])

  const handleSave = async () => {
    if (!webhook || webhook.trim() === '') {
      alert('Webhook email must not be empty.')
      return
    }

    try {
      if (!data?.email || data?.email === '') {
        await useAddSysEmail(webhook, isActive)
        alert('Email added successfully.')
      } else {
        await useUpdateSysEmail(webhook, isActive)
        alert('Email updated successfully.')
      }
    } catch (error) {
      alert('Failed to save email settings.')
      console.error(error)
    }
  }

  return (
    <div className="p-6 bg-white shadow-lg w-full">
      {/* Title */}
      <div className="flex items-center mb-4">
        <h4 className="text-lg font-semibold text-gray-800 mr-2">Attack Alert</h4>
        <Info size={16} className="text-blue-500" />
      </div>

      {/* Radio buttons */}
      <div className="flex gap-6 mb-4">
        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="radio"
            name="alertType"
            value="Email"
            checked={alertType === 'Email'}
            onChange={() => setAlertType('Email')}
            className="w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"
          />
          <span className="text-gray-700">Email</span>
        </label>
      </div>

      {/* Webhook input */}
      <input
        type="text"
        placeholder="waf-alert-server-noreply@org.com"
        value={webhook}
        onChange={e => setWebhook(e.target.value)}
        className="w-full px-4 py-2 mb-4 border border-gray-300 text-sm placeholder-gray-400"
      />

      {/* Checkboxes */}
      <div className="flex gap-6 mb-4">
        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            checked={isActive}
            onChange={() => setIsActive(!isActive)}
            className="w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"
          />
          <span className="text-gray-700">Allow notifications via this email</span>
        </label>
      </div>

      {/* Save Button */}
      <div className="mt-4 flex justify-end">
        <button
          onClick={handleSave}
          className="bg-blue-600 text-white font-semibold px-6 py-2 hover:bg-blue-700 transition"
        >
          Save
        </button>
      </div>
    </div>
  )
}

export default AttackAlertSettings
