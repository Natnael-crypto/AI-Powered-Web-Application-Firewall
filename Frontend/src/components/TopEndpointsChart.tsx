import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid } from 'recharts';

function TopEndpointsChart() {
  const mockEndpointStats = [
    { application_name: 'App A', request_url: '/login', count: 120 },
    { application_name: 'App B', request_url: '/admin', count: 95 },
    { application_name: 'App A', request_url: '/api/user', count: 85 },
    { application_name: 'App C', request_url: '/upload', count: 70 },
    { application_name: 'App D', request_url: '/config', count: 60 },
  ];

  return (
    <div className="p-4">
      <h2 className="text-lg mb-4">Top 5 Targeted Endpoints</h2>
      <ResponsiveContainer width="100%" height={300}>
        <BarChart data={mockEndpointStats}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="request_url" />
          <YAxis />
          <Tooltip />
          <Bar dataKey="count" fill="#3B82F6" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}

export default TopEndpointsChart;
