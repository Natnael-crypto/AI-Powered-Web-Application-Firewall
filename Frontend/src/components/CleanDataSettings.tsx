import {useState} from 'react'

const options = ['Not Clean', '1 days', '3 days', '7 days', '15 days', '30 days']

const CleanDataSettings = () => {
  const [protectionLog, setProtectionLog] = useState('Not Clean')
  const [statisticsData, setStatisticsData] = useState('Not Clean')

  return (
    <div className="p-6 bg-white  shadow-lg w-full">
      <h4 className="text-2xl font-bold text-gray-800 mb-6">Clean Data</h4>

      {/* Clean Protection Logs */}
      <div className="mb-8">
        <label className="block text-lg font-semibold text-gray-700 mb-3">
          Clean Protection Logs
        </label>
        <div className="flex flex-wrap gap-4">
          {options.map(option => (
            <label key={option} className="flex items-center cursor-pointer group">
              <input
                type="radio"
                name="protectionLog"
                value={option}
                checked={protectionLog === option}
                onChange={() => setProtectionLog(option)}
                className="w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"
              />
              <span className="ml-2 text-gray-700 group-hover:text-gray-900">
                {option}
              </span>
            </label>
          ))}
        </div>
      </div>

      {/* Clean Statistics Data */}
      <div>
        <label className="block text-lg font-semibold text-gray-700 mb-3">
          Clean Statistics Data
        </label>
        <div className="flex flex-wrap gap-4">
          {options.map(option => (
            <label key={option} className="flex items-center cursor-pointer group">
              <input
                type="radio"
                name="statisticsData"
                value={option}
                checked={statisticsData === option}
                onChange={() => setStatisticsData(option)}
                className="w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"
              />
              <span className="ml-2 text-gray-700 group-hover:text-gray-900">
                {option}
              </span>
            </label>
          ))}
        </div>
      </div>
    </div>
  )
}

export default CleanDataSettings
