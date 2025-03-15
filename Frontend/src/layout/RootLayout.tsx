import {Outlet, useLocation, useNavigate} from 'react-router-dom'
import Sidebar from '../components/Sidebar'
import Navbar from '../components/Navbar'
import {useEffect} from 'react'

function RootLayout() {
  const {pathname} = useLocation()
  const navigate = useNavigate()
  // useEffect(() => {
  //   if (!localStorage.getItem('token')) {
  //     navigate('/login')
  //   }
  // }, [pathname])
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
