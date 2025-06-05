import {ChangeEvent, useState} from 'react'
import {logFilterType} from '../lib/types'
import {useLogFilter} from '../store/LogFilter'
import {X, ChevronDown, Filter, Plus, Check} from 'lucide-react'

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
  }>({value: '', label: ''})

  const [customRange, setCustomRange] = useState({
    startDate: '',
    endDate: '',
    startTime: '00:00',
    endTime: '23:59',
  })

  const [isPopoverOpen, setIsPopoverOpen] = useState(false)

  const timePresets = [
    {label: 'Last 1 hour', value: 'last_1_hour'},
    {label: 'Last 24 hours', value: 'last_24_hours'},
    {label: 'Last 7 days', value: 'last_7_days'},
    {label: 'Last 30 days', value: 'last_30_days'},
    {label: 'Custom Range', value: 'custom'},
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
    field: 'date' | 'time',
  ) => {
    setCustomRange(prev => ({
      ...prev,
      [`${type}${field.charAt(0).toUpperCase() + field.slice(1)}`]: e.target.value,
    }))
  }

  const applyCustomRange = () => {
    const {startDate, endDate, startTime, endTime} = customRange
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
      setTimeRange({value, label, start})
      setIsPopoverOpen(false)
      setFilter('last_hours', start.toString())
    } else {
      setTimeRange({value: 'custom', label: 'Custom Range'})
    }
  }

  const timesKey = ['last_hours', 'start_date', 'end_date']

  const formatTimestamp = (key: string, value: any) => {
    if (timesKey.includes(key)) {
      const timestamp = parseInt(value, 10)
      if (!isNaN(timestamp)) {
        const date = new Date(timestamp)
        return new Intl.DateTimeFormat('en-US', {
          weekday: 'short',
          year: 'numeric',
          month: 'short',
          day: 'numeric',
          hour: 'numeric',
          minute: 'numeric',
          second: 'numeric',
        }).format(date)
      }
    }
    return value
  }

  return (
    <div className="w-full bg-white p-6 rounded-xl shadow-md border border-gray-100 flex flex-col gap-6">
      <div className="flex items-center gap-2">
        <Filter className="w-5 h-5 text-blue-600" />
        <h3 className="text-lg font-semibold text-gray-800">Attack Log Filters</h3>
      </div>

      <div className="flex flex-col md:flex-row gap-4">
        <div className="flex-1">
          <label className="block mb-2 text-sm font-medium text-gray-700">
            Filter Type
          </label>
          <div className="relative">
            <select
              value={tempFilter.key}
              onChange={handleTypeChange}
              className="w-full px-4 py-2.5 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 appearance-none bg-white"
            >
              <option value="">Select filter type</option>
              {filterKeys.map(key => (
                <option key={key} value={key}>
                  {logFilterType[key as keyof typeof logFilterType]}
                </option>
              ))}
            </select>
            <ChevronDown className="absolute right-3 top-3 h-4 w-4 text-gray-400 pointer-events-none" />
          </div>
        </div>

        <div className="flex-1">
          <label className="block mb-2 text-sm font-medium text-gray-700">
            Filter Value
          </label>
          <input
            type="text"
            value={tempFilter.value}
            onChange={handleValueChange}
            className="w-full px-4 py-2.5 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            placeholder="Enter filter value"
          />
        </div>

        <div className="flex-1 relative">
          <label className="block mb-2 text-sm font-medium text-gray-700">
            Time Range
          </label>
          <button
            onClick={() => setIsPopoverOpen(!isPopoverOpen)}
            className="w-full flex justify-between items-center px-4 py-2.5 border border-gray-200 rounded-lg bg-white text-sm hover:bg-gray-50 transition-all duration-200"
          >
            <span className={timeRange.label ? 'text-gray-800' : 'text-gray-400'}>
              {timeRange.label || 'Select time range'}
            </span>
            <ChevronDown
              className={`w-4 h-4 text-gray-500 transition-transform ${isPopoverOpen ? 'rotate-180' : ''}`}
            />
          </button>

          {isPopoverOpen && (
            <div className="absolute z-20 mt-1 right-0 w-full bg-white border border-gray-200 rounded-lg shadow-xl p-3 space-y-3">
              <div className="space-y-1">
                {timePresets.map(preset => (
                  <button
                    key={preset.value}
                    onClick={() => handlePresetSelect(preset.value)}
                    className={`w-full flex justify-between items-center px-3 py-2 text-sm rounded-md hover:bg-gray-50 transition-colors ${
                      timeRange.value === preset.value
                        ? 'bg-blue-50 text-blue-600'
                        : 'text-gray-700'
                    }`}
                  >
                    {preset.label}
                    {timeRange.value === preset.value && (
                      <Check className="w-4 h-4 text-blue-600" />
                    )}
                  </button>
                ))}
              </div>

              {timeRange.value === 'custom' && (
                <div className="space-y-3 pt-2 border-t border-gray-100">
                  <div className="grid grid-cols-2 gap-3">
                    <div>
                      <label className="block text-xs font-medium text-gray-500 mb-1">
                        Start Date
                      </label>
                      <input
                        type="date"
                        value={customRange.startDate}
                        onChange={e => handleDateChange(e, 'start', 'date')}
                        className="w-full text-sm px-3 py-2 border border-gray-200 rounded-lg"
                      />
                    </div>
                    <div>
                      <label className="block text-xs font-medium text-gray-500 mb-1">
                        Start Time
                      </label>
                      <input
                        type="time"
                        value={customRange.startTime}
                        onChange={e => handleDateChange(e, 'start', 'time')}
                        className="w-full text-sm px-3 py-2 border border-gray-200 rounded-lg"
                      />
                    </div>
                  </div>
                  <div className="grid grid-cols-2 gap-3">
                    <div>
                      <label className="block text-xs font-medium text-gray-500 mb-1">
                        End Date
                      </label>
                      <input
                        type="date"
                        value={customRange.endDate}
                        onChange={e => handleDateChange(e, 'end', 'date')}
                        min={customRange.startDate}
                        className="w-full text-sm px-3 py-2 border border-gray-200 rounded-lg"
                      />
                    </div>
                    <div>
                      <label className="block text-xs font-medium text-gray-500 mb-1">
                        End Time
                      </label>
                      <input
                        type="time"
                        value={customRange.endTime}
                        onChange={e => handleDateChange(e, 'end', 'time')}
                        className="w-full text-sm px-3 py-2 border border-gray-200 rounded-lg"
                      />
                    </div>
                  </div>
                  <button
                    className="w-full mt-1 bg-blue-600 text-white py-2 rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors flex items-center justify-center gap-2"
                    onClick={applyCustomRange}
                  >
                    <Check className="w-4 h-4" />
                    Apply Custom Range
                  </button>
                </div>
              )}
            </div>
          )}
        </div>

        <div className="flex gap-2 items-end">
          <button
            onClick={addFilter}
            className="px-4 py-2.5 bg-slate-900 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center gap-2"
          >
            <Plus className="w-4 h-4" />
            Add
          </button>
          <button
            onClick={applyFilters}
            className="px-4 py-2.5 bg-slate-800 text-white rounded-lg hover:bg-green-700 transition-colors flex items-center gap-2"
          >
            <Check className="w-4 h-4" />
            Apply
          </button>
        </div>
      </div>

      {activeFilters.length > 0 && (
        <div className="space-y-3">
          <h4 className="text-sm font-medium text-gray-700 flex items-center gap-2">
            <span className="bg-blue-100 text-blue-800 w-5 h-5 rounded-full flex items-center justify-center text-xs">
              {activeFilters.length}
            </span>
            Active Filters
          </h4>
          <div className="flex flex-wrap gap-2">
            {activeFilters.map(([key, value]) => (
              <div
                key={key}
                className="bg-blue-50 text-blue-700 px-3 py-1.5 rounded-full text-sm flex items-center gap-2 border border-blue-100"
              >
                <span className="font-medium">
                  {logFilterType[key as keyof typeof logFilterType]}:
                </span>
                <span>{formatTimestamp(key, value)}</span>
                <button
                  className="text-blue-500 hover:text-blue-700 transition-colors"
                  onClick={() => removeFilter(key)}
                >
                  <X size={14} />
                </button>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

export default AttackLogFilter
