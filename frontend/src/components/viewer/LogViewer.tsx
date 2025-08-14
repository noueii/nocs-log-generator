"use client"

import React, { useEffect, useRef, useState, useCallback } from 'react'
import { cn } from '@/lib/utils'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { ScrollArea } from '@/components/ui/scroll-area'
import { 
  Pause, 
  Play, 
  Square, 
  Download, 
  Maximize2, 
  Minimize2,
  Search,
  Filter
} from 'lucide-react'
import type { 
  ISpecificGameEvent, 
  TEventType,
  IKillEvent,
  IRoundStartEvent,
  IRoundEndEvent,
  IBombPlantEvent,
  IBombDefuseEvent,
  IPlayerHurtEvent,
  IChatEvent,
  IItemPurchaseEvent
} from '@/types/events'
import { EVENT_COLORS, formatEventForDisplay, getEventIcon } from '@/types/events'

interface LogViewerProps {
  events: ISpecificGameEvent[]
  isStreaming?: boolean
  onPause?: () => void
  onResume?: () => void
  onStop?: () => void
  autoScroll?: boolean
  className?: string
  showControls?: boolean
  fullscreen?: boolean
  onToggleFullscreen?: () => void
}

interface LogLine {
  id: string
  timestamp: string
  tick: number
  round: number
  event: ISpecificGameEvent
  formattedText: string
  level: 'info' | 'warning' | 'error' | 'success'
  highlighted?: boolean
}

const EVENT_LOG_LEVELS: Record<TEventType, 'info' | 'warning' | 'error' | 'success'> = {
  'player_death': 'error',
  'round_start': 'success',
  'round_end': 'success',
  'bomb_planted': 'warning',
  'bomb_defused': 'success',
  'bomb_exploded': 'error',
  'player_hurt': 'warning',
  'player_connect': 'info',
  'player_disconnect': 'info',
  'item_purchase': 'info',
  'grenade_thrown': 'info',
  'weapon_fire': 'info',
  'flashbang_detonate': 'warning',
  'chat_message': 'info',
  'team_switch': 'info',
  'server_command': 'info',
}

export function LogViewer({
  events = [],
  isStreaming = false,
  onPause,
  onResume,
  onStop,
  autoScroll = true,
  className,
  showControls = true,
  fullscreen = false,
  onToggleFullscreen,
}: LogViewerProps) {
  const [isPaused, setIsPaused] = useState(false)
  const [searchTerm, setSearchTerm] = useState('')
  const [showFilters, setShowFilters] = useState(false)
  const [selectedPlayers, setSelectedPlayers] = useState<string[]>([])
  const scrollAreaRef = useRef<HTMLDivElement>(null)
  const endOfLogRef = useRef<HTMLDivElement>(null)

  // Convert events to log lines
  const logLines: LogLine[] = events.map((event, index) => ({
    id: `${event.tick}-${event.type}-${index}`,
    timestamp: formatTimestamp(event.timestamp),
    tick: event.tick,
    round: event.round,
    event,
    formattedText: formatEventForDisplay(event),
    level: EVENT_LOG_LEVELS[event.type],
    highlighted: searchTerm && formatEventForDisplay(event).toLowerCase().includes(searchTerm.toLowerCase())
  }))

  // Filter log lines based on search and selected players
  const filteredLines = logLines.filter(line => {
    if (searchTerm && !line.formattedText.toLowerCase().includes(searchTerm.toLowerCase())) {
      return false
    }
    
    if (selectedPlayers.length > 0) {
      // Check if event involves any selected players
      const eventText = line.formattedText.toLowerCase()
      const hasSelectedPlayer = selectedPlayers.some(player => 
        eventText.includes(player.toLowerCase())
      )
      if (!hasSelectedPlayer) return false
    }
    
    return true
  })

  // Auto-scroll to bottom when new events arrive
  useEffect(() => {
    if (autoScroll && !isPaused && endOfLogRef.current) {
      endOfLogRef.current.scrollIntoView({ behavior: 'smooth' })
    }
  }, [filteredLines.length, autoScroll, isPaused])

  const handlePause = useCallback(() => {
    setIsPaused(true)
    onPause?.()
  }, [onPause])

  const handleResume = useCallback(() => {
    setIsPaused(false)
    onResume?.()
  }, [onResume])

  const handleStop = useCallback(() => {
    setIsPaused(false)
    onStop?.()
  }, [onStop])

  const downloadLogs = useCallback(() => {
    const logText = filteredLines.map(line => 
      `[${line.timestamp}] [TICK:${line.tick}] [R${line.round}] ${line.formattedText}`
    ).join('\n')
    
    const blob = new Blob([logText], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `cs2-match-logs-${new Date().toISOString().slice(0, 19).replace(/:/g, '-')}.txt`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  }, [filteredLines])

  // Get unique player names from events for filtering
  const playerNames = Array.from(new Set(
    events.flatMap(event => {
      const players: string[] = []
      switch (event.type) {
        case 'player_death':
          const killEvent = event as IKillEvent
          players.push(killEvent.attacker.name, killEvent.victim.name)
          if (killEvent.assister) players.push(killEvent.assister.name)
          break
        case 'player_hurt':
          const hurtEvent = event as IPlayerHurtEvent
          players.push(hurtEvent.attacker.name, hurtEvent.victim.name)
          break
        case 'bomb_planted':
        case 'bomb_defused':
          const bombEvent = event as IBombPlantEvent | IBombDefuseEvent
          players.push(bombEvent.player.name)
          break
        case 'item_purchase':
          const purchaseEvent = event as IItemPurchaseEvent
          players.push(purchaseEvent.player.name)
          break
        case 'chat_message':
          const chatEvent = event as IChatEvent
          if (chatEvent.player) players.push(chatEvent.player.name)
          break
      }
      return players
    })
  )).sort()

  return (
    <Card className={cn(
      'bg-gray-900 border-gray-700 text-green-400 font-mono text-sm',
      fullscreen && 'fixed inset-0 z-50 rounded-none',
      className
    )}>
      {/* Header with controls */}
      {showControls && (
        <div className="flex items-center justify-between p-3 border-b border-gray-700 bg-gray-800">
          <div className="flex items-center gap-2">
            <div className="flex items-center gap-1">
              <div className={cn(
                'w-3 h-3 rounded-full',
                isStreaming && !isPaused ? 'bg-green-400 animate-pulse' : 'bg-gray-500'
              )} />
              <span className="text-xs text-gray-300">
                {isStreaming ? (isPaused ? 'PAUSED' : 'STREAMING') : 'STOPPED'}
              </span>
            </div>
            <Badge variant="outline" className="text-xs">
              {filteredLines.length} lines
            </Badge>
          </div>

          <div className="flex items-center gap-2">
            {/* Search */}
            <div className="relative">
              <Search className="absolute left-2 top-1/2 transform -translate-y-1/2 h-3 w-3 text-gray-400" />
              <input
                type="text"
                placeholder="Search logs..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-7 pr-2 py-1 text-xs bg-gray-700 border border-gray-600 rounded focus:ring-1 focus:ring-blue-500 focus:border-blue-500 text-white placeholder-gray-400"
              />
            </div>

            {/* Filter button */}
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setShowFilters(!showFilters)}
              className={cn(
                'h-7 px-2 text-gray-300 hover:bg-gray-700',
                showFilters && 'bg-gray-700'
              )}
            >
              <Filter className="h-3 w-3" />
            </Button>

            {/* Playback controls */}
            {isStreaming ? (
              <div className="flex gap-1">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={isPaused ? handleResume : handlePause}
                  className="h-7 px-2 text-gray-300 hover:bg-gray-700"
                >
                  {isPaused ? <Play className="h-3 w-3" /> : <Pause className="h-3 w-3" />}
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={handleStop}
                  className="h-7 px-2 text-gray-300 hover:bg-gray-700"
                >
                  <Square className="h-3 w-3" />
                </Button>
              </div>
            ) : null}

            {/* Download logs */}
            <Button
              variant="ghost"
              size="sm"
              onClick={downloadLogs}
              className="h-7 px-2 text-gray-300 hover:bg-gray-700"
              disabled={filteredLines.length === 0}
            >
              <Download className="h-3 w-3" />
            </Button>

            {/* Fullscreen toggle */}
            {onToggleFullscreen && (
              <Button
                variant="ghost"
                size="sm"
                onClick={onToggleFullscreen}
                className="h-7 px-2 text-gray-300 hover:bg-gray-700"
              >
                {fullscreen ? <Minimize2 className="h-3 w-3" /> : <Maximize2 className="h-3 w-3" />}
              </Button>
            )}
          </div>
        </div>
      )}

      {/* Filter panel */}
      {showFilters && (
        <div className="p-3 border-b border-gray-700 bg-gray-800">
          <div className="space-y-2">
            <h4 className="text-xs font-medium text-gray-300 mb-2">Filter by Players:</h4>
            <div className="flex flex-wrap gap-1 max-h-20 overflow-y-auto">
              {playerNames.map(player => (
                <Button
                  key={player}
                  variant="ghost"
                  size="sm"
                  onClick={() => {
                    setSelectedPlayers(prev => 
                      prev.includes(player)
                        ? prev.filter(p => p !== player)
                        : [...prev, player]
                    )
                  }}
                  className={cn(
                    'h-6 px-2 text-xs',
                    selectedPlayers.includes(player)
                      ? 'bg-blue-600 text-white hover:bg-blue-700'
                      : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                  )}
                >
                  {player}
                </Button>
              ))}
            </div>
            {selectedPlayers.length > 0 && (
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setSelectedPlayers([])}
                className="h-6 px-2 text-xs text-gray-400 hover:text-white"
              >
                Clear filters
              </Button>
            )}
          </div>
        </div>
      )}

      {/* Log content */}
      <ScrollArea 
        className={cn(
          'flex-1 p-0',
          fullscreen ? 'h-[calc(100vh-140px)]' : 'h-96'
        )}
        ref={scrollAreaRef}
      >
        <div className="p-4 space-y-0.5">
          {filteredLines.length === 0 ? (
            <div className="text-center text-gray-500 py-8">
              {events.length === 0 ? 'No log events yet...' : 'No events match your filters'}
            </div>
          ) : (
            filteredLines.map((line) => (
              <LogLineComponent
                key={line.id}
                line={line}
                highlighted={line.highlighted}
              />
            ))
          )}
          <div ref={endOfLogRef} />
        </div>
      </ScrollArea>
    </Card>
  )
}

interface LogLineProps {
  line: LogLine
  highlighted?: boolean
}

function LogLineComponent({ line, highlighted }: LogLineProps) {
  const { timestamp, tick, round, event, formattedText, level } = line

  // Get color based on event type
  const eventColor = EVENT_COLORS[event.type] || '#6b7280'
  const icon = getEventIcon(event.type)

  // Syntax highlighting
  const highlightedText = highlightEventText(formattedText, event)

  return (
    <div 
      className={cn(
        'flex items-start gap-2 py-0.5 hover:bg-gray-800/50 transition-colors',
        highlighted && 'bg-yellow-900/30'
      )}
    >
      {/* Line number / Timestamp */}
      <div className="flex-shrink-0 text-xs text-gray-500 font-mono w-16 text-right">
        {tick}
      </div>

      {/* Timestamp */}
      <div className="flex-shrink-0 text-xs text-gray-400 font-mono w-24">
        {timestamp}
      </div>

      {/* Round */}
      <div className="flex-shrink-0 text-xs text-blue-400 font-mono w-8">
        R{round}
      </div>

      {/* Event icon */}
      <div className="flex-shrink-0 text-xs w-6">
        <span style={{ color: eventColor }}>
          {icon}
        </span>
      </div>

      {/* Event text */}
      <div 
        className="flex-1 text-sm font-mono leading-tight"
        dangerouslySetInnerHTML={{ __html: highlightedText }}
      />
    </div>
  )
}

function formatTimestamp(timestamp: string): string {
  try {
    const date = new Date(timestamp)
    return date.toLocaleTimeString('en-US', { 
      hour12: false, 
      hour: '2-digit', 
      minute: '2-digit', 
      second: '2-digit' 
    })
  } catch {
    return timestamp.slice(-8) // Fallback to last 8 characters
  }
}

function highlightEventText(text: string, event: ISpecificGameEvent): string {
  let highlighted = text

  // Color player names
  switch (event.type) {
    case 'player_death':
      const killEvent = event as IKillEvent
      highlighted = highlighted
        .replace(new RegExp(killEvent.attacker.name, 'g'), `<span class="text-red-400 font-medium">${killEvent.attacker.name}</span>`)
        .replace(new RegExp(killEvent.victim.name, 'g'), `<span class="text-orange-400 font-medium">${killEvent.victim.name}</span>`)
        .replace(/\(headshot\)/, '<span class="text-red-500 font-bold">(headshot)</span>')
      if (killEvent.assister) {
        highlighted = highlighted.replace(new RegExp(killEvent.assister.name, 'g'), `<span class="text-yellow-400 font-medium">${killEvent.assister.name}</span>`)
      }
      break

    case 'player_hurt':
      const hurtEvent = event as IPlayerHurtEvent
      highlighted = highlighted
        .replace(new RegExp(hurtEvent.attacker.name, 'g'), `<span class="text-orange-400 font-medium">${hurtEvent.attacker.name}</span>`)
        .replace(new RegExp(hurtEvent.victim.name, 'g'), `<span class="text-yellow-400 font-medium">${hurtEvent.victim.name}</span>`)
        .replace(/(\d+) HP/, '<span class="text-red-400 font-medium">$1 HP</span>')
      break

    case 'bomb_planted':
    case 'bomb_defused':
      const bombEvent = event as IBombPlantEvent | IBombDefuseEvent
      highlighted = highlighted
        .replace(new RegExp(bombEvent.player.name, 'g'), `<span class="text-green-400 font-medium">${bombEvent.player.name}</span>`)
        .replace(/site [AB]/, '<span class="text-yellow-400 font-medium">$&</span>')
      break

    case 'round_start':
    case 'round_end':
      highlighted = highlighted
        .replace(/Round \d+/, '<span class="text-blue-400 font-medium">$&</span>')
        .replace(/(CT|TERRORIST)/, '<span class="text-green-400 font-medium">$1</span>')
        .replace(/\(CT: \d+, T: \d+\)/, '<span class="text-cyan-400">$&</span>')
      break

    case 'item_purchase':
      const purchaseEvent = event as IItemPurchaseEvent
      highlighted = highlighted
        .replace(new RegExp(purchaseEvent.player.name, 'g'), `<span class="text-blue-400 font-medium">${purchaseEvent.player.name}</span>`)
        .replace(/\$\d+/, '<span class="text-green-400 font-medium">$&</span>')
      break

    case 'chat_message':
      const chatEvent = event as IChatEvent
      highlighted = highlighted
        .replace(/\[TEAM\]/, '<span class="text-yellow-400 font-medium">[TEAM]</span>')
        .replace(/\[ALL\]/, '<span class="text-blue-400 font-medium">[ALL]</span>')
      if (chatEvent.player) {
        highlighted = highlighted.replace(new RegExp(chatEvent.player.name, 'g'), `<span class="text-cyan-400 font-medium">${chatEvent.player.name}</span>`)
      }
      break
  }

  // Highlight weapon names
  highlighted = highlighted.replace(/(ak47|m4a4|m4a1|awp|deagle|glock|usp|p250|tec9|five-seven|cz75|dual_berettas|mac10|mp9|ump45|p90|bizon|mag7|sawed_off|nova|xm1014|m249|negev|knife|hegrenade|flashbang|smokegrenade|incgrenade|molotov|decoy)/gi, 
    '<span class="text-purple-400 font-medium">$1</span>')

  // Highlight Steam IDs (make them clickable)
  highlighted = highlighted.replace(/(STEAM_[0-5]:[01]:\d+)/g, 
    '<span class="text-blue-300 underline cursor-pointer hover:text-blue-200" title="Click to copy Steam ID">$1</span>')

  return highlighted
}

export default LogViewer