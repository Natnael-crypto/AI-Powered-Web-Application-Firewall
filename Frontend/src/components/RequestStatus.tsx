import CustomLineChart from './CustomLineChart'

const data = [
  {name: 'Mon', uv: 4000},
  {name: 'Tue', uv: 3000},
  {name: 'Wed', uv: 2000},
  {name: 'Thu', uv: 2780},
  {name: 'Fri', uv: 1890},
  {name: 'Sat', uv: 2390},
  {name: 'Sun', uv: 3490},
]

const data2 = [
  {name: 'Mon', uv: 2000},
  {name: 'Tue', uv: 5000},
  {name: 'Wed', uv: 4000},
  {name: 'Thu', uv: 1780},
  {name: 'Fri', uv: 4890},
  {name: 'Sat', uv: 3390},
  {name: 'Sun', uv: 1490},
]

const RequestStatus = () => (
  <div className="h-full flex flex-col gap-4 bg-white  shadow-md p-4">
    <h2 className="text-lg font-semibold text-gray-800 px-2">Traffic Analysis</h2>
    <div className="flex-1 flex flex-col gap-4">
      <CustomLineChart data={data} title="Blocking Status" max={4000} />
      <CustomLineChart data={data2} title="Request Status" max={5000} />
    </div>
  </div>
)

export default RequestStatus
