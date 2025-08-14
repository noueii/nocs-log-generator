"use client"

import React, { useState } from 'react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { 
  Search, 
  Filter,
  X,
  Download,
  RefreshCw,
  Eye,
  EyeOff
} from 'lucide-react'
import type { TEventType, IEventFilter } from '@/types/events'
import { EVENT_COLORS, getEventIcon } from '@/types/events'

interface EventFilterProps {
  onFilterChange: (filter: IEventFilter) => void
  totalEvents: number
  filteredEvents: number
  className?: string
}

interface FilterPreset {
  id: string
  name: string
  description: string
  types: TEventType[]
  icon: string
}

const EVENT_TYPE_GROUPS = {
  combat: {
    name: 'Combat Events',
    icon: '⚔️',
    types: ['player_death', 'player_hurt', 'weapon_fire'] as TEventType[]
  },
  objectives: {
    name: 'Objective Events', 
    icon: '💣',
    types: ['bomb_planted', 'bomb_defused', 'bomb_exploded'] as TEventType[]
  },
  rounds: {
    name: 'Round Events',
    icon: '🏁', 
    types: ['round_start', 'round_end'] as TEventType[]
  },
  economy: {
    name: 'Economy Events',
    icon: '💰',
    types: ['item_purchase'] as TEventType[]
  },
  utility: {
    name: 'Utility Events',
    icon: '🎯',
    types: ['grenade_thrown', 'flashbang_detonate'] as TEventType[]
  },
  social: {
    name: 'Social Events',
    icon: '💬',
    types: ['chat_message', 'player_connect', 'player_disconnect'] as TEventType[]
  },
  system: {
    name: 'System Events',
    icon: '⚙️',
    types: ['server_command', 'team_switch'] as TEventType[]
  }
}

const FILTER_PRESETS: FilterPreset[] = [
  {
    id: 'all',
    name: 'All Events',
    description: 'Show all log events',
    types: Object.values(EVENT_TYPE_GROUPS).flatMap(group => group.types),
    icon: '📝'
  },
  {
    id: 'kills-only',
    name: 'Kills Only',
    description: 'Show only kill events',
    types: ['player_death'],
    icon: '☠️'
  },
  {
    id: 'important',
    name: 'Important Events',
    description: 'Kills, rounds, and bomb events',
    types: ['player_death', 'round_start', 'round_end', 'bomb_planted', 'bomb_defused', 'bomb_exploded'],
    icon: '⭐'
  },
  {
    id: 'no-spam',
    name: 'No Spam',
    description: 'Hide weapon fire and hurt events',
    types: Object.values(EVENT_TYPE_GROUPS)
      .flatMap(group => group.types)
      .filter(type => !['weapon_fire', 'player_hurt'].includes(type)),
    icon: '🔇'
  },
  {
    id: 'economy-focus',
    name: 'Economy Focus',
    description: 'Purchases, kills, and round events',
    types: ['item_purchase', 'player_death', 'round_start', 'round_end'],
    icon: '💎'
  }
]

export function EventFilter({
  onFilterChange,
  totalEvents,
  filteredEvents,
  className
}: EventFilterProps) {
  const [searchTerm, setSearchTerm] = useState('')
  const [selectedTypes, setSelectedTypes] = useState<Set<TEventType>>(
    new Set(Object.values(EVENT_TYPE_GROUPS).flatMap(group => group.types))
  )
  const [selectedPlayers, setSelectedPlayers] = useState<string[]>([])
  const [selectedRounds, setSelectedRounds] = useState<number[]>([])
  const [tickRange, setTickRange] = useState<{ start?: number; end?: number }>({})
  const [activePreset, setActivePreset] = useState<string>('all')
  const [collapsed, setCollapsed] = useState(false)

  // Apply filters whenever selections change
  React.useEffect(() => {
    const filter: IEventFilter = {
      types: selectedTypes.size > 0 ? Array.from(selectedTypes) : undefined,
      players: selectedPlayers.length > 0 ? selectedPlayers : undefined,
      rounds: selectedRounds.length > 0 ? selectedRounds : undefined,
      startTick: tickRange.start,
      endTick: tickRange.end
    }

    onFilterChange(filter)
  }, [selectedTypes, selectedPlayers, selectedRounds, tickRange, onFilterChange])

  const handlePresetChange = (preset: FilterPreset) => {
    setActivePreset(preset.id)
    setSelectedTypes(new Set(preset.types))
    
    // Clear other filters when applying preset
    if (preset.id === 'all') {
      setSelectedPlayers([])
      setSelectedRounds([])
      setTickRange({})
    }
  }

  const handleTypeToggle = (type: TEventType, checked: boolean) => {
    setActivePreset('custom')
    const newSelected = new Set(selectedTypes)
    if (checked) {
      newSelected.add(type)
    } else {
      newSelected.delete(type)
    }
    setSelectedTypes(newSelected)
  }

  const handleGroupToggle = (groupTypes: TEventType[], checked: boolean) => {
    setActivePreset('custom')
    const newSelected = new Set(selectedTypes)
    groupTypes.forEach(type => {
      if (checked) {
        newSelected.add(type)
      } else {
        newSelected.delete(type)
      }
    })
    setSelectedTypes(newSelected)
  }

  const clearAllFilters = () => {
    setActivePreset('all')
    setSelectedTypes(new Set(Object.values(EVENT_TYPE_GROUPS).flatMap(group => group.types)))
    setSelectedPlayers([])
    setSelectedRounds([])
    setTickRange({})
    setSearchTerm('')
  }

  const exportFilters = () => {
    const filterConfig = {
      preset: activePreset,
      types: Array.from(selectedTypes),
      players: selectedPlayers,
      rounds: selectedRounds,
      tickRange,
      searchTerm
    }
    
    const blob = new Blob([JSON.stringify(filterConfig, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'event-filters.json'
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  }

  if (collapsed) {
    return (
      <Card className={cn('border-gray-700', className)}>
        <div className="flex items-center justify-between p-3">
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setCollapsed(false)}
              className="h-6 w-6 p-0"
            >
              <Eye className="h-4 w-4" />
            </Button>
            <Filter className="h-4 w-4 text-muted-foreground" />
            <span className="text-sm text-muted-foreground">
              Showing {filteredEvents} / {totalEvents} events
            </span>
          </div>
          
          <div className="flex items-center gap-1">
            {activePreset !== 'all' && (
              <Badge variant="secondary" className="text-xs">
                {FILTER_PRESETS.find(p => p.id === activePreset)?.name || 'Custom'}
              </Badge>
            )}
          </div>
        </div>
      </Card>
    )
  }

  return (
    <Card className={cn('border-gray-700', className)}>
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2 text-base">
            <Filter className="h-4 w-4" />
            Event Filters
            <Badge variant="outline" className="ml-2">
              {filteredEvents} / {totalEvents}
            </Badge>
          </CardTitle>
          
          <div className="flex items-center gap-1">
            <Button
              variant="ghost"
              size="sm"
              onClick={exportFilters}
              className="h-6 px-2 text-xs"
            >
              <Download className="h-3 w-3" />
            </Button>
            
            <Button
              variant="ghost"  
              size="sm"
              onClick={clearAllFilters}
              className="h-6 px-2 text-xs"
            >
              <RefreshCw className="h-3 w-3" />
            </Button>
            
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setCollapsed(true)}
              className="h-6 w-6 p-0"
            >
              <EyeOff className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </CardHeader>

      <CardContent className="space-y-4">
        {/* Search */}
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Search event text..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="pl-9"
          />
          {searchTerm && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setSearchTerm('')}
              className="absolute right-1 top-1/2 transform -translate-y-1/2 h-6 w-6 p-0"
            >
              <X className="h-3 w-3" />
            </Button>
          )}
        </div>

        {/* Quick Filter Presets */}
        <div className="space-y-2">
          <h4 className="text-sm font-medium">Quick Filters:</h4>
          <div className="flex flex-wrap gap-1">
            {FILTER_PRESETS.map(preset => (
              <Button
                key={preset.id}
                variant={activePreset === preset.id ? "default" : "outline"}
                size="sm"
                onClick={() => handlePresetChange(preset)}
                className={cn(
                  'h-7 px-2 text-xs',
                  activePreset === preset.id && 'bg-blue-600 hover:bg-blue-700'
                )}
              >
                <span className="mr-1">{preset.icon}</span>
                {preset.name}
              </Button>
            ))}
          </div>
        </div>

        {/* Event Type Groups */}
        <div className="space-y-3">
          <h4 className="text-sm font-medium">Event Types:</h4>
          
          {Object.entries(EVENT_TYPE_GROUPS).map(([groupId, group]) => {
            const groupSelected = group.types.every(type => selectedTypes.has(type))
            const groupPartial = group.types.some(type => selectedTypes.has(type)) && !groupSelected
            
            return (
              <div key={groupId} className="space-y-2">
                {/* Group header */}
                <div className="flex items-center gap-2">
                  <Checkbox
                    checked={groupSelected}
                    indeterminate={groupPartial}
                    onCheckedChange={(checked) => 
                      handleGroupToggle(group.types, checked as boolean)
                    }
                  />
                  <span className="text-sm font-medium flex items-center gap-1">
                    <span>{group.icon}</span>
                    {group.name}
                  </span>
                  <Badge variant="secondary" className="text-xs">
                    {group.types.filter(type => selectedTypes.has(type)).length}/{group.types.length}
                  </Badge>
                </div>

                {/* Individual event types */}
                <div className="ml-6 grid grid-cols-2 gap-1">
                  {group.types.map(type => (
                    <div key={type} className="flex items-center gap-2">
                      <Checkbox
                        checked={selectedTypes.has(type)}
                        onCheckedChange={(checked) => 
                          handleTypeToggle(type, checked as boolean)
                        }
                      />
                      <span className="text-xs flex items-center gap-1">
                        <span style={{ color: EVENT_COLORS[type] }}>
                          {getEventIcon(type)}
                        </span>
                        {type.replace('_', ' ')}
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            )
          })}
        </div>

        {/* Advanced Filters */}
        <div className="space-y-3 pt-2 border-t border-gray-200">
          <h4 className="text-sm font-medium">Advanced Filters:</h4>
          
          {/* Tick Range */}
          <div className="grid grid-cols-2 gap-2">
            <div>
              <label className="text-xs text-muted-foreground">Start Tick</label>
              <Input
                type="number"
                placeholder="0"
                value={tickRange.start || ''}
                onChange={(e) => setTickRange(prev => ({
                  ...prev,
                  start: e.target.value ? parseInt(e.target.value) : undefined
                }))}
                className="text-xs"
              />
            </div>
            <div>
              <label className="text-xs text-muted-foreground">End Tick</label>
              <Input
                type="number" 
                placeholder="999999"
                value={tickRange.end || ''}
                onChange={(e) => setTickRange(prev => ({
                  ...prev,
                  end: e.target.value ? parseInt(e.target.value) : undefined
                }))}
                className="text-xs"
              />
            </div>
          </div>

          {/* Round Filter */}
          <div>
            <label className="text-xs text-muted-foreground">Specific Rounds (comma-separated)</label>
            <Input
              placeholder="1, 5, 12, 16"
              value={selectedRounds.join(', ')}
              onChange={(e) => {
                const rounds = e.target.value
                  .split(',')
                  .map(r => parseInt(r.trim()))
                  .filter(r => !isNaN(r))
                setSelectedRounds(rounds)
              }}
              className="text-xs"
            />
          </div>

          {/* Player Filter */}
          <div>
            <label className="text-xs text-muted-foreground">Player Names (comma-separated)</label>
            <Input
              placeholder="player1, player2"
              value={selectedPlayers.join(', ')}
              onChange={(e) => {
                const players = e.target.value
                  .split(',')
                  .map(p => p.trim())
                  .filter(p => p.length > 0)
                setSelectedPlayers(players)
              }}
              className="text-xs"
            />
          </div>
        </div>

        {/* Filter Summary */}
        {(selectedPlayers.length > 0 || selectedRounds.length > 0 || tickRange.start || tickRange.end) && (
          <div className="pt-2 border-t border-gray-200">
            <h4 className="text-xs text-muted-foreground mb-2">Active Advanced Filters:</h4>
            <div className="flex flex-wrap gap-1">
              {selectedPlayers.map(player => (
                <Badge key={player} variant="secondary" className="text-xs">
                  Player: {player}
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setSelectedPlayers(prev => prev.filter(p => p !== player))}
                    className="ml-1 h-3 w-3 p-0 hover:bg-transparent"
                  >
                    <X className="h-2 w-2" />
                  </Button>
                </Badge>
              ))}
              
              {selectedRounds.map(round => (
                <Badge key={round} variant="secondary" className="text-xs">
                  Round: {round}
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setSelectedRounds(prev => prev.filter(r => r !== round))}
                    className="ml-1 h-3 w-3 p-0 hover:bg-transparent"
                  >
                    <X className="h-2 w-2" />
                  </Button>
                </Badge>
              ))}
              
              {(tickRange.start || tickRange.end) && (
                <Badge variant="secondary" className="text-xs">
                  Ticks: {tickRange.start || '0'} - {tickRange.end || '∞'}
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setTickRange({})}
                    className="ml-1 h-3 w-3 p-0 hover:bg-transparent"
                  >
                    <X className="h-2 w-2" />
                  </Button>
                </Badge>
              )}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}

export default EventFilter