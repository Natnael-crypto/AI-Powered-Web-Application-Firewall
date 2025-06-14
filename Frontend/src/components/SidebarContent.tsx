import {AiOutlineRight} from 'react-icons/ai'
import {Link, useLocation} from 'react-router-dom'
import SidebarItem from './SidebarItem'
import {useEffect, useState} from 'react'
import {useUserInfo} from '../store/UserInfo'
import {Roles} from '../lib/types'

interface SidebarItemType {
  title: string
  href: string
}

interface SidebarItemProps {
  title: string
  href: string
  children?: SidebarItemType[]
  className?: string
  openItem?: SidebarItemType
  changeOpenItem?: (item: SidebarItemType) => void
  Icon?: React.ComponentType<{className?: string}>
}

function SidebarContent({title, href, children, changeOpenItem, Icon}: SidebarItemProps) {
  const location = useLocation()
  const [isOpen, setOpen] = useState(false)
  const [_, setActive] = useState(false)
  const {user} = useUserInfo()

  const normalizeString = (str: string) => {
    return str
      .toLowerCase()
      .replace(/\s+/g, '-')
      .replace(/[^a-z0-9-]/g, '')
  }

  useEffect(() => {
    const currentSection = location.pathname.split('/')[1]?.toLowerCase()
    const normalizedTitle = normalizeString(title)
    const normalizedHref = normalizeString(href)
    const active = currentSection === normalizedTitle || currentSection === normalizedHref
    setActive(active)
    setOpen(active)
  }, [location.pathname, title, href])

  const hasChildren = children && children.length > 0

  if (title.toLocaleLowerCase() === 'system' && user?.role != Roles.SUPER_ADMIN)
    return null

  return (
    <div className="transition-all duration-200 ease-in-out w-full ">
      <Link
        to={href}
        onClick={() => changeOpenItem?.({title, href})}
        className={`group flex w-full items-center justify-between gap-x-3 px-6 py-3  transition-all duration-300 cursor-pointer
          ${isOpen ? 'bg-[#303750] text-white shadow-md' : 'text-gray-500 hover:bg-gray-600 hover:text-white'}
        `}
      >
        <div className={`flex items-center gap-x-3 `}>
          {Icon && (
            <Icon
              className={`w-5 h-5 shrink-0 transition-colors duration-300 ${
                isOpen ? 'text-white' : 'text-gray-400 group-hover:text-white'
              }`}
            />
          )}
          <h2
            className={`text-[15px] font-normal truncate ${
              isOpen ? 'text-white' : 'text-gray-400'
            }`}
          >
            {title}
          </h2>
        </div>

        {hasChildren && (
          <AiOutlineRight
            className={`w-4 h-4 transition-transform duration-300 ease-in-out ${
              isOpen ? 'rotate-90 text-white' : 'text-gray-400'
            }`}
          />
        )}
      </Link>

      {isOpen && hasChildren && (
        <div className="flex flex-col gap-1 mt-2 pl-9 border-l-2 border-green-300">
          {children.map(child => (
            <SidebarItem
              key={child.title}
              title={child.title}
              href={child.href}
              isActive={location.pathname === child.href}
              className={`text-sm py-2 px-3  transition-colors duration-200 ${
                location.pathname.replace(/^\/+/, '') === child.href.replace(/^\/+/, '')
                  ? 'bg-[#303750] text-white'
                  : 'text-gray-500 hover:bg-gray-600 text-white'
              }`}
            />
          ))}
        </div>
      )}
    </div>
  )
}

export default SidebarContent
