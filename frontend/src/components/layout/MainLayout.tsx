import { type ReactNode } from 'react'
import { SidebarProvider } from '@/components/ui/sidebar'
import { Header } from './Header'
import { AppSidebar } from './Sidebar'

interface MainLayoutProps {
  children: ReactNode
  className?: string
}

export function MainLayout({ children, className }: MainLayoutProps) {
  return (
    <SidebarProvider>
      <div className={`min-h-screen flex w-full ${className || ''}`}>
        {/* Sidebar */}
        <AppSidebar />
        
        {/* Main content area */}
        <div className="flex flex-1 flex-col overflow-hidden">
          {/* Header */}
          <Header />
          
          {/* Page content */}
          <main className="flex-1 overflow-y-auto">
            <div className="container mx-auto p-6 space-y-6">
              {children}
            </div>
          </main>
        </div>
      </div>
    </SidebarProvider>
  )
}