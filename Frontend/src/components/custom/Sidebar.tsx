"use client"

import * as React from "react"
import {
  PieChart,
  Settings2,
  SquareTerminal,
  Logs,
  Shield
} from "lucide-react"

import MainSidebarContent from "./MainSidebarContent"
import {
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarRail,
} from "../ui/sidebar"

const data = {
  navMain: [
    {
      title: "Statistics",
      url: "/",
      icon: PieChart,
      isActive: true,
      items: [
        {
          title: "Dashboard",
          url: "#",
        }
      ]
    },
    {
      title: "Logs",
      url: "#",
      icon: Logs,
      isActive: false,
      items: [
        {
          title: "Attacks",
          url: "#",
        },
        {
          title: "Rate Limiting",
          url: "#",
        },
        {
          title: "anti bot",
          url: "#",
        },
        {
          title: "auth",
          url: "#",
        },
        {
          title: "Waiting Room",
          url: "#",
        }
      ],
    },
    {
      title: "Web Services",
      url: "#",
      icon: Shield,
      isActive: false,
      items: [
        {
          title: "web services",
          url: "#",
        },
        {
          title: "SSL Cert",
          url: "#",
        },
        {
          title: "Global Settings",
          url: "#",
        }
      ],
    },
    {
      title: "Protection",
      url: "#",
      icon: Settings2,
      isActive: false,
      items: [
        {
          title: "Rate Limiting",
          url: "#",
        },
        {
          title: "Custom Rules",
          url: "#",
        },
        {
          title: "Detection mod",
          url: "#",
        },
        {
          title: "Settings",
          url: "#",
        },
      ],
    },
    {
      title: "System",
      url: "#",
      icon: SquareTerminal,
    }
  ],
}

export function CustomSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    
    <Sidebar  {...props} className="min-w-[300px] bg-slate-100 border-none">
      <SidebarHeader className=" text-center text-3xl my-5">
        LOGO
      </SidebarHeader>
      <SidebarContent className="">
        <MainSidebarContent items={data.navMain} />
      </SidebarContent>
      
      <SidebarRail />
    </Sidebar>
  )
}
