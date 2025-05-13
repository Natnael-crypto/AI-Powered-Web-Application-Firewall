import { useEffect, useState } from 'react'
import Card from '../components/Card'
import StatisticCard from '../components/StatisticCard'
import StatisticGroup from '../components/StatisticGroup'
import GlobeMap from '../components/GlobeMap'
import UserClientsCard from '../components/Devices-stat'
import ResponseStatus from '../components/ResponseStatus'
// import RequestStatus from '../components/RequestStatus'
import TopEndpointsChart from '../components/TopEndpointsChart'
import TopThreatsChart from '../components/TopThreatsChart'
import FilterBar from '../components/FilterBar'
import { useGetOverAllStat, useRateStat } from '../hooks/api/useDashboardStat'

const cardStyles =
  'bg-white shadow-md rounded-md transition-shadow duration-300 border border-gray-100 hover:shadow-lg'

function Dashboard() {
  const [selectedApp, setSelectedApp] = useState('')
  const [timeRange, setTimeRange] = useState({
    value: '',
    label: '',
    start: '',
    end: '',
  })

  const { data, refetch } = useGetOverAllStat(selectedApp, timeRange)
  const { data: rateData, refetch: refetchRate } = useRateStat(selectedApp, timeRange)

  const [liveRate, setLiveRate] = useState(rateData)

  useEffect(() => {
    refetch()
  }, [selectedApp, timeRange])

  useEffect(() => {
    setLiveRate(rateData)
  }, [rateData])

  useEffect(() => {
    const interval = setInterval(() => {
      refetchRate()
    }, 60000) // every 60 seconds

    return () => clearInterval(interval)
  }, [selectedApp, timeRange])

  return (
    <main className="flex flex-col gap-6 w-full">
      <FilterBar
        selectedApp={selectedApp}
        setSelectedApp={setSelectedApp}
        timeRange={timeRange}
        setTimeRange={setTimeRange}
      />

      {/* Top Security Overview */}
      <section className="grid grid-cols-6 gap-7">
        <div className="col-span-1 flex flex-col justify-center rounded-lg border bg-white shadow-sm">
          <div className="h-full w-full p-4 hover:bg-gray-50 transition-colors duration-300">
            <StatisticCard label="Total Requests" value={data?.total_requests} />
          </div>
        </div>

        <div className="col-span-2 rounded-lg border bg-white shadow-sm flex items-center">
          <div className="h-full w-full p-4 hover:bg-gray-50 transition-colors duration-300">
            <StatisticGroup
              stats={[
                { label: 'Blocked Requests', value: data?.blocked_requests },
                { label: 'Malicious IPs Blocked', value: data?.malicious_ips_blocked },
              ]}
            />
          </div>
        </div>

        <div className="col-span-2 rounded-lg border bg-white shadow-sm flex items-center">
          <div className="h-full w-full p-4 hover:bg-gray-50 transition-colors duration-300">
            <StatisticGroup
              stats={[
                { label: 'AI-Based Detections', value: data?.ai_based_detections },
                { label: 'Rule-Based Detections', value: data?.rule_based_detections },
              ]}
            />
          </div>
        </div>

        <div className="col-span-1 flex flex-col justify-center rounded-lg border bg-white shadow-sm">
          <div className="h-full w-full p-4 hover:bg-gray-50 transition-colors duration-300">
            <StatisticCard label="Live Traffic Rate" value={liveRate} />
          </div>
        </div>
      </section>

      <section>
        <Card className={cardStyles}>
          <div className="h-[600px] w-full p-4">
            <GlobeMap selectedApp={selectedApp} timeRange={timeRange} />
          </div>
        </Card>
      </section>

      <section className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="rounded-lg border border-gray-200 bg-white shadow-sm">
          <TopEndpointsChart selectedApp={selectedApp} timeRange={timeRange} />
        </div>
        <div className="rounded-lg border border-gray-200 bg-white shadow-sm">
          <TopThreatsChart selectedApp={selectedApp} timeRange={timeRange} />
        </div>
      </section>

      <section className="">
        <div className="grid grid-cols-2 gap-4">
          <div className="w-full">
            <UserClientsCard selectedApp={selectedApp} timeRange={timeRange} />
          </div>
          <div className="w-full">
            <ResponseStatus selectedApp={selectedApp} timeRange={timeRange} />
          </div>
        </div>
      </section>
    </main>
  )
}

export default Dashboard
