import { 
  Play, 
  FileVideo, 
  History, 
  Settings, 
  Target,
  BarChart3,
  Home
} from 'lucide-react'
import { Link, useLocation } from 'react-router-dom'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarSeparator,
} from '@/components/ui/sidebar'

// Menu items
const menuItems = [
  {
    title: "Home",
    url: "/",
    icon: Home,
    description: "Dashboard and overview"
  },
  {
    title: "Generate Match",
    url: "/generate",
    icon: Play,
    description: "Create synthetic CS2 match logs"
  },
  {
    title: "Match History",
    url: "/history",
    icon: History,
    description: "View generated matches"
  },
  {
    title: "Parse Demo",
    url: "/parse", 
    icon: FileVideo,
    description: "Convert demo files to logs (coming soon)"
  },
]

const toolItems = [
  {
    title: "Statistics",
    url: "/stats",
    icon: BarChart3,
    description: "View match statistics (coming soon)"
  },
  {
    title: "CS2 Events",
    url: "/events",
    icon: Target,
    description: "Event type reference (coming soon)"
  },
]

const settingsItems = [
  {
    title: "Settings",
    url: "/settings",
    icon: Settings,
    description: "Application settings (coming soon)"
  },
]

interface AppSidebarProps {
  className?: string
}

export function AppSidebar({ className }: AppSidebarProps) {
  const location = useLocation();
  
  return (
    <Sidebar variant="inset" className={className}>
      <SidebarHeader>
        <div className="flex items-center gap-2 px-2 py-1">
          <div className="size-6 rounded bg-gradient-to-br from-cs-orange to-cs-blue flex items-center justify-center">
            <span className="text-white font-bold text-xs">CS2</span>
          </div>
          <span className="font-semibold text-sm">Log Generator</span>
        </div>
      </SidebarHeader>

      <SidebarContent>
        {/* Main Actions */}
        <SidebarGroup>
          <SidebarGroupLabel>Main Actions</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {menuItems.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild tooltip={item.description} isActive={location.pathname === item.url}>
                    <Link to={item.url} className="flex items-center gap-2">
                      <item.icon className="size-4" />
                      <span>{item.title}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        <SidebarSeparator />

        {/* Tools & Analytics */}
        <SidebarGroup>
          <SidebarGroupLabel>Tools & Analytics</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {toolItems.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild tooltip={item.description} isActive={location.pathname === item.url}>
                    <Link to={item.url} className="flex items-center gap-2 opacity-50 pointer-events-none">
                      <item.icon className="size-4" />
                      <span>{item.title}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        <SidebarSeparator />

        {/* Configuration */}
        <SidebarGroup>
          <SidebarGroupLabel>Configuration</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {settingsItems.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild tooltip={item.description} isActive={location.pathname === item.url}>
                    <Link to={item.url} className="flex items-center gap-2 opacity-50 pointer-events-none">
                      <item.icon className="size-4" />
                      <span>{item.title}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter>
        <div className="px-2 py-1 text-xs text-muted-foreground">
          <div className="flex items-center justify-between">
            <span>CS2 Log Generator</span>
            <span>v1.0.0</span>
          </div>
        </div>
      </SidebarFooter>
    </Sidebar>
  )
}