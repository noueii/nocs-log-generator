import { MainLayout } from '@/components/layout'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

function App() {
  return (
    <MainLayout>
      {/* Welcome Section */}
      <div className="space-y-2">
        <h1 className="text-3xl font-bold tracking-tight">Welcome to CS2 Log Generator</h1>
        <p className="text-muted-foreground">
          Generate synthetic Counter-Strike 2 match logs and parse demo files for testing and analysis.
        </p>
      </div>

      {/* Feature Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <div className="w-2 h-2 bg-cs-blue rounded-full"></div>
              Generate Match Logs
            </CardTitle>
            <CardDescription>
              Create synthetic CS2 match logs with customizable teams and players
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex gap-2">
              <Button variant="ct">Start Generation</Button>
              <Button variant="outline">Configure Teams</Button>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <div className="w-2 h-2 bg-cs-orange rounded-full"></div>
              Parse Demo Files
            </CardTitle>
            <CardDescription>
              Convert CS2 demo files (.dem) to HTTP log format
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex gap-2">
              <Button variant="t">Upload Demo</Button>
              <Button variant="outline">Parse Settings</Button>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* CS2 Theme Demo */}
      <Card>
        <CardHeader>
          <CardTitle>CS2 Theme Components</CardTitle>
          <CardDescription>
            Showcasing CS2-specific colors and components with the new layout
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
          <CardTitle>Layout Setup Status</CardTitle>
          <CardDescription>
            Verification of layout components and responsive design
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-2 text-sm">
            <div className="flex items-center gap-2">
              <span className="w-2 h-2 bg-green-500 rounded-full"></span>
              <span>Header component with theme toggle</span>
            </div>
            <div className="flex items-center gap-2">
              <span className="w-2 h-2 bg-green-500 rounded-full"></span>
              <span>Sidebar with navigation menu</span>
            </div>
            <div className="flex items-center gap-2">
              <span className="w-2 h-2 bg-green-500 rounded-full"></span>
              <span>Responsive mobile layout</span>
            </div>
            <div className="flex items-center gap-2">
              <span className="w-2 h-2 bg-green-500 rounded-full"></span>
              <span>MainLayout component integration</span>
            </div>
            <div className="flex items-center gap-2">
              <span className="w-2 h-2 bg-green-500 rounded-full"></span>
              <span>shadcn/ui components working</span>
            </div>
          </div>
        </CardContent>
      </Card>
    </MainLayout>
  )
}

export default App