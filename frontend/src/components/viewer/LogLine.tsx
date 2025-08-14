"use client"

import React, { useState, useCallback } from 'react'
import { cn } from '@/lib/utils'
import { Copy, ChevronDown, ChevronRight, Info } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { 
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { toast } from 'sonner'
import type { 
  ISpecificGameEvent,
  IKillEvent,
  IPlayerHurtEvent,
  IBombPlantEvent,
  IBombDefuseEvent,
  IItemPurchaseEvent,
  IChatEvent
} from '@/types/events'
import { EVENT_COLORS, getEventIcon } from '@/types/events'

interface LogLineProps {
  event: ISpecificGameEvent
  lineNumber: number
  highlighted?: boolean
  searchTerm?: string
  showDetails?: boolean
  onToggleDetails?: () => void
  className?: string
}

export function LogLine({
  event,
  lineNumber,
  highlighted = false,
  searchTerm,
  showDetails = false,
  onToggleDetails,
  className
}: LogLineProps) {
  const [isExpanded, setIsExpanded] = useState(false)

  const eventColor = EVENT_COLORS[event.type] || '#6b7280'
  const icon = getEventIcon(event.type)
  
  const handleCopyLine = useCallback(() => {
    const lineText = `[${formatTimestamp(event.timestamp)}] [TICK:${event.tick}] [R${event.round}] ${formatEventText(event)}`
    navigator.clipboard.writeText(lineText)
    toast.success('Log line copied to clipboard')
  }, [event])

  const handleToggleExpanded = useCallback(() => {
    setIsExpanded(prev => !prev)
    onToggleDetails?.()
  }, [onToggleDetails])

  const hasExpandableDetails = hasEventDetails(event)

  return (
    <div 
      className={cn(
        'group relative flex items-start gap-2 py-1 px-2 hover:bg-gray-800/50 transition-colors rounded-sm font-mono text-sm',
        highlighted && 'bg-yellow-900/30 border-l-2 border-yellow-500',
        className
      )}
    >
      {/* Line number */}
      <div className="flex-shrink-0 w-12 text-right text-xs text-gray-500 select-none">
        {lineNumber}
      </div>

      {/* Timestamp */}
      <div className="flex-shrink-0 w-20 text-xs text-gray-400">
        {formatTimestamp(event.timestamp)}
      </div>

      {/* Round */}
      <div className="flex-shrink-0 w-8 text-xs text-blue-400">
        R{event.round}
      </div>

      {/* Tick */}
      <div className="flex-shrink-0 w-16 text-xs text-gray-500">
        {event.tick}
      </div>

      {/* Event type icon */}
      <div className="flex-shrink-0 w-6 text-center">
        <span style={{ color: eventColor }} className="text-base">
          {icon}
        </span>
      </div>

      {/* Event type badge */}
      <div className="flex-shrink-0">
        <Badge 
          variant="outline" 
          className="text-xs px-2 py-0 h-5"
          style={{ 
            borderColor: eventColor + '40',
            backgroundColor: eventColor + '10',
            color: eventColor
          }}
        >
          {event.type.replace('_', ' ')}
        </Badge>
      </div>

      {/* Main event text */}
      <div className="flex-1 min-w-0">
        <div 
          className="leading-tight"
          dangerouslySetInnerHTML={{ 
            __html: highlightEventText(formatEventText(event), searchTerm) 
          }}
        />
        
        {/* Expanded details */}
        {isExpanded && hasExpandableDetails && (
          <div className="mt-2 pl-4 border-l-2 border-gray-700 space-y-1">
            {renderEventDetails(event)}
          </div>
        )}
      </div>

      {/* Action buttons */}
      <div className="flex-shrink-0 flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
        {hasExpandableDetails && (
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={handleToggleExpanded}
                  className="h-6 w-6 p-0 text-gray-400 hover:text-white"
                >
                  {isExpanded ? (
                    <ChevronDown className="h-3 w-3" />
                  ) : (
                    <ChevronRight className="h-3 w-3" />
                  )}
                </Button>
              </TooltipTrigger>
              <TooltipContent>
                <p>{isExpanded ? 'Hide details' : 'Show details'}</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        )}
        
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant="ghost"
                size="sm"
                onClick={handleCopyLine}
                className="h-6 w-6 p-0 text-gray-400 hover:text-white"
              >
                <Copy className="h-3 w-3" />
              </Button>
            </TooltipTrigger>
            <TooltipContent>
              <p>Copy line</p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      </div>
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
    // Fallback to extracting time portion
    return timestamp.slice(-8)
  }
}

function formatEventText(event: ISpecificGameEvent): string {
  switch (event.type) {
    case 'player_death':
      const killEvent = event as IKillEvent
      let killText = `${killEvent.attacker.name} killed ${killEvent.victim.name}`
      if (killEvent.weapon) killText += ` with ${killEvent.weapon}`
      if (killEvent.headshot) killText += ' (headshot)'
      if (killEvent.assister) killText += ` (assist: ${killEvent.assister.name})`
      if (killEvent.distance > 0) killText += ` (${Math.round(killEvent.distance)}u)`
      return killText

    case 'round_start':
      const roundStart = event
      return `Round ${event.round} started`

    case 'round_end':
      const roundEnd = event
      return `Round ${event.round} ended`

    case 'bomb_planted':
      const plantEvent = event as IBombPlantEvent
      return `${plantEvent.player.name} planted the bomb at site ${plantEvent.site}`

    case 'bomb_defused':
      const defuseEvent = event as IBombDefuseEvent
      return `${defuseEvent.player.name} defused the bomb${defuseEvent.with_kit ? ' (with kit)' : ''}`

    case 'bomb_exploded':
      return 'The bomb has exploded'

    case 'player_hurt':
      const hurtEvent = event as IPlayerHurtEvent
      return `${hurtEvent.attacker.name} damaged ${hurtEvent.victim.name} for ${hurtEvent.damage} HP (${hurtEvent.health} HP remaining)`

    case 'item_purchase':
      const purchaseEvent = event as IItemPurchaseEvent
      return `${purchaseEvent.player.name} purchased ${purchaseEvent.item} ($${purchaseEvent.cost})`

    case 'chat_message':
      const chatEvent = event as IChatEvent
      const chatPrefix = chatEvent.team ? '[TEAM]' : '[ALL]'
      return `${chatPrefix} ${chatEvent.player?.name || 'Server'}: ${chatEvent.message}`

    default:
      return `${event.type} event`
  }
}

function highlightEventText(text: string, searchTerm?: string): string {
  if (!searchTerm) return text

  const regex = new RegExp(`(${escapeRegex(searchTerm)})`, 'gi')
  return text.replace(regex, '<mark class="bg-yellow-400/50 text-yellow-100 rounded px-1">$1</mark>')
}

function escapeRegex(text: string): string {
  return text.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}

function hasEventDetails(event: ISpecificGameEvent): boolean {
  const detailedEvents = [
    'player_death',
    'player_hurt', 
    'bomb_planted',
    'bomb_defused',
    'item_purchase'
  ]
  return detailedEvents.includes(event.type)
}

function renderEventDetails(event: ISpecificGameEvent): React.ReactNode {
  switch (event.type) {
    case 'player_death':
      const killEvent = event as IKillEvent
      return (
        <div className="space-y-1 text-xs text-gray-300">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <div className="text-gray-400">Attacker:</div>
              <div className="text-red-400">{killEvent.attacker.name}</div>
              <div className="text-gray-500">Team: {killEvent.attacker.team}</div>
            </div>
            <div>
              <div className="text-gray-400">Victim:</div>
              <div className="text-orange-400">{killEvent.victim.name}</div>
              <div className="text-gray-500">Team: {killEvent.victim.team}</div>
            </div>
          </div>
          
          <div className="pt-2 border-t border-gray-700">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <div className="text-gray-400">Weapon:</div>
                <div className="text-purple-400">{killEvent.weapon}</div>
              </div>
              <div>
                <div className="text-gray-400">Distance:</div>
                <div>{Math.round(killEvent.distance)} units</div>
              </div>
            </div>
          </div>
          
          <div className="flex gap-4 pt-1">
            {killEvent.headshot && (
              <Badge variant="destructive" className="text-xs">
                Headshot
              </Badge>
            )}
            {killEvent.penetrated > 0 && (
              <Badge variant="secondary" className="text-xs">
                Wallbang ({killEvent.penetrated})
              </Badge>
            )}
            {killEvent.no_scope && (
              <Badge variant="secondary" className="text-xs">
                No-scope
              </Badge>
            )}
            {killEvent.attacker_blind && (
              <Badge variant="secondary" className="text-xs">
                Blind Kill
              </Badge>
            )}
          </div>
          
          {killEvent.assister && (
            <div className="pt-2 border-t border-gray-700">
              <div className="text-gray-400">Assist by:</div>
              <div className="text-yellow-400">{killEvent.assister.name}</div>
            </div>
          )}
        </div>
      )

    case 'player_hurt':
      const hurtEvent = event as IPlayerHurtEvent
      return (
        <div className="space-y-1 text-xs text-gray-300">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <div className="text-gray-400">Damage:</div>
              <div className="text-red-400">{hurtEvent.damage} HP</div>
              {hurtEvent.damage_armor > 0 && (
                <div className="text-blue-400">{hurtEvent.damage_armor} Armor</div>
              )}
            </div>
            <div>
              <div className="text-gray-400">Remaining:</div>
              <div className="text-green-400">{hurtEvent.health} HP</div>
              <div className="text-blue-400">{hurtEvent.armor} Armor</div>
            </div>
          </div>
          <div>
            <div className="text-gray-400">Weapon:</div>
            <div className="text-purple-400">{hurtEvent.weapon}</div>
          </div>
          <div>
            <div className="text-gray-400">Hit Group:</div>
            <div>{getHitGroupName(hurtEvent.hitgroup)}</div>
          </div>
        </div>
      )

    case 'bomb_planted':
    case 'bomb_defused':
      const bombEvent = event as IBombPlantEvent | IBombDefuseEvent
      return (
        <div className="space-y-1 text-xs text-gray-300">
          <div>
            <div className="text-gray-400">Player:</div>
            <div className="text-green-400">{bombEvent.player.name}</div>
          </div>
          <div>
            <div className="text-gray-400">Site:</div>
            <div className="text-yellow-400">Site {bombEvent.site}</div>
          </div>
          {'with_kit' in bombEvent && (
            <div>
              <div className="text-gray-400">Defuse Kit:</div>
              <div className={bombEvent.with_kit ? 'text-green-400' : 'text-red-400'}>
                {bombEvent.with_kit ? 'Yes' : 'No'}
              </div>
            </div>
          )}
          <div>
            <div className="text-gray-400">Position:</div>
            <div className="text-gray-500 font-mono">
              {Math.round(bombEvent.position.x)}, {Math.round(bombEvent.position.y)}, {Math.round(bombEvent.position.z)}
            </div>
          </div>
        </div>
      )

    case 'item_purchase':
      const purchaseEvent = event as IItemPurchaseEvent
      return (
        <div className="space-y-1 text-xs text-gray-300">
          <div>
            <div className="text-gray-400">Player:</div>
            <div className="text-blue-400">{purchaseEvent.player.name}</div>
          </div>
          <div>
            <div className="text-gray-400">Item:</div>
            <div className="text-purple-400">{purchaseEvent.item}</div>
          </div>
          <div>
            <div className="text-gray-400">Cost:</div>
            <div className="text-green-400">${purchaseEvent.cost}</div>
          </div>
        </div>
      )

    default:
      return null
  }
}

function getHitGroupName(hitgroup: number): string {
  const hitGroups: Record<number, string> = {
    0: 'Generic',
    1: 'Head',
    2: 'Chest', 
    3: 'Stomach',
    4: 'Left Arm',
    5: 'Right Arm',
    6: 'Left Leg',
    7: 'Right Leg'
  }
  return hitGroups[hitgroup] || `Unknown (${hitgroup})`
}

export default LogLine