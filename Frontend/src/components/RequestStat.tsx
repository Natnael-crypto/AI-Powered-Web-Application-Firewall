import {requestData} from '../lib/Constants'
import Card from './Card'
import countries from 'i18n-iso-countries'
import localeRegistration from 'i18n-iso-countries/langs/en.json'

countries.registerLocale(localeRegistration)

function RequestStat({className}: {className?: string}) {
  const formattedData = Object.entries(requestData)
    .map(([key, value]) => ({
      label: countries.getName(key, 'en'),
      value,
    }))
    .sort((a, b) => b.value - a.value) // Sort by value descending

  const maxValue = Math.max(...formattedData.map(item => item.value))

  return (
    <Card
      className={`bg-gradient-to-r from-green-100 to-green-50 p-6 w-full rounded-xl shadow-md ${className}`}
    >
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-xl font-semibold text-gray-800">Request Statistics</h3>
        <span className="text-sm text-gray-500 font-medium">
          {formattedData.length} Countries
        </span>
      </div>

      <div className="flex gap-4 w-full pb-2">
        {formattedData.map((data, index) => {
          const percentage = (data.value / maxValue) * 100
          const computedHeight = `${Math.max(percentage, 15)}%`

          return (
            <div key={index} className="flex flex-col items-center min-w-[100px] group">
              <div className="relative h-40 w-4 bg-gray-200 rounded-full mb-2">
                <div
                  className="absolute bottom-0 w-full bg-green-600 rounded-full transition-all duration-300 group-hover:bg-green-500"
                  style={{height: computedHeight}}
                />
              </div>

              <div className="flex flex-col items-center">
                <span className="font-medium text-sm text-gray-800 mb-1">
                  {data.value.toLocaleString()}
                </span>
                <span
                  className="text-xs text-gray-600 text-center whitespace-nowrap max-w-[120px] truncate"
                  title={data.label}
                >
                  {data.label}
                </span>
              </div>
            </div>
          )
        })}
      </div>

      <div className="mt-4 flex justify-between text-xs text-gray-500">
        <span>0</span>
        <span>{maxValue.toLocaleString()}</span>
      </div>
    </Card>
  )
}

export default RequestStat
