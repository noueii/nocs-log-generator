import { 
  Play, 
  FileVideo, 
  History, 
  Settings, 
  Target,
  BarChart3
} from 'lucide-react'
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
    title: "Generate Match",
    url: "#generate",
    icon: Play,
    description: "Create synthetic CS2 match logs"
  },
  {
    title: "Parse Demo",
    url: "#parse", 
    icon: FileVideo,
    description: "Convert demo files to logs"
  },
  {
    title: "Match History",
    url: "#history",
    icon: History,
    description: "View generated matches"
  },
]

const toolItems = [
  {
    title: "Statistics",
    url: "#stats",
    icon: BarChart3,
    description: "View match statistics"
  },
  {
    title: "CS2 Events",
    url: "#events",
    icon: Target,
    description: "Event type reference"
  },
]

const settingsItems = [
  {
    title: "Settings",
    url: "#settings",
    icon: Settings,
    description: "Application settings"
  },
]

interface AppSidebarProps {
  className?: string
}

export function AppSidebar({ className }: AppSidebarProps) {
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
                  <SidebarMenuButton asChild tooltip={item.description}>
                    <a href={item.url} className="flex items-center gap-2">
                      <item.icon className="size-4" />
                      <span>{item.title}</span>
                    </a>
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
                  <SidebarMenuButton asChild tooltip={item.description}>
                    <a href={item.url} className="flex items-center gap-2">
                      <item.icon className="size-4" />
                      <span>{item.title}</span>
                    </a>
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
                  <SidebarMenuButton asChild tooltip={item.description}>
                    <a href={item.url} className="flex items-center gap-2">
                      <item.icon className="size-4" />
                      <span>{item.title}</span>
                    </a>
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