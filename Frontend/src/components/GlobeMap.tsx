import { useEffect, useState, useRef } from 'react'
import { VectorMap } from '@react-jvectormap/core'
import { worldMill } from '@react-jvectormap/world'
import { useMapState } from '../hooks/api/useDashboardStat'
import { whereCountry } from 'iso-3166-1'

interface GlobeMapProps {
  selectedApp: string
  timeRange: any
}

const GlobeMap = ({ selectedApp, timeRange }: GlobeMapProps) => {
  const [filter, setFilter] = useState<'all' | 'blocked'>('all')
  const threat = filter === 'all' ? '' : 'blocked'
  const { data, refetch } = useMapState(selectedApp, timeRange, threat)
  const [country, setCountry] = useState<any>({})
  const countryRef = useRef<any>({})

  useEffect(() => {
    refetch()
  }, [filter])

  useEffect(() => {
    const updatedCountryData: any = {}
    if (data) {
      Object.keys(data).forEach((countryName) => {
        const countryCode = whereCountry(countryName)
        if (countryCode) {
          updatedCountryData[countryCode.alpha2] = data[countryName]
        }
      })
    }
    setCountry(updatedCountryData)
  }, [data])

  useEffect(() => {
    countryRef.current = country
  }, [country])

  return (
    <div className="h-full flex flex-col">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-lg text-gray-800">Global Traffic Map</h2>
        <div className="flex gap-2 items-center">
          <select
            value={filter}
            onChange={e => setFilter(e.target.value as 'all' | 'blocked')}
            className="text-xs font-medium text-gray-700 bg-gray-100 px-2 py-1 rounded border border-gray-300 focus:outline-none"
          >
            <option value="all">All Requests</option>
            <option value="blocked">Blocked Requests</option>
          </select>
          <div className="text-xs font-medium text-gray-500 bg-gray-100 px-3 py-1 rounded-full">
            Live Data
          </div>
        </div>
      </div>

      <div className="flex-1 rounded-lg overflow-hidden">
        <VectorMap
          backgroundColor="#f8fafc"
          className="h-full w-full"
          zoomOnScroll={false}
          regionStyle={{
            initial: { fill: '#E2E8F0', stroke: '#fff', strokeWidth: 1 },
            hover: { fill: '#3B82F6', cursor: 'pointer' },
          }}
          series={{
            regions: [
              {
                values: country,
                scale: ['#E2E8F0', '#1E40AF'],
                normalizeFunction: 'linear',
                attribute: 'fill',
              },
            ],
          }}
          onRegionTipShow={(_, el: any, code) => {
            const requests = countryRef.current[code] || 0
            el.html(
              `<div style="font-family: 'Oxygen', sans-serif; padding: 12px 16px; background-color:rgb(23, 36, 43); color: #fff; border-radius: 8px; box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); max-width: 200px;">
                <div style="font-size: 14px; font-weight: 600; color: #fff; margin-bottom: 8px;">
                  <span style="color: #BBE1F7;">Country: </span>${el.html()}
                </div>
                <div style="font-size: 13px; color: #B0C4DE;">
                  <strong>Requests:</strong> <span style="color: #4CAF50;">${requests.toLocaleString()}</span>
                </div>
              </div>`
            )
          }}
          map={worldMill}
        />
      </div>
    </div>
  )
}

export default GlobeMap
