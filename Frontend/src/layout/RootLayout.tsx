import {Outlet, useLocation, useNavigate} from 'react-router-dom'
import Sidebar from '../components/Sidebar'
import Navbar from '../components/Navbar'
import {useEffect} from 'react'
import axios from 'axios'

function RootLayout() {
  const {pathname} = useLocation()
  const navigate = useNavigate()
  const backendUrl = import.meta.env.VITE_BACKEND_URL;
  const checkToken = ()=>{
    const token = localStorage.getItem('token')
    if (token) {
      axios.get(`${backendUrl}/is-logged-in`, {
        headers: {
          Authorization: token,
        },
      })
      .then((response) => {
        if (response.status !== 200) {
          navigate('/login');
        }
      })
      .catch((error) => {
        navigate('/login');
      });
    } else {
      navigate('/login');
    }

  }

  useEffect(() => {
    console.log('got you 1')
    checkToken()
    console.log('got you 2')
  }, [pathname])

  return (
    <div className="h-screen w-full flex bg-gradient-to-r from-slate-100 to-slate-50">
      <Sidebar />
      <div className="flex flex-col w-full overflow-hidden">
        <Navbar />
        <main className="flex-1 overflow-y-auto p-6" style={{backgroundColor:"#F3F6FE"}}>
          <Outlet />
        </main>
      </div>
    </div>
  )
}

export default RootLayout
