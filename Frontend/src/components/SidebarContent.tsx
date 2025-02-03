import {AiOutlineRight} from 'react-icons/ai'
import {Link, useLocation} from 'react-router-dom'
import SidebarItem from './SidebarItem'

interface sidebarItem {
  title: string
  href: string
}

interface sidebarItemProps {
  title: string
  href: string
  children?: sidebarItem[]
  className?: string
  openItem?: sidebarItem
  changeOpenItem?: (item: sidebarItem) => void
}

function SidebarContent({
  title,
  href,
  children,
  className,
  openItem,
  changeOpenItem,
}: sidebarItemProps) {
  const isOpen = openItem?.title === title
  const url = useLocation()

  const isActive = url.pathname.split('/')[-1] === href.split('/')[-1]

  return (
    <div className="transition-all duration-200 ease-in-out w-full ">
      <Link
        to={href}
        onClick={() => changeOpenItem && changeOpenItem({title, href})}
        className={`flex w-full justify-between items-center gap-4 ${children?.length ? 'py-3 ' : ''} px-7 py-4 hover:bg-green-50 hover:text-green-600 rounded-xl transition-colors duration-200 ${
          isOpen ? 'bg-green-50 text-green-600 mb-3' : 'text-gray-700'
        }`}
      >
        <h2 className={`text-lg font-semibold ${className}`}>{title}</h2>
        {children && children.length > 0 && (
          <AiOutlineRight
            className={`w-5 h-5 transition-transform duration-200 ${
              isOpen ? 'transform rotate-90' : ''
            }`}
          />
        )}
      </Link>
      {isOpen && children && (
        <div className={`flex flex-col gap-1 `}>
          {children.map(child => (
            <SidebarItem
              isActive={isActive}
              key={child.title}
              href={child.href}
              title={child.title}
              className={`font-light px-8 py-4 text-gray-600`}
            />
          ))}
        </div>
      )}
    </div>
  )
}

export default SidebarContent
