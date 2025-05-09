import Card from '../components/Card'
import StatisticCard from '../components/StatisticCard'
import StatisticGroup from '../components/StatisticGroup'
import GlobeMap from '../components/GlobeMap'
import UserClientsCard from '../components/Devices-stat'
import ResponseStatus from '../components/ResponseStatus'
import RequestStatus from '../components/RequestStatus'
import TopEndpointsChart from '../components/TopEndpointsChart'
import TopThreatsChart from '../components/TopThreatsChart'
import {useState} from 'react'
import FilterBar from '../components/FilterBar'

const cardStyles =
  'bg-white shadow-lg rounded-xl transition-shadow duration-300 border border-gray-100 hover:shadow-xl'

function Dashboard() {
  const [selectedApp, setSelectedApp] = useState('All')
  const [timeRange, setTimeRange] = useState('24h')

  return (
    <main className="flex flex-col gap-6 w-full">
      {/* Top Statistics */}
      <section className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card className={`flex flex-col justify-center ${cardStyles}`}>
          <StatisticCard
            className="h-full w-full p-4 hover:bg-gray-50 transition-colors duration-300"
            label="Total Requests"
            value="254320"
          />
        </Card>

        {/* Today's Security Summary */}
        <Card className={`col-span-2 items-center ${cardStyles}`}>
          <StatisticGroup
            stats={[
              {label: 'Blocked Requests', value: '1024'},
              {label: 'Malicious IPs Blocked', value: 36},
            ]}
            className="h-full w-full p-4 hover:bg-gray-50 transition-colors duration-300"
          />
        </Card>

        {/* Detection Method Breakdown */}
        <Card className={`col-span-2 items-center ${cardStyles}`}>
          <StatisticGroup
            stats={[
              {label: 'Today Blocked', value: 700},
              {label: 'Today Attack IP', value: 20},
            ]}
            className="h-full w-full p-4 hover:bg-gray-50 transition-colors duration-300"
          />
        </Card>

        {/* Live Traffic Rate */}
        <Card className={`col-span-1 flex flex-col justify-center ${cardStyles}`}>
          <StatisticCard
            className="h-full w-full p-4 hover:bg-gray-50 transition-colors duration-300"
            label="Live Traffic Rate"
            value="122"
          />
        </Card>
      </section>

      {/* Globe Map */}
      <section>
        <Card className={cardStyles}>
          <div className="h-[600px] w-full p-4">
            <GlobeMap />
          </div>
        </Card>
      </section>

      {/* Charts */}
      <section className="">
        <div className="h-[500px] w-full">
          <RequestStatus />
        </div>
      </section>

      {/* Endpoint & Threat Charts */}
      <section className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card className={cardStyles}>
          <TopEndpointsChart />
        </Card>
        <Card className={cardStyles}>
          <TopThreatsChart />
        </Card>
      </section>

      <section className="">
        <div className="grid grid-cols-2 gap-4">
          <div className=" w-full p-4">
            <UserClientsCard />
          </div>
          <div className=" w-full p-4">
            <ResponseStatus />
          </div>
        </div>
      </section>
    </main>
  )
}

export default Dashboard
