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
  // const stats = [
  //   {label: 'Requests', value: 2400},
  //   {label: 'Today Blocked', value: 700},
  //   {label: 'Today Attack IP', value: 700},
  //   {label: 'Requests', value: 700},
  // ]

  const requestData = [
    {lat: 37.7749, lng: -122.4194, intensity: 100}, // San Francisco, USA
    {lat: 48.8566, lng: 2.3522, intensity: 200}, // Paris, France
    {lat: 51.5074, lng: -0.1278, intensity: 150}, // London, UK
    {lat: 35.6895, lng: 139.6917, intensity: 300}, // Tokyo, Japan
    {lat: -33.8688, lng: 151.2093, intensity: 250}, // Sydney, Australia
    {lat: 55.7558, lng: 37.6173, intensity: 180}, // Moscow, Russia
    {lat: -23.5505, lng: -46.6333, intensity: 220}, // São Paulo, Brazil
  ]

  return (
    <div className="px-5 h-screen flex flex-col gap-5 w-full overflow-y-scroll">
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
                    {label: 'Today Blocked', value: 700},
                    {label: 'Today Attack IP', value: 20},
                  ]}
                  className=" h-full w-full"
                />
              </Card>
            </div>
            {/* <Card className="bg-white">
              <StatisticGroup className="py-5" stats={stats} />
            </Card> */}
          </div>
          <Card className="h-full w-[55%] shadow-lg bg-white">Query per second</Card>
        </div>
      )}
      <div className="flex min-h-[500px] gap-8">
        <Card className="bg-white">
          <GlobeMap />
        </Card>
        <TodaySummary />
      </div>
    </div>
  )
}

export default Dashboard
