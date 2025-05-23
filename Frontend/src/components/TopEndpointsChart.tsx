import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  CartesianGrid,
} from 'recharts'
import { useGetMostTargetedEndpoint } from '../hooks/api/useDashboardStat'
import { useEffect } from 'react'

interface TopEndpointsChartProps {
  selectedApp: string
  timeRange: any
}

function TopEndpointsChart({ selectedApp, timeRange }: TopEndpointsChartProps) {
  const { data, isLoading, isError,refetch } = useGetMostTargetedEndpoint(selectedApp, timeRange)

  useEffect(()=>{
    refetch()
  },[selectedApp,timeRange])

  if (isLoading) return <p>Loading...</p>
  if (isError) return <p>Something went wrong</p>

  return (
    <div className="p-4">
      <h2 className="text-lg mb-4">Top 5 Targeted Endpoints</h2>
      <ResponsiveContainer width="100%" height={300}>
        <BarChart data={data}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis
            dataKey="request_url"
            tickFormatter={(value: string) =>
              value.length > 10 ? `${value.slice(0, 10)}...` : value
            }
          />
          <YAxis />
          <Tooltip
            formatter={(value: any, name: any, props: any) => [value, 'Requests']}
            labelFormatter={(label: string) => `URL: ${label}`}
          />
          <Bar dataKey="count" fill="#3B82F6" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  )
}

export default TopEndpointsChart
