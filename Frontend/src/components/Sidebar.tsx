import {useState} from 'react'
import {SiderbarContentItems} from '../lib/Constants'
import SidebarContent from './SidebarContent'

function Sidebar() {
  const [openItem, setOpenItem] = useState({
    title: SiderbarContentItems[0].title,
    href: SiderbarContentItems[0].href,
  })

  const handleItemClick = (item: {title: string; href: string}) => {
    setOpenItem(item)
  }

  return (
    <div className="h-full  flex flex-col items-center py-7 w-[15%]">
      <div className=" text-3xl text-black font-semibold">Logo</div>
      <div className="flex flex-col justify-center py-10 gap-3 overflow-y-scroll w-[80%]">
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
