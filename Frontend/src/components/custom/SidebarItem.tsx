import {
  Collapsible,
  CollapsibleTrigger,
  CollapsibleContent,
} from "@radix-ui/react-collapsible";
import { ChevronRight, LucideIcon } from "lucide-react";
import {
  SidebarMenuItem,
  SidebarMenuButton,
  SidebarMenuSub,
  SidebarMenuSubItem,
  SidebarMenuSubButton,
} from "../ui/sidebar";

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

function SidebarItem({
  item,
  isOpen,
  onToggle,
}: {
  item: NavItem;
  isOpen: boolean;
  onToggle: () => void;
}) {
  return (
    <Collapsible asChild open={isOpen} className="group/collapsible">
      <SidebarMenuItem>
        <CollapsibleTrigger asChild>
          <SidebarMenuButton
            tooltip={item.title}
            className={`text-xl ${isOpen ? " bg-green-300" : ""} p-7 px-10 `}
            onClick={onToggle}
          >
            <div className="block mr-5">
              {item.icon && <item.icon size={25} />}
            </div>
            <span>{item.title}</span>
            {(item.items?.length ?? 0) > 0 && (
              <ChevronRight
                className={`ml-auto transition-transform duration-200 ${
                  isOpen ? "rotate-90" : ""
                }`}
              />
            )}
          </SidebarMenuButton>
        </CollapsibleTrigger>
        <CollapsibleContent>
          <SidebarMenuSub className="my-3">
            {item.items?.map((subItem) => (
              <SidebarMenuSubItem
                key={subItem.title}
                className="px-2 py-2 bg-slate-100 rounded-lg"
              >
                <SidebarMenuSubButton asChild>
                  <a href={subItem.url}>
                    <span className="text-lg">{subItem.title}</span>
                  </a>
                </SidebarMenuSubButton>
              </SidebarMenuSubItem>
            ))}
          </SidebarMenuSub>
        </CollapsibleContent>
      </SidebarMenuItem>
    </Collapsible>
  );
}

export default SidebarItem;
