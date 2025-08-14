import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

function App() {
  const [isDark, setIsDark] = useState(false)

  const toggleTheme = () => {
    setIsDark(!isDark)
    if (!isDark) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }

  return (
    <div className="min-h-screen bg-background text-foreground p-8">
      <div className="container mx-auto max-w-4xl space-y-8">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-4xl font-bold text-foreground">CS2 Log Generator</h1>
            <p className="text-muted-foreground mt-2">Counter-Strike 2 HTTP Log Generator</p>
          </div>
          <Button onClick={toggleTheme} variant="outline">
            {isDark ? '☀️' : '🌙'} Toggle Theme
          </Button>
        </div>

        {/* Feature Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <Card>
            <CardHeader>
              <CardTitle>Generate Match Logs</CardTitle>
              <CardDescription>
                Create synthetic CS2 match logs with customizable teams and players
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex gap-2">
                <Button>Start Generation</Button>
                <Button variant="outline">Configure Teams</Button>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Parse Demo Files</CardTitle>
              <CardDescription>
                Convert CS2 demo files (.dem) to HTTP log format
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex gap-2">
                <Button variant="secondary">Upload Demo</Button>
                <Button variant="outline">Parse Settings</Button>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* CS2 Theme Demo */}
        <Card>
          <CardHeader>
            <CardTitle>CS2 Theme Demo</CardTitle>
            <CardDescription>
              Showcasing CS2-specific colors and components
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex gap-2 flex-wrap">
              <Button variant="ct">Counter-Terrorist</Button>
              <Button variant="t">Terrorist</Button>
              <Button variant="default">Default</Button>
              <Button variant="destructive">Destructive</Button>
              <Button variant="outline">Outline</Button>
              <Button variant="secondary">Secondary</Button>
              <Button variant="ghost">Ghost</Button>
              <Button variant="link">Link</Button>
            </div>
            
            <div className="grid grid-cols-2 gap-4 mt-6">
              <div className="p-4 rounded-lg bg-ct text-white text-center">
                Counter-Terrorist Theme
              </div>
              <div className="p-4 rounded-lg bg-t text-white text-center">
                Terrorist Theme
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Status */}
        <Card>
          <CardHeader>
            <CardTitle>Setup Status</CardTitle>
            <CardDescription>
              Tailwind CSS 4 and shadcn/ui configuration verification
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-2 text-sm">
              <div className="flex items-center gap-2">
                <span className="w-2 h-2 bg-green-500 rounded-full"></span>
                <span>Tailwind CSS 4 configured with CS2 colors</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="w-2 h-2 bg-green-500 rounded-full"></span>
                <span>shadcn/ui components installed</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="w-2 h-2 bg-green-500 rounded-full"></span>
                <span>Dark/Light theme switching working</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="w-2 h-2 bg-green-500 rounded-full"></span>
                <span>CS2-specific button variants available</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

export default App
