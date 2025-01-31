import {useState} from 'react'
import Filter from './Filter'
import Table from './Table'

interface TableWithFiltersProps {
  columns: {key: string; header: string}[]
  data: Record<string, any>[]
  filterConfig: {type: string; name: string; placeholder: string}[]
}

const TableWithFilters = ({columns, data, filterConfig}: TableWithFiltersProps) => {
  const [filters, setFilters] = useState<Record<string, string>>({})

  const handleFilterChange = (updatedFilters: Record<string, string>) => {
    setFilters(updatedFilters)
  }

  const filteredData = data.filter(row => {
    return Object.entries(filters).every(([key, value]) => {
      if (!value) return true
      if (key === 'startDate' || key === 'endDate') {
        return row.startAt >= filters.startDate && row.startAt <= filters.endDate
      }
      return String(row[key]).toLowerCase().includes(value.toLowerCase())
    })
  })

  return (
    <div className="bg-gray-50 overflow-y-scroll shadow-lg rounded-lg h-full">
      <Filter filterConfig={filterConfig} onFilterChange={handleFilterChange} />
      <Table columns={columns} data={filteredData} />
    </div>
  )
}

export default TableWithFilters
