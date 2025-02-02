import {Column} from 'react-table'
import Table from '../components/Table'
import LogFilter from '../components/Logs/LogFilter'
import {AttackLog, FilterValues} from '../lib/types'
import {useState} from 'react'
import LogLogs from '../components/Logs/LogLogs'

const LogAttack = () => {
  const columns: Column<AttackLog>[] = [
    {Header: 'IP Address', accessor: 'ipAddress'},
    {Header: 'Application', accessor: 'application'},
    {Header: 'Attack Count', accessor: 'attackCount'},
    {Header: 'Duration', accessor: 'duration'},
    {Header: 'Start At', accessor: 'startAt'},
  ]

  const data: AttackLog[] = [
    {
      ipAddress: '192.168.1.1',
      application: 'https://webserver.example.com',
      attackCount: 12,
      duration: '2h 30m',
      startAt: '2023-10-01 14:00',
    },
    {
      ipAddress: '10.0.0.1',
      application: 'https://database.example.com',
      attackCount: 5,
      duration: '1h 15m',
      startAt: '2023-10-02 09:30',
    },
    {
      ipAddress: '172.16.0.1',
      application: 'https://apigateway.example.com',
      attackCount: 8,
      duration: '3h 45m',
      startAt: '2023-10-03 18:20',
    },
  ]

  const handleFilter = (filters: FilterValues) => {
    console.log('Applied Filters:', filters)
    // Implement your filter logic here (e.g., filter logs based on the criteria)
  }

  const [logType, setLogType] = useState<'log' | 'event'>('log')
  const toggleType = () => setLogType(logType === 'log' ? 'event' : 'log')

  return (
    <div className="px-10 overflow-y-scroll">
      <LogFilter onFilter={handleFilter} logtype={logType} onLogtypeChange={toggleType} />
      {logType === 'event' ? <Table columns={columns} data={data} /> : <LogLogs />}
    </div>
  )
}

export default LogAttack
