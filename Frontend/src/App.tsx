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
import PageNotFound from './pages/PageNotFound'
import Application from './pages/Application'
import CustomeRules from './pages/CustomeRules'

const router = createBrowserRouter(
  createRoutesFromElements(
    <>
      <Route path="/" element={<RootLayout />}>
        <Route path="statistics/dashboard" element={<Dashboard />} />
        <Route path="logs/attacks" element={<LogAttack />} />
        <Route path="application/applications" element={<Application />} />
        <Route path="protection/custom_rules" element={<CustomeRules />} />
        <Route path="*" element={<PageNotFound />} />
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
