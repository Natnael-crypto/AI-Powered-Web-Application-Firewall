import {PieChart, Pie, Tooltip, Cell, ResponsiveContainer, Legend} from 'recharts'
import {useGetTopThreatTypes} from '../hooks/api/useDashboardStat'

const COLORS = ['#EF4444', '#F59E0B', '#10B981', '#3B82F6', '#8B5CF6']

function TopThreatsChart() {
  const {data, isLoading, isError} = useGetTopThreatTypes()

  if (isLoading) return <p>Loading ...</p>
  if (isError) return <p>Something went wrong</p>

  const chartData = data.map((item, index) => ({
    ...item,
    short_label: `#${index + 1}: ${item.threat_type.slice(0, 40)}${item.threat_type.length > 40 ? '...' : ''}`,
  }))

  return (
    <div className="p-4">
      <h2 className="text-lg mb-4">Top 5 Threat Types</h2>
      <ResponsiveContainer width="100%" height={300}>
        <PieChart>
          <Pie
            data={chartData}
            dataKey="count"
            nameKey="short_label"
            cx="50%"
            cy="50%"
            outerRadius={90}
            fill="#8884d8"
            label
          >
            {chartData.map((_, index) => (
              <Cell key={index} fill={COLORS[index % COLORS.length]} />
            ))}
          </Pie>
          <Tooltip
            formatter={(value, name, props) => [value, props.payload.threat_type]}
          />
          <Legend />
        </PieChart>
      </ResponsiveContainer>
    </div>
  )
}

export default TopThreatsChart
