import Button from './atoms/Button'
interface dashboardMenuProps {
  selectedMenu: 'basic' | 'advanced'
  changeMenu: () => void
}

function DashboardMenu({changeMenu, selectedMenu}: dashboardMenuProps) {
  return (
    <div className=" flex items-center py-2 px-2 gap-3 shadow-md bg-white ">
      <Button
        classname={`px-4 py-1 ${selectedMenu == 'basic' ? 'bg-green-500' : 'border border-1 border-gray-400'}  text-lg `}
        onClick={changeMenu}
      >
        Basic
      </Button>
      <Button
        classname={`px-4 py-1 ${selectedMenu == 'advanced' ? 'bg-green-500' : 'border border-1 border-gray-400'} text-lg border-gray-300 `}
        onClick={changeMenu}
      >
        Advanced
      </Button>
      <Button classname={`px-4 py-1 text-lg border-gray-300 `} variant="outline">
        Data Dashboard
      </Button>
    </div>
  )
}

export default DashboardMenu
