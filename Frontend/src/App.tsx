import {
  createBrowserRouter,
  createRoutesFromElements,
  Route,
  RouterProvider,
} from 'react-router-dom'
import RootLayout from './layout/RootLayout'
import Dashboard from './pages/Dashboard'
import LoginPage from './pages/Login'
import LogAttack from './pages/LogAttacks'

const router = createBrowserRouter(
  createRoutesFromElements(
    <>
      <Route path="/" element={<RootLayout />}>
        <Route path="statistics/dashboard" element={<Dashboard />} />
        <Route path="logs/attacks" element={<LogAttack />} />
      </Route>
      <Route path="/login" element={<LoginPage />} />
    </>,
  ),
)

function App() {
  return (
    <>
      <RouterProvider router={router} />
    </>
  )
}

export default App
