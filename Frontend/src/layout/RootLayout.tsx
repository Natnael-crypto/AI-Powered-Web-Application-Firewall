import {Outlet} from 'react-router-dom'
import {CustomSidebar} from '../components/custom/Sidebar'
import Navbar from '../components/custom/Navbar'

function RootLayout() {
  return (
    <div className="h-full w-full flex">
      <CustomSidebar />
      <div>
        <Navbar />
        <Outlet />
      </div>
    </div>
  )
}

export default RootLayout
