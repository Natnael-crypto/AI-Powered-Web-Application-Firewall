import {PieChart, Pie, Cell, ResponsiveContainer, Tooltip} from 'recharts'

const OS_DATA = [
  {name: 'Windows', value: 362, color: '#3366FF'},
  {name: 'Linux', value: 50, color: '#7E88FF'},
  {name: 'Android', value: 36, color: '#7CDDDD'},
  {name: 'macOS', value: 28, color: '#FF6B81'},
  {name: 'iOS', value: 6, color: '#B6F0E2'},
]

const BROWSER_DATA = [
  {name: 'Chrome', value: 342, color: '#B6F0E2'},
  {name: 'Custom-AsyncHttpClient', value: 86, color: '#FF6B81'},
  {name: 'Firefox', value: 39, color: '#FFD166'},
  {name: 'Go-http-client', value: 35, color: '#7CDDDD'},
  {name: 'Headless Chrome', value: 28, color: '#7E88FF'},
]

export default function UserClientsCard() {
  return (
    <div className="w-full bg-white xl shadow-md p-6">
      <div className="flex justify-between items-center p-3">
        <p className="text-lg">User clients</p>
        <button className="text-xs text-blue-600 font-semibold hover:text-blue-700">
          MORE
        </button>
      </div>
      <div className="flex items-center gap-6">
        {/* Chart */}
        <div className="w-full sm:w-60 h-60">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie
                data={OS_DATA}
                dataKey="value"
                outerRadius={70}
                innerRadius={45}
                startAngle={90}
                endAngle={-270}
              >
                {OS_DATA.map((entry, index) => (
                  <Cell key={`os-${index}`} fill={entry.color} />
                ))}
              </Pie>
              <Pie
                data={BROWSER_DATA}
                dataKey="value"
                outerRadius={90}
                innerRadius={75}
                startAngle={90}
                endAngle={-270}
              >
                {BROWSER_DATA.map((entry, index) => (
                  <Cell key={`browser-${index}`} fill={entry.color} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </div>

        {/* Legend */}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-x-32 gap-y-2 text-sm text-slate-700">
          {[...OS_DATA, ...BROWSER_DATA].map(({name, value, color}) => (
            <div key={name} className="flex items-center gap-2">
              <span
                className="w-2.5 h-2.5 ull flex-shrink-0"
                style={{backgroundColor: color}}
              />
              <span className="truncate max-w-[8rem]">{name}</span>
              <span className="font-semibold ml-auto">{value}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
