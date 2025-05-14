import {useState, useCallback} from 'react'
import {SiderbarContentItems} from '../lib/Constants'
import SidebarContent from './SidebarContent'
import {useUserInfo} from '../store/UserInfo'

function Sidebar() {
  const {user} = useUserInfo()
  const [openItem, setOpenItem] = useState({
    title: SiderbarContentItems[0].title,
    href: SiderbarContentItems[0].href,
  })

  const handleItemClick = useCallback((item: {title: string; href: string}) => {
    setOpenItem(item)
  }, [])

  return (
    <div
      className="h-full hidden lg:flex flex-col items-center w-[16%] shadow-xl py-6 relative"
      style={{backgroundColor: '#1F263E'}}
    >
      <div className="flex flex-col w-full h-full">
        {/* Logo Section */}
        <br />
        <div className="mb-6 w-[70%] self-center">
          <p className='text-white text-3xl text-center'>GASHA WAF</p>
        </div>

        <br />
        {/* Navigation Items */}
        <div className="flex flex-col gap-4 w-full px-4 overflow-y-auto scrollbar-thin scrollbar-thumb-green-300 scrollbar-track-transparent flex-grow">
          {SiderbarContentItems.map(item => (
            <SidebarContent
              key={item.title}
              title={item.title}
              href={item.href}
              children={item.children}
              openItem={openItem}
              changeOpenItem={handleItemClick}
              Icon={item.icon}
            />
          ))}
        </div>

        {/* User Info Section - Fixed at the bottom */}
        {user && (
          <div className="mt-auto pt-4 border-t border-gray-600 px-4 w-full">
            <div className="text-white font-medium truncate">{user.username}</div>
            <div className="text-gray-400 text-sm truncate">
              {user.role || 'No email available'}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

export default Sidebar
