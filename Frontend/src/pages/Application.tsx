import TableWithFilters from '../components/TableWithFilter'

function Application() {
  const columns = [
    {key: 'application', header: 'Application'},
    {key: 'port', header: 'Port'},
    {key: 'runMode', header: 'Run Mode'},
    {key: 'security', header: 'Security'},
    {key: 'today', header: 'Today'},
  ]

  const data = [
    {
      ipAddress: '192.168.1.1',
      application: 'Web Server',
      attackCount: 12,
      duration: '2h 30m',
      startAt: '2023-10-01 14:00',
    },
    {
      ipAddress: '10.0.0.1',
      application: 'Database',
      attackCount: 5,
      duration: '1h 15m',
      startAt: '2023-10-02 09:30',
    },
    {
      ipAddress: '172.16.0.1',
      application: 'API Gateway',
      attackCount: 8,
      duration: '3h 45m',
      startAt: '2023-10-03 18:20',
    },
    {
      ipAddress: '192.168.1.2',
      application: 'Web Server',
      attackCount: 20,
      duration: '4h 10m',
      startAt: '2023-10-04 12:00',
    },
  ]

  const filterConfig: {type: string; name: string; placeholder: string}[] = [
    {type: 'text', name: 'ipAddress', placeholder: 'IP Address'},
    {type: 'text', name: 'application', placeholder: 'Application'},
    {type: 'text', name: 'port', placeholder: 'Port'},
    {type: 'date', name: 'startDate', placeholder: 'Start Date'},
    {type: 'date', name: 'endDate', placeholder: 'End Date'},
  ]

  return (
    <div className="px-10  py-2">
      <TableWithFilters columns={columns} data={data} filterConfig={filterConfig} />
    </div>
  )
}

export default Application
