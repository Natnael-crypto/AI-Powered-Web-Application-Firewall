import {useState} from 'react'
import {SiderbarContentItems} from '../lib/Constants'
import SidebarContent from './SidebarContent'
import logo from '../assets/waf-Logo.png'

function Sidebar() {
  const [openItem, setOpenItem] = useState({
    title: SiderbarContentItems[0].title,
    href: SiderbarContentItems[0].href,
  })

  const handleItemClick = (item: {title: string; href: string}) => {
    setOpenItem(item)
  }

  return (
    <div className="h-full hidden lg:flex flex-col items-center w-[16%] bg-gradient-to-b from-green-800  to-green-300 shadow-xl py-6">
      <div className="mb-6 w-[70%]">
        <img src={logo} alt="Logo" className="h-20 w-full object-contain" />
      </div>
      <div className="flex flex-col gap-4 w-full px-4 overflow-y-auto scrollbar-thin scrollbar-thumb-green-300 scrollbar-track-transparent">
        {SiderbarContentItems.map(item => (
          <SidebarContent
            key={item.title}
            title={item.title}
            href={item.href}
            children={item.children}
            openItem={openItem}
            changeOpenItem={handleItemClick}
          />
        ))}
      </div>
    </div>
  )
}

export default Sidebar
