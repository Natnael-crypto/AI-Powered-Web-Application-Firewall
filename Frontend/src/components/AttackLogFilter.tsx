import {ChangeEvent} from 'react'
import {filterOperations, logFilterType} from '../lib/types'
import {useLogFilter} from '../store/LogFilter'

function AttackLogFilter() {
  const {filterType, filterOperation, selectFilterOperation, selectFilterType} =
    useLogFilter()

  const handleSetFilterOperation = (e: ChangeEvent<HTMLSelectElement>) => {
    const input = e.target.value
    const matched = Object.values(filterOperations).find(val => val === input)
    selectFilterOperation(matched as filterOperations)
  }

  const handleSetfilterType = (e: ChangeEvent<HTMLSelectElement>) => {
    const input = e.target.value
    const matched = Object.values(logFilterType).find(val => val === input)
    selectFilterType(matched as logFilterType)
  }

  const filterTypeOption = Object.values(logFilterType).filter(
    val => typeof val !== 'number',
  )

  const filterOperationOption = Object.values(filterOperations).filter(
    val => typeof val !== 'number',
  )

  return (
    <div className="w-full bg-white p-6   flex flex-col gap-4 md:flex-row md:items-center">
      {/* Filter Type */}
      <div className="flex-1">
        <label className="block mb-1 text-sm font-medium text-gray-700">
          Filter Type
        </label>
        <select
          value={filterType ?? ''}
          onChange={handleSetfilterType}
          className="w-full px-4 py-2  border border-gray-300 text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="">Select Type</option>
          {filterTypeOption.map((val, index) => (
            <option key={index} value={val}>
              {val}
            </option>
          ))}
        </select>
      </div>

      {/* Filter Operation */}
      <div className="flex-1">
        <label className="block mb-1 text-sm font-medium text-gray-700">Operation</label>
        <select
          value={filterOperation ?? ''}
          onChange={handleSetFilterOperation}
          className="w-full px-4 py-2  border border-gray-300 text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="">Select Operation</option>
          {filterOperationOption.map((op, index) => (
            <option key={index} value={op}>
              {op}
            </option>
          ))}
        </select>
      </div>

      {/* Apply Button */}
      <div className="pt-1 md:pt-5">
        <button className="w-full md:w-auto px-6 py-2  bg-blue-600 text-white font-semibold hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500">
          Apply Filter
        </button>
      </div>
    </div>
  )
}

export default AttackLogFilter
