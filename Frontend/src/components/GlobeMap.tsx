import React from 'react'
import {ComposableMap, Geographies, Geography, Marker} from 'react-simple-maps'
import {scaleSequential} from 'd3-scale'
import {interpolatePlasma} from 'd3-scale-chromatic'
import CountryData from './CountryData'

interface RequestData {
  lat: number
  lng: number
  intensity: number
}

interface GlobeMapProps {
  data: RequestData[]
}

const GlobeMap: React.FC<GlobeMapProps> = ({data}) => {
  const maxIntensity = Math.max(...data.map(d => d.intensity))
  const colorScale = scaleSequential(interpolatePlasma).domain([0, maxIntensity])

  return (
    <div className="w-full  h-full flex    bg-white">
      <ComposableMap
        projection="geoMercator"
        projectionConfig={{
          scale: 70,
        }}
        width={500}
        height={300}
        className="w-2/3"
      >
        <Geographies geography="https://cdn.jsdelivr.net/npm/world-atlas@2/countries-110m.json">
          {({geographies}) =>
            geographies.map(geo => (
              <Geography
                key={geo.rsmKey}
                geography={geo}
                fill="#EAEAEC"
                stroke="#D6D6DA"
              />
            ))
          }
        </Geographies>
        {data.map(({lat, lng, intensity}, index) => (
          <Marker key={index} coordinates={[lng, lat]}>
            <circle
              r={Math.sqrt(intensity) / 2}
              fill={colorScale(intensity)}
              stroke="#000"
              strokeWidth={0.5}
            />
          </Marker>
        ))}
      </ComposableMap>
      <div className="flex items-center justify-between">
        <CountryData />
      </div>
    </div>
  )
}

export default GlobeMap
