import {Outlet, useLocation, useNavigate} from 'react-router-dom'
import Sidebar from '../components/Sidebar'
import Navbar from '../components/Navbar'
import {useEffect} from 'react'

function RootLayout() {
  const {pathname} = useLocation()
  const navigate = useNavigate()

  useEffect(() => {
    if (!localStorage.getItem('token')) {
      navigate('/login')
    }
  }, [pathname])

  return (
    <div className="h-screen w-full flex bg-gradient-to-r from-slate-100 to-slate-50">
      <Sidebar />
      <div className="flex flex-col w-full overflow-hidden">
        <Navbar />
        <main className="flex-1 overflow-y-auto p-6">
          <Outlet />
        </main>
      </div>
    </div>
  )
}

export default RootLayout
