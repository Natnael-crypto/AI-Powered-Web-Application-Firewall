import {useState} from 'react'
import Button from './atoms/Button'

const CountryData = () => {
  const [isBlocked, setIsBlocked] = useState(true)
  const toggleIsBlocked = () => setIsBlocked(!isBlocked)
  const data = [
    {country: 'China', value: '2.0m'},
    {country: 'Singapore', value: '8.4k'},
    {country: 'India', value: '5.3k'},
    {country: 'United States', value: '3.6k'},
    {country: 'Australia', value: '3.3k'},
    {country: 'Japan', value: '2.8k'},
    {country: 'United Kingdom', value: '1.8k'},
  ]

  return (
    <div className="px-5 py-3  w-max space-y-10">
      <div className="border border-1 p-1 border-gray-300 flex self-end">
        <Button
          size="s"
          classname={`w-1/2 text-gray-500 text-md font-thin ${isBlocked ? 'bg-green-500 text-white' : ''}`}
          onClick={!isBlocked ? toggleIsBlocked : () => {}}
        >
          Requests
        </Button>
        <Button
          size="s"
          classname={`w-1/2 text-gray-500 text-md font-thin ${!isBlocked ? 'bg-green-500 text-white' : ''}`}
          onClick={isBlocked ? toggleIsBlocked : () => {}}
        >
          Blocked
        </Button>
      </div>
      <h2 className="text-2xl font-bold mb-4 text-gray-800">Country Data</h2>
      <ul className="space-y-2">
        {data.map((item, index) => (
          <li key={index} className="p-3 space-y-2    hover:shadow-md transition-shadow">
            <div className="flex justify-between items-center gap-16">
              <span className="text-gray-700 font-medium">{item.country}</span>
              <span className="text-gray-600">{item.value}</span>
            </div>
            <div className="h-1 w-full bg-gray-300 relative">
              <div className="absolute h-1 w-[50%] bg-blue-600" />
            </div>
          </li>
        ))}
      </ul>
    </div>
  )
}

export default CountryData
