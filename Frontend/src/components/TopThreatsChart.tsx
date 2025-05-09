import { PieChart, Pie, Tooltip, Cell, ResponsiveContainer, Legend } from 'recharts';

const COLORS = ['#EF4444', '#F59E0B', '#10B981', '#3B82F6', '#8B5CF6'];

function TopThreatsChart() {
  const mockAttackStats = [
    { threat_type: 'SQL Injection', count: 150 },
    { threat_type: 'XSS', count: 130 },
    { threat_type: 'RCE', count: 110 },
    { threat_type: 'LFI', count: 90 },
    { threat_type: 'CSRF', count: 75 },
  ];

  return (
    <div className="p-4">
      <h2 className="text-lg mb-4">Top 5 Threat Types</h2>
      <ResponsiveContainer width="100%" height={300}>
        <PieChart>
          <Pie
            data={mockAttackStats}
            dataKey="count"
            nameKey="threat_type"
            cx="50%"
            cy="50%"
            outerRadius={90}
            fill="#8884d8"
            label
          >
            {mockAttackStats.map((_, index) => (
              <Cell key={index} fill={COLORS[index % COLORS.length]} />
            ))}
          </Pie>
          <Tooltip />
          <Legend />
        </PieChart>
      </ResponsiveContainer>
    </div>
  );
}

export default TopThreatsChart;
