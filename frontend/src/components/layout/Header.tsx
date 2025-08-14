import { useState } from 'react'
import { Moon, Sun } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SidebarTrigger } from '@/components/ui/sidebar'

interface HeaderProps {
  className?: string
}

export function Header({ className }: HeaderProps) {
  const [isDark, setIsDark] = useState(() => {
    if (typeof window !== 'undefined') {
      return document.documentElement.classList.contains('dark')
    }
    return false
  })

  const toggleTheme = () => {
    setIsDark(!isDark)
    if (!isDark) {
      document.documentElement.classList.add('dark')
      localStorage.setItem('theme', 'dark')
    } else {
      document.documentElement.classList.remove('dark')
      localStorage.setItem('theme', 'light')
    }
  }

  return (
    <header className={`sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 ${className || ''}`}>
      <div className="container flex h-14 items-center">
        {/* Mobile sidebar trigger */}
        <SidebarTrigger className="md:hidden" />
        
        {/* Logo and title */}
        <div className="flex items-center gap-2 flex-1 md:flex-none">
          <div className="flex items-center gap-2">
            <div className="size-8 rounded bg-gradient-to-br from-cs-orange to-cs-blue flex items-center justify-center">
              <span className="text-white font-bold text-sm">CS2</span>
            </div>
            <div className="flex flex-col">
              <h1 className="text-lg font-semibold">CS2 Log Generator</h1>
            </div>
          </div>
        </div>

        {/* Desktop navigation */}
        <nav className="hidden md:flex items-center gap-6 text-sm font-medium flex-1 justify-center">
          <a
            href="#generate"
            className="transition-colors hover:text-foreground/80 text-foreground/60"
          >
            Generate
          </a>
          <a
            href="#parse"
            className="transition-colors hover:text-foreground/80 text-foreground/60"
          >
            Parse
          </a>
          <a
            href="#history"
            className="transition-colors hover:text-foreground/80 text-foreground/60"
          >
            History
          </a>
        </nav>

        {/* Theme toggle */}
        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            size="icon"
            onClick={toggleTheme}
            className="size-9"
            aria-label="Toggle theme"
          >
            {isDark ? (
              <Sun className="size-4" />
            ) : (
              <Moon className="size-4" />
            )}
          </Button>
        </div>
      </div>
    </header>
  )
}