import TableWithFilters from '../components/TableWithFilter'

const LogAttack = () => {
  const columns = [
    {key: 'ipAddress', header: 'IP Address'},
    {key: 'application', header: 'Application'},
    {key: 'attackCount', header: 'Attack Count'},
    {key: 'duration', header: 'Duration'},
    {key: 'startAt', header: 'Start At'},
  ]

  const data = [
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
    {
      ipAddress: '192.168.1.2',
      application: 'https://webserver.example.com',
      attackCount: 20,
      duration: '4h 10m',
      startAt: '2023-10-04 12:00',
    },
    {
      ipAddress: '192.168.1.3',
      application: 'https://database.example.com',
      attackCount: 7,
      duration: '1h 45m',
      startAt: '2023-10-05 08:15',
    },
    {
      ipAddress: '10.0.0.2',
      application: 'https://apigateway.example.com',
      attackCount: 15,
      duration: '2h 50m',
      startAt: '2023-10-06 16:40',
    },
    {
      ipAddress: '172.16.0.2',
      application: 'https://webserver.example.com',
      attackCount: 10,
      duration: '3h 20m',
      startAt: '2023-10-07 11:00',
    },
    {
      ipAddress: '192.168.1.4',
      application: 'https://database.example.com',
      attackCount: 3,
      duration: '0h 50m',
      startAt: '2023-10-08 14:30',
    },
    {
      ipAddress: '10.0.0.3',
      application: 'https://apigateway.example.com',
      attackCount: 18,
      duration: '5h 10m',
      startAt: '2023-10-09 19:00',
    },
    {
      ipAddress: '172.16.0.3',
      application: 'https://webserver.example.com',
      attackCount: 9,
      duration: '2h 15m',
      startAt: '2023-10-10 10:45',
    },
    {
      ipAddress: '192.168.1.5',
      application: 'https://database.example.com',
      attackCount: 6,
      duration: '1h 30m',
      startAt: '2023-10-11 13:20',
    },
    {
      ipAddress: '10.0.0.4',
      application: 'https://apigateway.example.com',
      attackCount: 14,
      duration: '3h 55m',
      startAt: '2023-10-12 17:10',
    },
    {
      ipAddress: '172.16.0.4',
      application: 'https://webserver.example.com',
      attackCount: 11,
      duration: '2h 40m',
      startAt: '2023-10-13 12:05',
    },
    {
      ipAddress: '192.168.1.6',
      application: 'https://database.example.com',
      attackCount: 4,
      duration: '1h 10m',
      startAt: '2023-10-14 09:50',
    },
    {
      ipAddress: '10.0.0.5',
      application: 'https://apigateway.example.com',
      attackCount: 16,
      duration: '4h 20m',
      startAt: '2023-10-15 20:30',
    },
    {
      ipAddress: '172.16.0.5',
      application: 'https://webserver.example.com',
      attackCount: 8,
      duration: '2h 05m',
      startAt: '2023-10-16 11:15',
    },
    {
      ipAddress: '192.168.1.7',
      application: 'https://database.example.com',
      attackCount: 5,
      duration: '1h 25m',
      startAt: '2023-10-17 14:40',
    },
    {
      ipAddress: '10.0.0.6',
      application: 'https://apigateway.example.com',
      attackCount: 19,
      duration: '5h 30m',
      startAt: '2023-10-18 18:00',
    },
    {
      ipAddress: '172.16.0.6',
      application: 'https://webserver.example.com',
      attackCount: 7,
      duration: '1h 50m',
      startAt: '2023-10-19 10:20',
    },
    {
      ipAddress: '192.168.1.8',
      application: 'https://database.example.com',
      attackCount: 10,
      duration: '2h 15m',
      startAt: '2023-10-20 13:00',
    },
    {
      ipAddress: '10.0.0.7',
      application: 'https://apigateway.example.com',
      attackCount: 13,
      duration: '3h 10m',
      startAt: '2023-10-21 16:45',
    },
    {
      ipAddress: '172.16.0.7',
      application: 'https://webserver.example.com',
      attackCount: 6,
      duration: '1h 35m',
      startAt: '2023-10-22 09:10',
    },
    {
      ipAddress: '192.168.1.9',
      application: 'https://database.example.com',
      attackCount: 8,
      duration: '2h 00m',
      startAt: '2023-10-23 12:30',
    },
    {
      ipAddress: '10.0.0.8',
      application: 'https://apigateway.example.com',
      attackCount: 17,
      duration: '4h 45m',
      startAt: '2023-10-24 19:20',
    },
    {
      ipAddress: '172.16.0.8',
      application: 'https://webserver.example.com',
      attackCount: 9,
      duration: '2h 25m',
      startAt: '2023-10-25 11:50',
    },
    {
      ipAddress: '192.168.1.10',
      application: 'https://database.example.com',
      attackCount: 4,
      duration: '1h 15m',
      startAt: '2023-10-26 14:15',
    },
    {
      ipAddress: '10.0.0.9',
      application: 'https://apigateway.example.com',
      attackCount: 20,
      duration: '5h 00m',
      startAt: '2023-10-27 17:30',
    },
    {
      ipAddress: '172.16.0.9',
      application: 'https://webserver.example.com',
      attackCount: 11,
      duration: '2h 50m',
      startAt: '2023-10-28 10:05',
    },
    {
      ipAddress: '192.168.1.11',
      application: 'https://database.example.com',
      attackCount: 6,
      duration: '1h 40m',
      startAt: '2023-10-29 13:25',
    },
    {
      ipAddress: '10.0.0.10',
      application: 'https://apigateway.example.com',
      attackCount: 18,
      duration: '4h 10m',
      startAt: '2023-10-30 20:00',
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
    <div className="px-10 overflow-y-scroll">
      <TableWithFilters columns={columns} data={data} filterConfig={filterConfig} />
    </div>
  )
}

export default LogAttack
