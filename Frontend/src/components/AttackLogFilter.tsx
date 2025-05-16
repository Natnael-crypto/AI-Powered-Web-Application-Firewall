import { ChangeEvent } from 'react'
import { logFilterType } from '../lib/types'
import { useLogFilter } from '../store/LogFilter'
import { X } from 'lucide-react'

function AttackLogFilter() {
  const {
    filters,
    tempFilter,
    setTempFilter,
    addFilter,
    removeFilter,
    applyFilters,
  } = useLogFilter()

  const handleTypeChange = (e: ChangeEvent<HTMLSelectElement>) => {
    setTempFilter('key', e.target.value)
  }

  const handleValueChange = (e: ChangeEvent<HTMLInputElement>) => {
    setTempFilter('value', e.target.value)
  }

  const filterKeys = Object.keys(logFilterType)

  const activeFilters = Object.entries(filters).filter(
    ([_, value]) => value !== ''
  )

  return (
    <div className="w-full bg-white p-6 rounded-lg shadow flex flex-col gap-6">
      <div className="flex flex-col md:flex-row gap-4">
        <div className="flex-1">
          <label className="block mb-1 text-sm font-semibold text-gray-700">
            Filter Type
          </label>
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
          <label className="block mb-1 text-sm font-semibold text-gray-700">
            Value
          </label>
          <input
            type="text"
            value={tempFilter.value}
            onChange={handleValueChange}
            className="w-full px-4 py-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
          />
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
                {logFilterType[key as keyof typeof logFilterType]}: {value}
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
