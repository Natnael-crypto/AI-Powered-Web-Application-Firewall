import {useState} from 'react'
import {Info} from 'lucide-react'

const AttackAlertSettings = () => {
  const [alertType, setAlertType] = useState<'Telegram' | 'Email'>('Email')
  const [webhook, setWebhook] = useState('')
  const [attackEvent, setAttackEvent] = useState(false)
  const [rateLimiting, setRateLimiting] = useState(false)

  return (
    <div className="p-6 bg-white  shadow-lg w-full">
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
        placeholder="https://api.xxx.com/robot/send?access_token=xxxxxx"
        value={webhook}
        onChange={e => setWebhook(e.target.value)}
        className="w-full px-4 py-2 mb-4 border border-gray-300  text-sm placeholder-gray-400"
      />

      {/* Checkboxes */}
      <div className="flex gap-6">
        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            checked={attackEvent}
            onChange={() => setAttackEvent(!attackEvent)}
            className="w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"
          />
          <span className="text-gray-700">Attack Event</span>
        </label>
        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            checked={rateLimiting}
            onChange={() => setRateLimiting(!rateLimiting)}
            className="w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"
          />
          <span className="text-gray-700">Rate Limiting</span>
        </label>
      </div>
    </div>
  )
}

export default AttackAlertSettings
