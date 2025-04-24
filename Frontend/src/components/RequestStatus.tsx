import CustomLineChart from './CustomLineChart'

const data = [
  {name: 'Page A', uv: 4000},
  {name: 'Page B', uv: 3000},
  {name: 'Page C', uv: 2000},
  {name: 'Page D', uv: 2780},
  {name: 'Page E', uv: 1890},
  {name: 'Page F', uv: 2390},
  {name: 'Page G', uv: 3490},
]

const data2 = [
  {name: 'Page A', uv: 2000},
  {name: 'Page B', uv: 5000},
  {name: 'Page C', uv: 4000},
  {name: 'Page D', uv: 1780},
  {name: 'Page E', uv: 4890},
  {name: 'Page F', uv: 3390},
  {name: 'Page G', uv: 1490},
]
const RequestStatus = () => (
  <>
    <div className="w-2/6 h-full flex flex-col  gap-5 bg-gray-50 p-3">
      <CustomLineChart data={data} title="Blocking Status" max={4000} />
      <CustomLineChart data={data2} title="Request Status" max={5000} />
    </div>
  </>
)

export default RequestStatus
