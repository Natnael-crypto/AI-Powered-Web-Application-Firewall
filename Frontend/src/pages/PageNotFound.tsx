import {Link} from 'react-router-dom'

const PageNotFound = () => {
  return (
    <div className="flex flex-col items-center justify-center h-screen bg-gray-100 text-gray-800">
      <h1 className="text-4xl font-bold">Under Construction</h1>
      <p className="mt-2 text-lg">This page is currently unavailable.</p>
      <Link to="/dashboard">
        <button className="mt-4 px-6 py-2 bg-blue-600 text-white  shadow-md hover:bg-blue-700">
          Go Home
        </button>
      </Link>
    </div>
  )
}

export default PageNotFound
