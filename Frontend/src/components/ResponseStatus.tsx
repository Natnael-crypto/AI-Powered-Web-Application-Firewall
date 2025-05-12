import {PieChart, Pie, Cell, ResponsiveContainer, Tooltip} from 'recharts'

const OS_DATA = [
  {name: '404', value: 362, color: '#3366FF'},
  {name: '403', value: 50, color: '#7E88FF'},
  {name: '200', value: 36, color: '#7CDDDD'},
  {name: '467', value: 28, color: '#FF6B81'},
  {name: '400', value: 6, color: '#B6F0E2'},
]

export default function ResponseStatus() {
  return (
    <div className="w-full bg-white xl shadow-md p-6">
      <div className="flex justify-between items-center p-3">
        <p className="text-lg">Response Status Code</p>
      </div>
      <div className="flex  items-center gap-6">
        <div className="w-full sm:w-60 h-60 ">
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

              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </div>

        <div className="grid grid-cols-1  gap-x-32 gap-y-2 text-sm text-slate-700">
          {OS_DATA.map(os => (
            <div key={os.name} className="flex items-center  gap-32">
              <div className="flex items-center gap-2">
                <span
                  className="w-2.5 h-2.5 ull flex-shrink-0 "
                  style={{backgroundColor: os.color}}
                />
                <span className="truncate max-w-[8rem]">{os.name}</span>
              </div>

              <span className="font-semibold ml-auto">{os.value}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
