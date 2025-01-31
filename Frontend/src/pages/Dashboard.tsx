import Card from '../components/Card'
import DashboardMenu from '../components/DashboardMenu'
import StatisticCard from '../components/StatisticCard'
import StatisticGroup from '../components/StatisticGroup'
import GlobeMap from '../components/GlobeMap'
import {useState} from 'react'
import TodaySummary from '../components/TodaySummary'

function Dashboard() {
  const [selectedDashboardMenu, setSelectedDashboardMenu] = useState<
    'basic' | 'advanced'
  >('basic')

  const changeDashboardMenu = () =>
    setSelectedDashboardMenu(selectedDashboardMenu == 'basic' ? 'advanced' : 'basic')
  const stats = [
    {label: 'Requests', value: 700},
    {label: 'Requests', value: 700},
    {label: 'Requests', value: 700},
    {label: 'Requests', value: 700},
  ]

  const requestData = [
    {lat: 37.7749, lng: -122.4194, intensity: 100}, // San Francisco
    {lat: 40.7128, lng: -74.006, intensity: 200}, // New York
    {lat: 51.5074, lng: -0.1278, intensity: 150}, // London
    {lat: 35.6895, lng: 139.6917, intensity: 300}, // Tokyo
  ]

  return (
    <div className="px-5 flex flex-col gap-5 w-full overflow-y-scroll">
      <DashboardMenu
        changeMenu={changeDashboardMenu}
        selectedMenu={selectedDashboardMenu}
      />
      {selectedDashboardMenu == 'basic' && (
        <div className="flex gap-5">
          <div className="flex flex-col gap-5 w-full">
            <div className="flex gap-3">
              <Card className="gap-y-5 flex bg-white flex-col justify-center pl-5 py-5">
                <StatisticCard
                  className="h-full w-full py-3"
                  label="Requests"
                  value={700}
                />
              </Card>
              <Card className=" items-center bg-white">
                <StatisticGroup
                  stats={[
                    {label: 'Requests', value: 700},
                    {label: 'Requests', value: 700},
                  ]}
                  className=" h-full w-full"
                />
              </Card>
            </div>
            <Card className="bg-white">
              <StatisticGroup className="py-5" stats={stats} />
            </Card>
          </div>
          <Card className="h-full w-[55%] shadow-lg bg-white">Query per second</Card>
        </div>
      )}
      <div className="flex gap-8">
        <Card className="bg-white">
          <GlobeMap data={requestData} />
        </Card>
        <TodaySummary />
      </div>
    </div>
  )
}

export default Dashboard
