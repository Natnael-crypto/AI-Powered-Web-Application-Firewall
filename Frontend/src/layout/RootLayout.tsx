import {Outlet} from 'react-router-dom'
import Sidebar from '../components/Sidebar'
import Navbar from '../components/Navbar'

function RootLayout() {
  return (
    <div className="h-screen w-full flex bg-slate-100">
      <Sidebar />
      <div className="flex flex-col w-[85%]">
        <Navbar />
        <Outlet />
      </div>
    </div>
  )
}

export default RootLayout
