import { ChangeEvent, useEffect, useState } from 'react'
import { logFilterType } from '../lib/types'
import { useLogFilter } from '../store/LogFilter'
import { X, ChevronDown } from 'lucide-react'

function AttackLogFilter() {
  const {
    filters,
    tempFilter,
    setTempFilter,
    addFilter,
    removeFilter,
    applyFilters,
    setFilter,
  } = useLogFilter()

  const [timeRange, setTimeRange] = useState<{
    value: string
    label: string
    start?: number
    end?: number
  }>({ value: '', label: '' })

  const [customRange, setCustomRange] = useState({
    startDate: '',
    endDate: '',
    startTime: '00:00',
    endTime: '23:59',
  })

  const [isPopoverOpen, setIsPopoverOpen] = useState(false)

  const timePresets = [
    { label: 'Last 1 hour', value: 'last_1_hour' },
    { label: 'Last 24 hours', value: 'last_24_hours' },
    { label: 'Last 7 days', value: 'last_7_days' },
    { label: 'Last 30 days', value: 'last_30_days' },
    { label: 'Custom Range', value: 'custom' },
  ]

  const handleTypeChange = (e: ChangeEvent<HTMLSelectElement>) => {
    setTempFilter('key', e.target.value)
  }

  const handleValueChange = (e: ChangeEvent<HTMLInputElement>) => {
    setTempFilter('value', e.target.value)
  }

  const filterKeys = Object.keys(logFilterType)

  const activeFilters = Object.entries(filters).filter(([_, value]) => value !== '')

  const handleDateChange = (
    e: React.ChangeEvent<HTMLInputElement>,
    type: 'start' | 'end',
    field: 'date' | 'time'
  ) => {
    setCustomRange(prev => ({
      ...prev,
      [`${type}${field.charAt(0).toUpperCase() + field.slice(1)}`]: e.target.value,
    }))
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

      setFilter('start_date', start.toString())
      setFilter('end_date', end.toString())
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
    }

    const label = timePresets.find(p => p.value === value)?.label || ''

    if (value !== 'custom') {
      setTimeRange({ value, label, start })
      setIsPopoverOpen(false)

      // Apply to filters
      setFilter('last_hours', start.toString())
    } else {
      setTimeRange({ value: 'custom', label: 'Custom Range' })
    }
  }

  const timesKey = ['last_hours', 'start_date', 'end_date']

  const formatTimestamp = (key: string, value: any) => {
    if (timesKey.includes(key)) {
      const timestamp = parseInt(value, 10);
      if (!isNaN(timestamp)) {
        const date = new Date(timestamp);
        return new Intl.DateTimeFormat('en-US', {
          weekday: 'short',
          year: 'numeric',
          month: 'short',
          day: 'numeric',
          hour: 'numeric',
          minute: 'numeric',
          second: 'numeric',
        }).format(date);
      }
    }

    return value;
  }



  return (
    <div className="w-full bg-white p-6 rounded-lg shadow flex flex-col gap-6">
      <div className="flex flex-col md:flex-row gap-4">
        <div className="flex-1">
          <label className="block mb-1 text-sm font-semibold text-gray-700">Filter Type</label>
          <select
            value={tempFilter.key}
            onChange={handleTypeChange}
            className="w-full px-4 py-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
          >
            <option value="">Select Type</option>
            {filterKeys.map((key) => (
              <option key={key} value={key}>
                {logFilterType[key as keyof typeof logFilterType]}
              </option>
            ))}
          </select>
        </div>

        <div className="flex-1">
          <label className="block mb-1 text-sm font-semibold text-gray-700">Value</label>
          <input
            type="text"
            value={tempFilter.value}
            onChange={handleValueChange}
            className="w-full px-4 py-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div className="flex-1">
          <div className="relative">
            <button
              onClick={() => setIsPopoverOpen(!isPopoverOpen)}
              className="flex items-center gap-2 px-4 py-2 border rounded-md text-sm text-gray-700 hover:bg-gray-100"
            >
              <span>{timeRange.label || 'Select Time Range'}</span>
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
                        timeRange.value === preset.value ? 'bg-gray-100 font-medium' : ''
                      }`}
                    >
                      {preset.label}
                    </button>
                  ))}

                  {timeRange.value === 'custom' && (
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

        <div className="flex gap-2 items-end">
          <button
            onClick={addFilter}
            className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
          >
            Add
          </button>
          <button
            onClick={applyFilters}
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Apply
          </button>
        </div>
      </div>

      {activeFilters.length > 0 && (
        <div>
          <p className="font-semibold mb-2 text-gray-800">Active Filters:</p>
          <div className="flex flex-wrap gap-2">
            {activeFilters.map(([key, value]) => (
              <span
                key={key}
                className="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm flex items-center"
              >
                {logFilterType[key as keyof typeof logFilterType]}: {formatTimestamp(key,value)}
                <button
                  className="ml-2 text-red-500 hover:text-red-700"
                  onClick={() => removeFilter(key)}
                >
                  <X size={14} />
                </button>
              </span>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

export default AttackLogFilter
