import React from 'react'
import {scaleSequential} from 'd3-scale'
import {interpolateRgb} from 'd3-interpolate'
import {VectorMap} from '@react-jvectormap/core'
import {worldMill} from '@react-jvectormap/world'
import {requestData} from '../lib/Constants'
import RequestStat from './RequestStat'

const maxRequests = Math.max(...Object.values(requestData), 1)

const getColor = (requests: number) => {
  const scale = scaleSequential(interpolateRgb('#D3D3D3', '#006400')).domain([
    0,
    maxRequests,
  ])
  return scale(requests)
}

const GlobeMap: React.FC = () => {
  return (
    <div className="w-full h-full flex justify-between bg-white">
      <VectorMap
        backgroundColor="#f4f4f4"
        className="h-full w-full"
        zoomOnScroll={false}
        regionStyle={{
          initial: {fill: '#D3D3D3', stroke: '#fff', strokeWidth: 1},
          hover: {fill: '#228B22', cursor: 'pointer'},
        }}
        series={{
          regions: [
            {
              values: Object.fromEntries(
                Object.entries(requestData).map(([country, requests]) => [
                  country,
                  getColor(requests),
                ]),
              ),
              attribute: 'fill',
            },
          ],
        }}
        onRegionTipShow={(e, el: any, code) => {
          const requests = requestData[code] || 0
          el.html(
            `<strong>${el.html()}</strong><br>
            Requests: <b>${requests.toLocaleString()}</b>`,
          )
        }}
        map={worldMill}
      />
      <RequestStat />
    </div>
  )
}

export default GlobeMap
