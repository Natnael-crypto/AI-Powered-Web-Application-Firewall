import {AiOutlineRight} from 'react-icons/ai'
import {Link, useLocation} from 'react-router-dom'
import SidebarItem from './SidebarItem'
import {useEffect, useState} from 'react'

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
}

function SidebarContent({title, href, children, changeOpenItem}: SidebarItemProps) {
  const location = useLocation()
  const [isOpen, setOpen] = useState(false)
  const [_, setActive] = useState(false)

  useEffect(() => {
    const currentSection = location.pathname.split('/')[1]?.toLowerCase()
    const titleLower = title.toLowerCase()
    const active = currentSection === titleLower
    console.log('test', currentSection, titleLower)

    setActive(active)
    setOpen(active)
  }, [location.pathname, title])

  const hasChildren = children && children.length > 0

  return (
    <div className="transition-all duration-200 ease-in-out w-full">
      <Link
        to={href}
        onClick={() => changeOpenItem?.({title, href})}
        className={`group flex w-full justify-between items-center gap-3 px-5 ${
          hasChildren ? 'py-3' : 'py-3.5'
        } rounded-2xl transition-all duration-300 cursor-pointer 
    ${isOpen ? 'bg-green-200 text-green-800 shadow-md' : 'text-gray-700 hover:bg-green-100 hover:text-green-600'}
  `}
      >
        <h2
          className={`text-[15px] font-semibold transition duration-300 group-hover:scale-[1.02]`}
        >
          {title}
        </h2>
        {hasChildren && (
          <AiOutlineRight
            className={`w-4 h-4 transition-transform duration-300 ease-in-out ${isOpen ? 'rotate-90' : ''}`}
          />
        )}
      </Link>
      {isOpen && hasChildren && (
        <div className="flex flex-col gap-1 mt-1 pl-4 border-l-2 border-green-300 ml-2">
          {children.map(child => (
            <SidebarItem
              key={child.title}
              title={child.title}
              href={child.href}
              isActive={location.pathname === child.href}
              className={`text-sm py-2 px-4 rounded-xl transition-colors duration-200 ${
                location.pathname === child.href
                  ? 'bg-green-100 text-green-700'
                  : 'hover:bg-green-50 hover:text-green-600 text-gray-600'
              }`}
            />
          ))}
        </div>
      )}
    </div>
  )
}

export default SidebarContent
