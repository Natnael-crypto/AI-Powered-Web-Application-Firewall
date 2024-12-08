import { LucideIcon } from "lucide-react";
import { useState } from "react";
import { SidebarGroup, SidebarMenu } from "../ui/sidebar";
import SidebarItem from "./SidebarItem";

type NavItem = {
    title: string;
    url: string;
    icon?: LucideIcon;
    isActive?: boolean;
    items?: {
      title: string;
      url: string;
    }[];
  };
type MainSidebarContentProps = {
    items: NavItem[];
  };
function MainSidebarContent({ items }: MainSidebarContentProps) {
    const [openTab, setOpenTab] = useState<string | null>(
      items.find((item) => item.isActive)?.title || null
    );
  
    const handleToggle = (title: string) => {
      setOpenTab((prev) => (prev === title ? null : title));
    };
  
    return (
      
      <SidebarGroup className="gap-y-5">
        <SidebarMenu className="gap-y-5">
          {items.map((item) => (
            <SidebarItem
              key={item.title}
              item={item}
              isOpen={openTab === item.title}
              onToggle={() => handleToggle(item.title)}
            />
          ))}
        </SidebarMenu>
      </SidebarGroup>
    );
  }

  export default MainSidebarContent;