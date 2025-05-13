import React, { useState, useEffect, Dispatch, SetStateAction } from 'react'
import { ChevronDown } from 'lucide-react'
import { useGetApplications } from '../hooks/api/useApplication'

interface FilterBarProps {
  selectedApp: string
  setSelectedApp: Dispatch<SetStateAction<string>>
  timeRange: any
  setTimeRange: Dispatch<SetStateAction<any>>
}

const FilterBar = ({
  selectedApp,
  setSelectedApp,
  timeRange,
  setTimeRange,
}: FilterBarProps) => {
  const [customRange, setCustomRange] = useState({
    startDate: '',
    endDate: '',
    startTime: '00:00',
    endTime: '23:59',
  })
  const [isPopoverOpen, setIsPopoverOpen] = useState(false)
  const [_, setAvailableEndDates] = useState<string[]>([])

  const {
    data: applicationsData,
    isLoading,
    isError,
  } = useGetApplications()

  const timePresets = [
    { label: 'Last 1 hour', value: 'last_1_hour' },
    { label: 'Last 24 hours', value: 'last_24_hours' },
    { label: 'Last 7 days', value: 'last_7_days' },
    { label: 'Last 30 days', value: 'last_30_days' },
    { label: 'Custom Range', value: 'custom' },
  ]

  useEffect(() => {
    if (customRange.startDate) {
      const start = new Date(customRange.startDate)
      const end = new Date(start)
      end.setDate(start.getDate() + 30)
      setAvailableEndDates([end.toISOString().split('T')[0]])
    }
  }, [customRange.startDate])

  const handleDateChange = (
    e: React.ChangeEvent<HTMLInputElement>,
    type: 'start' | 'end',
    field: 'date' | 'time'
  ) => {
    setCustomRange(prev => {
      const updatedRange = { ...prev }
      updatedRange[
        `${type}${field.charAt(0).toUpperCase() + field.slice(1)}` as keyof typeof customRange
      ] = e.target.value
      return updatedRange
    })
  }

  const applyCustomRange = () => {
    const { startDate, endDate, startTime, endTime } = customRange
    if (startDate && endDate && startTime && endTime) {
      const start = new Date(`${startDate}T${startTime}`).getTime()
      const end = new Date(`${endDate}T${endTime}`).getTime()

      setTimeRange({
        value: 'custom',
        label: 'Custom Range',
        start,
        end,
      })
      setIsPopoverOpen(false)
    }
  }

  const handlePresetSelect = (value: string) => {
    const now = Date.now()
    let start = 0

    switch (value) {
      case 'last_1_hour':
        start = now - 1 * 60 * 60 * 1000
        break
      case 'last_24_hours':
        start = now - 24 * 60 * 60 * 1000
        break
      case 'last_7_days':
        start = now - 7 * 24 * 60 * 60 * 1000
        break
      case 'last_30_days':
        start = now - 30 * 24 * 60 * 60 * 1000
        break
      default:
        break
    }

    if (value !== 'custom') {
      const label = timePresets.find(p => p.value === value)?.label || ''
      setTimeRange({
        value,
        label,
        start,
        end: '', // dynamic/live until now
      })
      setIsPopoverOpen(false)
    } else {
      setTimeRange(value) // set to 'custom'
    }
  }

  return (
    <div className="flex items-center justify-between gap-6 p-4 bg-white shadow-md rounded-md">
      {/* Application Filter */}
      <div className="flex items-center gap-2">
        <span className="text-sm text-gray-700">Application:</span>
        <select
          value={selectedApp}
          onChange={e => setSelectedApp(e.target.value)}
          className="px-4 py-2 border rounded-md text-sm text-gray-700"
        >
          <option value="">All Applications</option>
          {isLoading ? (
            <option disabled>Loading...</option>
          ) : isError ? (
            <option disabled>Error loading applications</option>
          ) : (
            applicationsData?.map((app: any) => (
              <option key={app.application_id} value={app.application_id}>
                {app.hostname}
              </option>
            ))
          )}
        </select>
      </div>

      {/* Time Range Filter */}
      <div className="relative">
        <button
          onClick={() => setIsPopoverOpen(!isPopoverOpen)}
          className="flex items-center gap-2 px-4 py-2 border rounded-md text-sm text-gray-700 hover:bg-gray-100"
        >
          <span>
            {typeof timeRange === 'string'
              ? timePresets.find(p => p.value === timeRange)?.label
              : timeRange?.label || 'Select Time Range'}
          </span>
          <ChevronDown className="w-4 h-4 text-gray-500" />
        </button>

        {isPopoverOpen && (
          <div className="absolute z-10 mt-2 p-4 w-64 bg-white border rounded-md shadow-md right-0">
            <div className="flex flex-col gap-2">
              {timePresets.map(preset => (
                <button
                  key={preset.value}
                  onClick={() => handlePresetSelect(preset.value)}
                  className={`text-left px-3 py-1 text-sm rounded hover:bg-gray-100 ${
                    (typeof timeRange === 'string' && timeRange === preset.value) ||
                    (typeof timeRange === 'object' && timeRange?.value === preset.value)
                      ? 'bg-gray-100 font-medium'
                      : ''
                  }`}
                >
                  {preset.label}
                </button>
              ))}

              {(timeRange === 'custom' || timeRange?.value === 'custom') && (
                <div className="flex flex-col gap-1 mt-2">
                  <label className="text-xs text-gray-500">Start Date</label>
                  <input
                    type="date"
                    value={customRange.startDate}
                    onChange={e => handleDateChange(e, 'start', 'date')}
                    className="border rounded px-2 py-1 text-sm"
                  />

                  <label className="text-xs text-gray-500 mt-2">Start Time</label>
                  <input
                    type="time"
                    value={customRange.startTime}
                    onChange={e => handleDateChange(e, 'start', 'time')}
                    className="border rounded px-2 py-1 text-sm"
                  />

                  <label className="text-xs text-gray-500 mt-2">End Date</label>
                  <input
                    type="date"
                    value={customRange.endDate}
                    onChange={e => handleDateChange(e, 'end', 'date')}
                    min={customRange.startDate}
                    className="border rounded px-2 py-1 text-sm"
                  />

                  <label className="text-xs text-gray-500 mt-2">End Time</label>
                  <input
                    type="time"
                    value={customRange.endTime}
                    onChange={e => handleDateChange(e, 'end', 'time')}
                    min={customRange.startTime}
                    className="border rounded px-2 py-1 text-sm"
                  />

                  <button
                    className="mt-3 px-3 py-1 bg-blue-600 text-white text-sm rounded hover:bg-blue-700"
                    onClick={applyCustomRange}
                  >
                    Apply
                  </button>
                </div>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

export default FilterBar
