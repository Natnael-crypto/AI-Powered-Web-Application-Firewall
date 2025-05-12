import {Link, useLocation} from 'react-router-dom'
import {IconType} from 'react-icons'

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
  Icon?: IconType // Add this to accept an optional Icon
}

function SidebarItem({
  title,
  href,
  className,
  changeOpenItem,
  Icon, // Get the Icon prop
}: SidebarItemProps) {
  const url = useLocation()
  const isActive = url.pathname.slice(1) === href

  return (
    <div className="transition-all duration-200 ease-in-out">
      <Link
        to={href}
        onClick={() => changeOpenItem && changeOpenItem({title, href})}
        className={`flex w-full items-center gap-x-4 px-7 py-3 hover:bg-green-50 hover:text-gray-600  transition-colors duration-200 ${isActive ? '' : ''} ${className}`}
      >
        {Icon && <Icon className="w-5 h-5 text-gray-400" />}{' '}
        {/* Render the Icon if provided */}
        <h2 className="text-sm font-normal text-left flex-1">{title}</h2>{' '}
        {/* Ensure text aligns left and occupies available space */}
      </Link>
    </div>
  )
}

export default SidebarItem
