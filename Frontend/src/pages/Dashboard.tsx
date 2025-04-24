import Card from '../components/Card'
import StatisticCard from '../components/StatisticCard'
import StatisticGroup from '../components/StatisticGroup'
import GlobeMap from '../components/GlobeMap'
import UserClientsCard from '../components/Devices-stat'
import ResponseStatus from '../components/ResponseStatus'
import RequestStatus from '../components/RequestStatus'

function Dashboard() {
  return (
    <div className="px-5 h-screen flex flex-col gap-5 w-full overflow-y-scroll">
      <div className="flex gap-5">
        <div className="flex flex-col gap-5 w-full">
          <div className="flex gap-3">
            <Card className=" flex bg-white flex-col justify-center ">
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
        </div>
      </div>
      <div className="flex min-h-[550px] gap-8 ">
        <Card className="bg-white">
          <GlobeMap />
        </Card>
        <RequestStatus />
      </div>

      <div className="flex w-full gap-5 bg-white">
        <UserClientsCard />
        <ResponseStatus />
      </div>
    </div>
  )
}
;[]

export default Dashboard
