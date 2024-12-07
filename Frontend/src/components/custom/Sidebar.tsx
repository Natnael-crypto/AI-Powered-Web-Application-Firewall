"use client"

import * as React from "react"
import {
  PieChart,
  Settings2,
  SquareTerminal,
  Logs,
  Shield
} from "lucide-react"

import { NavMain } from "./NavMain"
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
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={data.navMain} />
      </SidebarContent>
      
      <SidebarRail />
    </Sidebar>
  )
}
