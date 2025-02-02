import React, {useState} from 'react'
import {FilterValues} from '../../lib/types'

interface LogFilterProps {
  onFilter: (filters: FilterValues) => void
  logtype?: 'log' | 'event'
  onLogtypeChange?: (newType: 'log' | 'event') => void
}

const LogFilter: React.FC<LogFilterProps> = ({onFilter, logtype, onLogtypeChange}) => {
  const [selectedLogtype, setSelectedLogtype] = useState<'log' | 'event'>(
    logtype ?? 'log',
  )

  const handleLogtypeChange = (newType: 'log' | 'event') => {
    if (newType !== selectedLogtype) {
      setSelectedLogtype(newType)
      if (onLogtypeChange) {
        onLogtypeChange(newType)
      }
    }
  }

  const [filters, setFilters] = useState<FilterValues>({
    ipAddress: '',
    port: '',
    domain: '',
    startAt: '',
    endAt: '',
  })

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const {name, value} = e.target
    setFilters({
      ...filters,
      [name]: value,
    })

    handleFilter()
  }

  const handleFilter = () => {
    onFilter(filters)
  }

  return (
    <div className="bg-white p-6 rounded-lg shadow-md mb-6 flex gap-10 items-end">
      {/* Log type selection buttons */}
      <div className="flex gap-2 ">
        <button
          onClick={() => handleLogtypeChange('log')}
          className={`px-4 py-2 rounded-md border border-gray-500 ${
            selectedLogtype === 'log'
              ? 'bg-blue-500 text-white'
              : 'bg-gray-200 text-gray-800'
          }`}
        >
          Log
        </button>
        <button
          onClick={() => handleLogtypeChange('event')}
          className={`px-4 py-2 rounded-md border border-gray-500 ${
            selectedLogtype === 'event'
              ? 'bg-blue-500 text-white'
              : 'bg-gray-200 text-gray-800'
          }`}
        >
          Event
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-4">
        {/* IP Address Input */}
        <div>
          <label htmlFor="ipAddress" className="block text-sm font-medium text-gray-700">
            IP Address
          </label>
          <input
            type="text"
            id="ipAddress"
            name="ipAddress"
            value={filters.ipAddress}
            onChange={handleInputChange}
            className=" block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        {/* Port Input */}
        <div>
          <label htmlFor="port" className="block text-sm font-medium text-gray-700">
            Port
          </label>
          <input
            type="text"
            id="port"
            name="port"
            value={filters.port}
            onChange={handleInputChange}
            className=" block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        {/* Domain Input */}
        <div>
          <label htmlFor="domain" className="block text-sm font-medium text-gray-700">
            Domain
          </label>
          <input
            type="text"
            id="domain"
            name="domain"
            value={filters.domain}
            onChange={handleInputChange}
            className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        {/* Start At Date Input */}
        <div>
          <label htmlFor="startAt" className="block text-sm font-medium text-gray-700">
            Start At
          </label>
          <input
            type="date"
            id="startAt"
            name="startAt"
            value={filters.startAt}
            onChange={handleInputChange}
            className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        {/* End At Date Input */}
        <div>
          <label htmlFor="endAt" className="block text-sm font-medium text-gray-700">
            End At
          </label>
          <input
            type="date"
            id="endAt"
            name="endAt"
            value={filters.endAt}
            onChange={handleInputChange}
            className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
        </div>
      </div>
    </div>
  )
}

export default LogFilter
