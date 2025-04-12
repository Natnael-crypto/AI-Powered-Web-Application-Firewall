import {requestData} from '../lib/Constants'
import Card from './Card'
import countries from 'i18n-iso-countries'
import localeRegistration from 'i18n-iso-countries/langs/en.json'

countries.registerLocale(localeRegistration)

function RequestStat({className}: {className?: string}) {
  const formattedData = Object.entries(requestData).map(([key, value]) => ({
    label: countries.getName(key, 'en'),
    value,
  }))
  const maxValue = Math.max(...formattedData.map(item => item.value)) // Get max value

  return (
    <Card
      className={`flex flex-col justify-center w-1/4 items-center px-5 py-3 ${className}`}
    >
      {formattedData.map((data, index) => {
        const computedWidth = `${(data.value / maxValue) * 100}%` // Calculate relative width

        return (
          <div
            key={index}
            className="flex flex-col w-full p-2 hover:shadow-lg rounded-lg gap-3"
          >
            <span>{data.label}</span>
            <div className="h-1 bg-green-600" style={{width: computedWidth}}></div>
          </div>
        )
      })}
    </Card>
  )
}

export default RequestStat
