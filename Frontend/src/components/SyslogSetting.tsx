import {useState} from 'react'

export default function SyslogSettings() {
  const [serverAddress, setServerAddress] = useState('')
  const [serverPort, setServerPort] = useState('')

  const handleSave = () => {
    alert(`Saved:\nAddress: ${serverAddress}\nPort: ${serverPort}`)
  }

  return (
    <div className="w-full p-6 bg-white shadow-md rounded-2xl">
      <h2 className="text-2xl font-semibold text-gray-800 mb-6 flex items-center gap-2">
        <span>🖥️ Syslog</span>
        <span className="tooltip text-sm text-gray-500">ℹ️</span>
      </h2>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Syslog server address
          </label>
          <input
            type="text"
            value={serverAddress}
            onChange={e => setServerAddress(e.target.value)}
            placeholder="192.168.10.10"
            className="w-full p-3 border border-gray-300 rounded-xl shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Syslog server port
          </label>
          <input
            type="number"
            min="1"
            max="65535"
            value={serverPort}
            onChange={e => setServerPort(e.target.value)}
            placeholder="Must be in range 1 ~ 65535"
            className="w-full p-3 border border-gray-300 rounded-xl shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>
      <div className="mt-6 flex justify-end">
        <button
          onClick={handleSave}
          className="bg-blue-600 text-white font-semibold px-6 py-2 rounded-xl hover:bg-blue-700 transition"
        >
          Save
        </button>
      </div>
    </div>
  )
}
