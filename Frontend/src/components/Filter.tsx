import React, {useState} from 'react'

interface FilterConfig {
  type: string
  name: string
  placeholder: string
}

interface FilterProps {
  filterConfig: FilterConfig[]
  onFilterChange: (filters: Record<string, string>) => void
}

const Filter = ({filterConfig, onFilterChange}: FilterProps) => {
  const [filters, setFilters] = useState<Record<string, string>>({})

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const {name, value} = e.target
    const updatedFilters = {...filters, [name]: value}
    setFilters(updatedFilters)
    onFilterChange(updatedFilters)
  }

  return (
    <div className="py-6 rounded-lg mb-2 sticky top-0 bg-white px-6">
      <h2 className="text-xl font-semibold mb-4 text-gray-800">Filters</h2>
      <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-4">
        {filterConfig.map(filter => (
          <input
            key={filter.name}
            type={filter.type}
            name={filter.name}
            placeholder={filter.placeholder}
            value={filters[filter.name] || ''}
            onChange={handleInputChange}
            className="p-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        ))}
      </div>
    </div>
  )
}

export default Filter
