# Product Requirements Document: CS2 HTTP Log Generator

## Executive Summary

The CS2 HTTP Log Generator is a comprehensive tool designed to generate and parse Counter-Strike 2 match logs in HTTP format. The system supports two primary modes: synthetic log generation for simulating matches with customizable parameters, and demo file parsing for converting recorded gameplay (.dem files) into HTTP log format. This tool enables developers, analysts, and server operators to test log processing systems, analyze match data, and simulate realistic CS2 server environments.

## Project Goals & Objectives

### Primary Goals
1. Generate realistic CS2 match logs that accurately simulate actual game events
2. Parse and convert CS2 demo files (.dem) to HTTP log format
3. Support comprehensive match simulation including server rollbacks and edge cases
4. Provide flexible configuration for various match scenarios and formats
5. Ensure compatibility with existing CS2 log parsers and analytics tools

### Success Criteria
- Generated logs are indistinguishable from actual CS2 server logs
- Support for all major CS2 game events and match formats
- Processing speed of <30 seconds for typical match generation
- Demo parsing performance of ~25 minutes of gameplay per second
- 100% compatibility with HLTV and common log parsing tools

## Technical Requirements

### Backend Requirements
- **Language**: Go (recommended for performance and existing ecosystem)
- **HTTP Server**: Gin framework for HTTP endpoints
- **Demo Parsing**: Integration with demoinfocs-golang or similar libraries
- **Output Format**: Standard CS2 HTTP log format with proper timestamps
- **Performance**: Support for concurrent log generation and parsing

### Frontend Requirements
- **Framework**: React 19 with TypeScript
- **Node Version**: 20+ LTS
- **Build Tool**: Vite 5+
- **Package Manager**: pnpm (recommended) or npm
- **Browser Support**: Chrome 90+, Firefox 88+, Safari 14+, Edge 90+

### Dependencies
#### Backend
- Go 1.21+
- Gin web framework
- Demo parsing library (demoinfocs-golang recommended)
- Protocol Buffers support for demo parsing
- JSON/YAML configuration parsing
- gorilla/websocket for WebSocket support

#### Frontend
- React 19
- TypeScript 5+
- Tailwind CSS 4 (with Oxide engine)
- shadcn/ui components
- Tanstack Query for server state
- Zustand for client state
- Axios or ky for API calls
- Socket.io-client for WebSocket
- Recharts or Tremor for data visualization
- Radix UI primitives (via shadcn/ui)
- Lucide React for icons

## Core Features

### 1. Synthetic Log Generation

#### Match Configuration
- **Teams**: Configure team names, countries, and rankings
- **Players**: 
  - Steam IDs (authentic format: STEAM_1:X:XXXXXXXX)
  - Player names and aliases
  - Skill ratings and performance profiles
  - Role assignments (AWPer, Entry, Support, etc.)
- **Match Format**:
  - MR12 (CS2 standard): 24 rounds, first to 13
  - MR15 (legacy CS:GO): 30 rounds, first to 16
  - Custom round counts
  - Overtime configuration (MR3, MR6)
- **Map Selection**: All active duty maps plus custom map support

#### Event Simulation

##### Player Events
```
- player_connect: Player joins server
- player_disconnect: Player leaves server
- player_spawn: Player spawns in round
- player_death: Kill event with attacker, victim, weapon, headshot, penetration
- player_hurt: Damage dealt with health/armor reduction
- player_team: Team switch events
```

##### Round Events
```
- round_start: Round begins with team economy states
- round_freeze_end: Buy time ends, round action begins
- round_end: Round conclusion with winner and reason
- round_mvp: MVP selection based on performance
```

##### Weapon/Combat Events
```
- weapon_fire: Shot fired with position and direction
- item_purchase: Weapon/equipment purchases
- grenade_thrown: Grenade deployment
- bullet_impact: Bullet hit registration
```

##### Objective Events
```
- bomb_planted: Bomb plant with planter and site
- bomb_defused: Successful defuse with defuser info
- bomb_exploded: Round end by explosion
- bomb_dropped/pickup: Bomb carrier changes
- hostage_rescued: Hostage mode events
```

#### Economy Simulation
- Starting money: $800 (pistol round)
- Win bonus progression: $3250 (standard), $3500 (after loss)
- Loss bonus: $1400 → $1900 → $2400 → $2900 → $3400
- Kill rewards: $300 (rifle), $600 (SMG), $100 (AWP)
- Objective bonuses: $300 (plant), $250 (defuse)
- Equipment costs matching current CS2 economy

#### Advanced Features
- **Server Rollbacks**: Simulate backup restoration scenarios
  - Random rollback probability configuration
  - Rollback to specific round numbers
  - State restoration simulation
- **Network Issues**: Lag, packet loss, disconnection events
- **Anti-cheat Events**: VAC authentication, kick/ban events
- **Admin Commands**: Server commands and configuration changes

### 2. Demo File Parsing

#### Input Support
- CS2 .dem files (Source 2 engine format)
- CS:GO legacy demo support
- POV and GOTV demo types
- Compressed and uncompressed formats

#### Event Extraction
- Complete event timeline with tick-accurate timestamps
- Player position and movement data
- Weapon states and inventory
- Grenade trajectories and effects
- Economy tracking per round
- Chat messages and radio commands

#### Conversion Process
1. Parse binary demo format using demoinfocs-golang
2. Extract events with tick timestamps
3. Convert ticks to server time (64 tick = 64 updates/second)
4. Format events as CS2 HTTP log lines
5. Apply proper timestamp formatting
6. Output via HTTP endpoint or file

### 3. Rollback Simulation

#### Backup System
- Automatic round backup generation
- Configurable backup frequency
- State preservation:
  - Scores and round count
  - Player statistics (K/D/A)
  - Economy state (money, equipment)
  - Round history

#### Rollback Scenarios
- **Technical Issues**: Server crash, network failure
- **Competitive Integrity**: Bug exploitation, rule violations
- **Admin Intervention**: Manual rollback requests
- **Random Occurrence**: Configurable probability for testing

#### Implementation
```
- Generate backup files: backup_round_XX.txt
- Trigger rollback event in log
- Restore state from backup
- Continue match from restored point
- Log restoration confirmation
```

## Log Format Specifications

### Standard Log Line Format
```
L MM/DD/YYYY - HH:MM:SS: "PlayerName<ID><STEAM_ID><TEAM>" action "details"
```

### Example Log Events

#### Kill Event
```
L 01/14/2025 - 15:44:36: "Player1<12><STEAM_1:0:12345678><CT>" killed "Player2<24><STEAM_1:0:87654321><TERRORIST>" with "ak47" (headshot)
```

#### Round Start
```
L 01/14/2025 - 15:42:00: World triggered "Round_Start"
L 01/14/2025 - 15:42:00: Team "CT" scored "7" with "5" players
L 01/14/2025 - 15:42:00: Team "TERRORIST" scored "5" with "5" players
```

#### Bomb Plant
```
L 01/14/2025 - 15:43:45: "Player3<15><STEAM_1:0:11111111><TERRORIST>" triggered "Planted_The_Bomb" at bombsite A
```

#### Purchase Event
```
L 01/14/2025 - 15:42:15: "Player1<12><STEAM_1:0:12345678><CT>" purchased "m4a1"
```

### HTTP Endpoint Format
```
POST /log
Content-Type: application/json
{
  "timestamp": "2025-01-14T15:44:36Z",
  "event": "player_death",
  "data": {
    "attacker": {...},
    "victim": {...},
    "weapon": "ak47",
    "headshot": true
  }
}
```

## API Design

### RESTful Endpoints

#### Synthetic Generation
```
POST /api/generate
{
  "mode": "synthetic",
  "config": {
    "teams": [...],
    "players": [...],
    "match_format": "mr12",
    "map": "de_mirage",
    "rollback_probability": 0.05
  }
}

Response: Stream of log events or downloadable file
```

#### Demo Parsing
```
POST /api/parse
{
  "mode": "demo",
  "demo_url": "https://example.com/match.dem",
  // OR
  "demo_base64": "...",
  "output_format": "http_log"
}

Response: Parsed log events
```

#### Configuration
```
GET /api/config/templates
Response: Available match templates

POST /api/config/validate
Body: Configuration object
Response: Validation results

GET /api/config/maps
Response: Available maps list

GET /api/config/weapons
Response: Weapons and equipment data
```

#### Match Management
```
GET /api/matches
Response: List of generated/parsed matches

GET /api/matches/{id}
Response: Match details and logs

DELETE /api/matches/{id}
Response: Deletion confirmation

GET /api/matches/{id}/download
Response: Download match logs
```

#### Real-time Streaming
```
WebSocket: /api/stream
- Real-time log event streaming
- Configurable event rate
- Pause/resume support
```

## Web UI Specifications

### Overview
A modern, responsive web interface for configuring and generating CS2 match logs. The UI provides intuitive controls for both synthetic generation and demo file parsing, with real-time preview and export capabilities.

### Technology Stack
- **Frontend Framework**: React 19 with TypeScript
- **UI Components**: shadcn/ui (Radix UI + Tailwind)
- **Styling**: Tailwind CSS 4 with Oxide engine
- **State Management**: 
  - Zustand for client state
  - Tanstack Query for server state
- **Build Tool**: Vite 5+
- **Charts/Visualization**: Tremor or Recharts
- **WebSocket Client**: Socket.io-client
- **Icons**: Lucide React
- **Forms**: React Hook Form + Zod validation
- **Tables**: Tanstack Table
- **Animations**: Framer Motion

### UI Architecture

#### Layout Structure
```
┌────────────────────────────────────────┐
│          Header / Navigation           │
├────────────┬───────────────────────────┤
│            │                           │
│  Sidebar   │      Main Content         │
│            │                           │
│ - Generate │   Match Configuration     │
│ - Parse    │        OR                 │
│ - History  │   Demo Upload/Parse       │
│ - Settings │        OR                 │
│            │   Match History/Logs      │
│            │                           │
└────────────┴───────────────────────────┘
```

### Core UI Features

#### 1. Match Configuration Interface

##### Team Setup (shadcn/ui Implementation)
```typescript
// Visual Team Builder with shadcn/ui
<Tabs defaultValue="team1" className="w-full">
  <TabsList className="grid w-full grid-cols-2">
    <TabsTrigger value="team1">Team 1</TabsTrigger>
    <TabsTrigger value="team2">Team 2</TabsTrigger>
  </TabsList>
  <TabsContent value="team1">
    <Card>
      <CardHeader>
        <CardTitle>Team Configuration</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <Input placeholder="Team Name" />
          <Input placeholder="Team Tag" />
        </div>
        <Select>
          <SelectTrigger>
            <SelectValue placeholder="Select Country" />
          </SelectTrigger>
          <SelectContent>
            {/* Country options */}
          </SelectContent>
        </Select>
        <RadioGroup defaultValue="ct">
          <div className="flex items-center space-x-2">
            <RadioGroupItem value="ct" id="ct" />
            <Label htmlFor="ct">Counter-Terrorist</Label>
          </div>
          <div className="flex items-center space-x-2">
            <RadioGroupItem value="t" id="t" />
            <Label htmlFor="t">Terrorist</Label>
          </div>
        </RadioGroup>
      </CardContent>
    </Card>
  </TabsContent>
</Tabs>
```

##### Player Configuration
- **Player Cards with shadcn/ui**
  - Avatar component with upload
  - Form with validation (React Hook Form + Zod)
  - Slider component for skill rating
  - Select component for role
  - Command palette for weapon selection

##### Match Settings
- **Format Selection with Toggle Group**
  ```typescript
  <ToggleGroup type="single" defaultValue="mr12">
    <ToggleGroupItem value="mr12">MR12</ToggleGroupItem>
    <ToggleGroupItem value="mr15">MR15</ToggleGroupItem>
    <ToggleGroupItem value="custom">Custom</ToggleGroupItem>
  </ToggleGroup>
  ```

- **Map Selection Grid**
  - Responsive grid with Tailwind 4
  - Dialog for map veto simulation
  - Command palette for custom maps

##### Advanced Options
- **Simulation Parameters with shadcn/ui**
  - Slider components for numeric values
  - Switch components for toggles
  - Popover with info tooltips
  - Accordion for grouped settings

#### 2. Demo Parser Interface

##### Upload Section (shadcn/ui Implementation)
```typescript
// Drag & Drop with shadcn/ui and Tailwind 4
<Card className="border-2 border-dashed">
  <CardContent className="flex flex-col items-center justify-center py-10">
    <Upload className="h-10 w-10 text-muted-foreground mb-4" />
    <div className="text-center">
      <p className="text-sm font-medium">Drop demo files here</p>
      <p className="text-xs text-muted-foreground">or click to browse</p>
    </div>
    <Input type="file" accept=".dem" className="hidden" />
    <Button variant="secondary" className="mt-4">
      Browse Files
    </Button>
  </CardContent>
</Card>

// File queue with DataTable
<DataTable columns={fileColumns} data={uploadQueue} />
```

##### Parsing Options
```typescript
// Output Configuration with shadcn/ui
<Card>
  <CardHeader>
    <CardTitle>Parsing Options</CardTitle>
  </CardHeader>
  <CardContent className="space-y-4">
    <Select>
      <SelectTrigger>
        <SelectValue placeholder="Output Format" />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="http">HTTP Logs</SelectItem>
        <SelectItem value="json">JSON</SelectItem>
        <SelectItem value="csv">CSV</SelectItem>
      </SelectContent>
    </Select>
    
    <div className="space-y-2">
      <Label>Event Filters</Label>
      <div className="grid grid-cols-2 gap-2">
        <Checkbox id="kills" label="Kills" />
        <Checkbox id="rounds" label="Rounds" />
        <Checkbox id="economy" label="Economy" />
        <Checkbox id="grenades" label="Grenades" />
      </div>
    </div>
  </CardContent>
</Card>
```

##### Processing Status
```typescript
// Real-time Progress with shadcn/ui
<Alert>
  <AlertDescription>
    <div className="space-y-2">
      <div className="flex justify-between text-sm">
        <span>Parsing demo...</span>
        <span>{progress}%</span>
      </div>
      <Progress value={progress} />
      <div className="flex justify-between text-xs text-muted-foreground">
        <span>{eventsExtracted} events</span>
        <span>~{timeRemaining}s remaining</span>
      </div>
    </div>
  </AlertDescription>
</Alert>
```

#### 3. Real-time Generation Monitor

##### Live Preview
- **Log Stream Display**
  - Syntax-highlighted log viewer
  - Auto-scroll with pause option
  - Search and filter bar
  - Event type color coding

##### Match Visualization
- **Mini-map View**
  - Real-time player positions
  - Kill feed overlay
  - Bomb/objective indicators
  - Round timer display

##### Statistics Dashboard
- **Live Stats**
  - Scoreboard with K/D/A
  - Economy graph
  - Round history timeline
  - Performance metrics

#### 4. Match History & Management

##### History Grid
- **Match List**
  - Sortable columns (date, teams, map, duration)
  - Quick actions (view, download, delete)
  - Status indicators (completed, in-progress, failed)
  - Search and filter options

##### Match Details
- **Detailed View**
  - Full match statistics
  - Round-by-round breakdown
  - Player performance charts
  - Event timeline
  - Log viewer with syntax highlighting

#### 5. Settings & Configuration

##### Global Settings
- **Application Preferences**
  - Theme selector (light/dark/auto)
  - Default output format
  - API endpoint configuration
  - Performance tuning options

##### Templates Management
- **Save/Load Configurations**
  - Template library
  - Import/export JSON configs
  - Share configurations
  - Version control

### User Flows

#### Flow 1: Generate Synthetic Match
1. User clicks "Generate Match" in sidebar
2. Selects or creates teams with players
3. Configures match settings
4. Sets simulation parameters
5. Clicks "Generate" button
6. Views real-time generation
7. Downloads or saves results

#### Flow 2: Parse Demo File
1. User clicks "Parse Demo" in sidebar
2. Drags demo file to upload zone
3. Configures parsing options
4. Initiates parsing
5. Monitors progress
6. Reviews parsed logs
7. Exports results

#### Flow 3: Quick Match
1. User clicks "Quick Match" button
2. Selects preset template
3. Makes minor adjustments
4. Generates match immediately
5. Views results

### UI Components Library

#### shadcn/ui Component Usage
```typescript
// Using shadcn/ui components with Tailwind 4
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet"
import { DataTable } from "@/components/ui/data-table"
import { toast } from "@/components/ui/use-toast"
```

#### Custom Components with Tailwind 4
```typescript
// Team Builder Component
export function TeamBuilder({ team, onChange, playerPool }) {
  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Users className="h-5 w-5" />
          Team Configuration
        </CardTitle>
      </CardHeader>
      <CardContent className="grid gap-4">
        {/* Drag and drop with Tailwind 4's new features */}
        <div className="grid grid-cols-5 gap-2">
          {/* Player slots */}
        </div>
      </CardContent>
    </Card>
  )
}

// Player Card with shadcn/ui
export function PlayerCard({ player, onEdit, onDelete }) {
  return (
    <Card className="group relative hover:shadow-lg transition-shadow">
      <CardContent className="p-4">
        <div className="flex items-center gap-3">
          <Avatar>
            <AvatarImage src={player.avatar} />
            <AvatarFallback>{player.name[0]}</AvatarFallback>
          </Avatar>
          <div className="flex-1">
            <p className="font-medium">{player.name}</p>
            <p className="text-sm text-muted-foreground">{player.steamId}</p>
          </div>
        </div>
        <div className="mt-4 flex gap-2">
          <Button size="sm" variant="outline" onClick={onEdit}>
            <Edit className="h-4 w-4" />
          </Button>
          <Button size="sm" variant="destructive" onClick={onDelete}>
            <Trash className="h-4 w-4" />
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}

// Log Viewer with syntax highlighting
export function LogViewer({ logs, filters }) {
  return (
    <div className="relative rounded-lg border bg-black/95 p-4">
      <ScrollArea className="h-[600px] w-full">
        <pre className="text-sm">
          <code className="language-log">
            {logs.map((log, i) => (
              <div key={i} className="hover:bg-white/5 px-2">
                <span className="text-gray-500">{log.timestamp}</span>
                <span className={cn(
                  "ml-2",
                  log.type === 'kill' && "text-red-400",
                  log.type === 'round_start' && "text-green-400",
                  log.type === 'bomb_planted' && "text-yellow-400"
                )}>
                  {log.content}
                </span>
              </div>
            ))}
          </code>
        </pre>
      </ScrollArea>
    </div>
  )
}
```

### Design System with Tailwind 4

#### Tailwind 4 Configuration
```javascript
// tailwind.config.ts with Oxide engine
export default {
  darkMode: ['class'],
  content: [
    './src/**/*.{ts,tsx}',
    './components/**/*.{ts,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        // CS2 themed colors
        'cs-orange': '#DE9B35',
        'cs-blue': '#5E98D9',
        'ct-blue': '#5E98D9',
        't-orange': '#DE9B35',
      },
      animation: {
        'pulse-slow': 'pulse 4s cubic-bezier(0.4, 0, 0.6, 1) infinite',
      },
    },
  },
  plugins: [
    require('@tailwindcss/container-queries'),
    require('tailwindcss-animate'),
  ],
}
```

#### CSS Variables for Theming
```css
/* Using Tailwind 4's improved CSS variable support */
@layer base {
  :root {
    --background: 0 0% 100%;
    --foreground: 222.2 84% 4.9%;
    --card: 0 0% 100%;
    --card-foreground: 222.2 84% 4.9%;
    --primary: 222.2 47.4% 11.2%;
    --primary-foreground: 210 40% 98%;
    /* CS2 specific */
    --cs-ct: 206 63% 61%;
    --cs-t: 33 63% 54%;
  }
  
  .dark {
    --background: 222.2 84% 4.9%;
    --foreground: 210 40% 98%;
    --card: 222.2 84% 4.9%;
    --card-foreground: 210 40% 98%;
  }
}
```

### Responsive Design

#### Breakpoints (Tailwind 4)
- **Desktop**: `2xl:` 1536px+ (full features)
- **Laptop**: `xl:` 1280px+ (standard layout)
- **Tablet**: `lg:` 1024px+ (condensed sidebar)
- **Mobile**: `<lg:` <1024px (mobile-optimized)

#### Container Queries (New in Tailwind 4)
```typescript
// Responsive components using container queries
<div className="@container">
  <div className="@sm:grid-cols-2 @lg:grid-cols-3 @xl:grid-cols-5 grid gap-4">
    {/* Player cards adapt to container size */}
  </div>
</div>
```

#### Mobile Adaptations
- Sheet component for mobile navigation
- Drawer for mobile settings
- Swipeable tabs for team switching
- Touch-optimized button sizes
- Responsive typography scaling

### Accessibility Features

#### WCAG 2.1 Compliance
- Keyboard navigation support
- Screen reader compatibility
- High contrast mode
- Focus indicators
- Alt text for images
- ARIA labels

#### User Preferences
- Font size adjustment
- Color blind modes
- Reduced motion option
- Keyboard shortcuts
- Tooltip delays

### Performance Optimization

#### Frontend Optimization
- Lazy loading for components
- Virtual scrolling for long lists
- Debounced search inputs
- Memoized calculations
- Code splitting by route

#### Data Management
- Pagination for match history
- Incremental log loading
- WebSocket connection pooling
- Local storage caching
- Optimistic UI updates

### Integration with Backend

#### API Communication with Tanstack Query
```typescript
// Using Tanstack Query for server state
const { data, error, isLoading } = useQuery({
  queryKey: ['matches'],
  queryFn: fetchMatches,
  staleTime: 5 * 60 * 1000, // 5 minutes
})

// Mutations with optimistic updates
const mutation = useMutation({
  mutationFn: generateMatch,
  onMutate: async (newMatch) => {
    await queryClient.cancelQueries({ queryKey: ['matches'] })
    const previousMatches = queryClient.getQueryData(['matches'])
    queryClient.setQueryData(['matches'], old => [...old, newMatch])
    return { previousMatches }
  },
  onError: (err, newMatch, context) => {
    queryClient.setQueryData(['matches'], context.previousMatches)
  },
  onSettled: () => {
    queryClient.invalidateQueries({ queryKey: ['matches'] })
  },
})
```

#### State Management with Zustand
```typescript
// Client state with Zustand
interface AppStore {
  // State
  configuration: MatchConfig
  ui: {
    theme: 'light' | 'dark' | 'system'
    sidebarOpen: boolean
    activeTab: string
  }
  
  // Actions
  setConfiguration: (config: Partial<MatchConfig>) => void
  toggleSidebar: () => void
  setTheme: (theme: string) => void
}

const useAppStore = create<AppStore>((set) => ({
  configuration: defaultConfig,
  ui: {
    theme: 'system',
    sidebarOpen: true,
    activeTab: 'generate',
  },
  
  setConfiguration: (config) =>
    set((state) => ({
      configuration: { ...state.configuration, ...config },
    })),
    
  toggleSidebar: () =>
    set((state) => ({
      ui: { ...state.ui, sidebarOpen: !state.ui.sidebarOpen },
    })),
    
  setTheme: (theme) =>
    set((state) => ({
      ui: { ...state.ui, theme },
    })),
}))
```

#### React 19 Features Usage
```typescript
// Using React 19's use() hook for data fetching
function MatchDetails({ matchId }) {
  const match = use(fetchMatch(matchId))
  
  return (
    <Card>
      <CardHeader>
        <CardTitle>{match.title}</CardTitle>
      </CardHeader>
      <CardContent>
        {/* Match details */}
      </CardContent>
    </Card>
  )
}

// Server Components for initial data
export default async function MatchesPage() {
  const matches = await getMatches()
  
  return (
    <div className="container mx-auto">
      <MatchList matches={matches} />
    </div>
  )
}
```

### Deployment Considerations

#### Build Configuration
- Production builds with minification
- Environment-specific configs
- Docker containerization
- CDN asset delivery
- Service worker for offline support

#### Monitoring
- Error tracking (Sentry)
- Analytics integration
- Performance monitoring
- User session recording
- A/B testing support

## Data Models

### Team Model
```go
type Team struct {
    Name        string
    Tag         string
    Country     string
    Ranking     int
    Players     []Player
    Side        string // "CT" or "TERRORIST"
    Score       int
    Economy     TeamEconomy
}
```

### Player Model
```go
type Player struct {
    Name        string
    SteamID     string
    Team        string
    Stats       PlayerStats
    Economy     PlayerEconomy
    Position    Vector3
    Health      int
    Armor       int
    Helmet      bool
    Kit         bool
    Weapons     []Weapon
}
```

### Event Models
```go
type GameEvent interface {
    GetTimestamp() time.Time
    GetType() string
    ToLogLine() string
}

type KillEvent struct {
    Timestamp    time.Time
    Attacker     *Player
    Victim       *Player
    Weapon       string
    Headshot     bool
    Penetrated   int
    AttackerPos  Vector3
    VictimPos    Vector3
}

type RoundEvent struct {
    Timestamp    time.Time
    RoundNumber  int
    Winner       string
    Reason       string // "elimination", "bomb", "time", "defuse"
    TeamScores   map[string]int
}
```

### Match Configuration
```yaml
match:
  format: mr12
  overtime: mr3
  map: de_mirage
  
teams:
  - name: "Team Alpha"
    players:
      - name: "Player1"
        steam_id: "STEAM_1:0:12345678"
        skill_rating: 2000
        
  - name: "Team Beta"
    players:
      - name: "Player6"
        steam_id: "STEAM_1:0:87654321"
        skill_rating: 1950

simulation:
  rollback_probability: 0.05
  network_issues: true
  realistic_economy: true
  skill_variance: 0.15
```

## Technical Architecture

### Component Architecture
```
┌─────────────────────────────────────┐
│      Web UI (React 19)              │
│  shadcn/ui + Tailwind CSS 4        │
└──────────────┬──────────────────────┘
               │ HTTP/WebSocket
┌──────────────▼──────────────────────┐
│         HTTP API Layer              │
│      (Gin Web Framework)            │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│       Core Engine                   │
├─────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐│
│  │  Synthetic   │  │    Demo      ││
│  │  Generator   │  │   Parser     ││
│  └──────────────┘  └──────────────┘│
└─────────────────────────────────────┘
               │
┌──────────────▼──────────────────────┐
│      Event Processing               │
│   (Format, Filter, Transform)       │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│         Output Layer                │
│   (HTTP, File, Stream, WebSocket)   │
└─────────────────────────────────────┘
```

### Module Structure
```
/backend
  /cmd
    /server       - HTTP server entry point
    /cli         - CLI tool for local generation
  /pkg
    /generator   - Synthetic log generation
    /parser      - Demo file parsing
    /models      - Data structures
    /events      - Event definitions
    /economy     - Economy simulation
    /rollback    - Rollback system
    /api         - API handlers
    /websocket   - WebSocket handlers
  
/frontend
  /src
    /components  - React components
    /pages       - Page components
    /hooks       - Custom React hooks
    /store       - State management
    /services    - API client services
    /utils       - Utility functions
    /styles      - Global styles
  /public        - Static assets
  
/shared
  /types         - Shared TypeScript types
  /constants     - Shared constants
  
/configs         - Sample configurations
/templates       - Match templates
/test           - Unit and integration tests
/docs           - Documentation
```

## Configuration Options

### Global Settings
- Log format (standard, JSON, custom)
- Timestamp format and timezone
- Output verbosity levels
- Performance tuning (concurrency, buffer sizes)

### Match Settings
- Team compositions and skill levels
- Map pool and veto process
- Round and overtime configuration
- Economy parameters
- Rollback frequency and triggers

### Simulation Parameters
- Skill variance and randomness
- Network condition simulation
- Anti-cheat event frequency
- Chat message generation
- Spectator events

## Integration Requirements

### Compatibility
- HLTV log format compatibility
- Support for common parsers (HLstatsX, etc.)
- WebSocket streaming for real-time applications
- Batch export for analysis tools

### Performance Metrics
- Generation speed: 1000+ events/second
- Demo parsing: 25+ minutes of gameplay/second
- Memory usage: <500MB for typical match
- Concurrent operations: 10+ simultaneous generations

## Development Phases

### Phase 1: Core Infrastructure (Weeks 1-2)
- Project setup and structure
- Basic data models
- HTTP server framework
- Configuration system
- Frontend project initialization
- Component library setup

### Phase 2: Backend Core (Weeks 3-4)
- Event generation engine
- Economy simulation
- Basic match flow
- Log formatting
- API endpoint implementation
- WebSocket server setup

### Phase 3: UI Development (Weeks 5-6)
- Match configuration interface
- Team and player builders
- Demo upload interface
- Real-time log viewer
- Basic navigation and layout

### Phase 4: Demo Parsing (Weeks 7-8)
- Demo file reader integration
- Event extraction
- Format conversion
- Performance optimization
- UI parsing progress display

### Phase 5: Advanced Features (Weeks 9-10)
- Rollback system
- Network simulation
- Live match visualization
- Statistics dashboard
- Match history management

### Phase 6: Integration & Polish (Weeks 11-12)
- Frontend-backend integration
- WebSocket streaming
- Responsive design implementation
- Accessibility features
- Performance optimization

### Phase 7: Testing & Deployment (Weeks 13-14)
- Comprehensive testing
- UI/UX testing
- Documentation
- Docker containerization
- Deployment guides

## Success Metrics

### Functional Metrics
- All CS2 event types supported
- Demo parsing accuracy >99%
- Log format validation passing
- Rollback scenarios working correctly
- UI supports all configuration options
- Real-time visualization working smoothly

### Performance Metrics
- Generation speed meeting targets
- Memory usage within limits
- Concurrent operation support
- API response times <100ms
- UI load time <2 seconds
- 60 FPS for live visualizations

### Quality Metrics
- Unit test coverage >80%
- Integration tests passing
- No critical bugs in production
- Documentation completeness
- Accessibility WCAG 2.1 AA compliance
- Mobile responsiveness on all devices

### User Experience Metrics
- Task completion rate >95%
- Average time to generate match <2 minutes
- User satisfaction score >4.5/5
- Support ticket rate <5%

## Risk Mitigation

### Technical Risks
- **Demo format changes**: Maintain version compatibility
- **Performance bottlenecks**: Profile and optimize critical paths
- **Memory leaks**: Implement proper resource management

### Operational Risks
- **Invalid configurations**: Robust validation and defaults
- **Large file handling**: Streaming and chunked processing
- **Concurrent access**: Thread-safe operations

## Future Enhancements

### Planned Features
- CS:GO legacy support
- Advanced statistics generation
- Machine learning for realistic behavior
- Cloud deployment with auto-scaling
- Mobile app for remote generation
- AI-powered match prediction
- Voice commentary generation
- Heat map generation from logs

### Potential Integrations
- Direct HLTV integration
- Faceit/ESEA compatibility
- Tournament system support
- Twitch/YouTube streaming integration
- Discord bot for match generation
- Slack notifications for completed matches
- Export to popular analytics platforms
- Integration with CS2 workshop maps

### UI Enhancements
- 3D match replay viewer
- Advanced data visualization dashboards
- Collaborative team building
- Social sharing features
- Custom themes and skins
- Keyboard-only navigation mode
- VR/AR match visualization
- Multi-language support

## Conclusion

The CS2 HTTP Log Generator provides a comprehensive solution for generating and parsing Counter-Strike 2 match logs with an intuitive web-based user interface. By supporting both synthetic generation and demo file parsing, the tool serves diverse use cases from testing log processors to analyzing recorded matches. 

The modern React-based UI makes the tool accessible to users of all technical levels, while the powerful backend ensures high performance and accuracy. The flexible configuration system, real-time visualization capabilities, and robust feature set ensure the tool can simulate realistic match scenarios while maintaining compatibility with existing CS2 infrastructure.

With its dual-mode operation, comprehensive event coverage, and user-friendly interface, this tool will become an essential resource for CS2 server operators, tournament organizers, developers, and analysts who need to work with CS2 match data in HTTP log format.