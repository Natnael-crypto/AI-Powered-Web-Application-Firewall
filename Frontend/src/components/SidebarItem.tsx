import {Link, useLocation} from 'react-router-dom'

interface SidebarItem {
  title: string
  href: string
}

interface SidebarItemProps {
  title: string
  href: string
  children?: SidebarItem[]
  className?: string
  openItem?: SidebarItem
  changeOpenItem?: (item: SidebarItem) => void
  isActive?: boolean
}

function SidebarItem({title, href, className, changeOpenItem}: SidebarItemProps) {
  const url = useLocation()
  const isActive = url.pathname.slice(1) == href

  return (
    <div className="transition-all duration-200 ease-in-out">
      <Link
        to={href}
        onClick={() => changeOpenItem && changeOpenItem({title, href})}
        className={`flex w-full justify-between items-center gap-4 px-7 hover:bg-green-50 hover:text-green-600 rounded-xl transition-colors duration-200 ${isActive ? '' : ''} ${className}`}
      >
        <h2 className={`text-lg font-semibold`}>{title}</h2>
      </Link>
    </div>
  )
}

export default SidebarItem
