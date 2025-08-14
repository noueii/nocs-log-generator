/**
 * CS2 Log Parser Utility
 * Handles parsing, formatting, and color coding of CS2 log entries
 */

import type { 
  ISpecificGameEvent, 
  TEventType,
  IKillEvent,
  IPlayerHurtEvent,
  IBombPlantEvent,
  IBombDefuseEvent,
  IItemPurchaseEvent,
  IChatEvent,
  IRoundStartEvent,
  IRoundEndEvent
} from '@/types/events'
import { EVENT_COLORS, getEventIcon } from '@/types/events'

/**
 * Parsed log line interface
 */
export interface IParsedLogLine {
  id: string
  timestamp: string
  tick: number
  round: number
  event: ISpecificGameEvent
  rawText: string
  formattedText: string
  htmlText: string
  level: LogLevel
  tags: string[]
  searchableText: string
}

export type LogLevel = 'info' | 'warning' | 'error' | 'success' | 'debug'

/**
 * Log parsing options
 */
export interface ILogParseOptions {
  enableColorCoding?: boolean
  enableSyntaxHighlighting?: boolean
  extractPlayerNames?: boolean
  extractWeaponNames?: boolean
  extractPositions?: boolean
}

/**
 * Color coding configuration
 */
export interface IColorConfig {
  players: string
  weapons: string
  teams: string
  numbers: string
  positions: string
  steamIds: string
  timestamps: string
  damages: string
  money: string
}

const DEFAULT_COLORS: IColorConfig = {
  players: 'text-cyan-400',
  weapons: 'text-purple-400',
  teams: 'text-green-400',
  numbers: 'text-yellow-400',
  positions: 'text-blue-400',
  steamIds: 'text-gray-400',
  timestamps: 'text-gray-500',
  damages: 'text-red-400',
  money: 'text-green-400'
}

/**
 * Regular expressions for parsing different elements
 */
const PATTERNS = {
  // Player names (assuming format: "PlayerName<userid><STEAM_ID><team>")
  player: /(\w+)<(\d+)><(STEAM_[0-5]:[01]:\d+)><(\w+)>/g,
  
  // Weapons
  weapon: /(ak47|m4a4|m4a1|awp|deagle|glock|usp|p250|tec9|five-seven|cz75|dual_berettas|mac10|mp9|ump45|p90|bizon|mag7|sawed_off|nova|xm1014|m249|negev|knife|hegrenade|flashbang|smokegrenade|incgrenade|molotov|decoy)/gi,
  
  // Steam IDs
  steamId: /(STEAM_[0-5]:[01]:\d+)/g,
  
  // Damage numbers
  damage: /(\d+) damage/gi,
  
  // HP/Health numbers
  health: /(\d+) HP/gi,
  
  // Money amounts
  money: /\$(\d+)/g,
  
  // Positions (coordinates)
  position: /\((-?\d+\.?\d*)\s+(-?\d+\.?\d*)\s+(-?\d+\.?\d*)\)/g,
  
  // Teams
  team: /(CT|TERRORIST|Counter-Terrorist|Terrorist)/gi,
  
  // Sites
  site: /(site [AB])/gi,
  
  // Round numbers
  round: /(Round \d+)/gi,
  
  // Timestamps
  timestamp: /(\d{2}:\d{2}:\d{2})/g,
  
  // Tick numbers
  tick: /(tick:?\s*\d+)/gi,
  
  // Special events
  headshot: /\(headshot\)/gi,
  penetration: /\(penetrated\)/gi,
  noScope: /\(noscope\)/gi,
  wallbang: /\(wallbang\)/gi
}

/**
 * Parse raw CS2 log text into structured format
 */
export function parseLogText(
  rawText: string, 
  options: ILogParseOptions = {}
): IParsedLogLine[] {
  const lines = rawText.split('\n').filter(line => line.trim())
  const parsedLines: IParsedLogLine[] = []

  lines.forEach((line, index) => {
    const parsed = parseLogLine(line, index, options)
    if (parsed) {
      parsedLines.push(parsed)
    }
  })

  return parsedLines
}

/**
 * Parse a single log line
 */
export function parseLogLine(
  rawLine: string,
  lineIndex: number,
  options: ILogParseOptions = {}
): IParsedLogLine | null {
  const trimmedLine = rawLine.trim()
  if (!trimmedLine) return null

  // Extract basic components (this is a simplified version)
  const parts = trimmedLine.split(' ')
  const timestamp = extractTimestamp(trimmedLine)
  const tick = extractTick(trimmedLine)
  const round = extractRound(trimmedLine)

  // Create a mock event for now (in real implementation, this would parse actual log format)
  const event = createMockEventFromLine(trimmedLine, tick, round, timestamp)
  if (!event) return null

  const formattedText = formatEventText(event)
  const htmlText = options.enableSyntaxHighlighting 
    ? applySyntaxHighlighting(formattedText, options)
    : formattedText

  const level = determineLogLevel(event)
  const tags = extractTags(event, trimmedLine)
  const searchableText = createSearchableText(event, formattedText)

  return {
    id: `line-${lineIndex}-${tick}-${event.type}`,
    timestamp,
    tick,
    round,
    event,
    rawText: rawLine,
    formattedText,
    htmlText,
    level,
    tags,
    searchableText
  }
}

/**
 * Apply syntax highlighting to formatted text
 */
export function applySyntaxHighlighting(
  text: string, 
  options: ILogParseOptions = {},
  colors: IColorConfig = DEFAULT_COLORS
): string {
  let highlighted = text

  if (options.extractPlayerNames !== false) {
    // Highlight player names
    highlighted = highlighted.replace(
      PATTERNS.player,
      (match, name, userid, steamid, team) => 
        `<span class="${colors.players} font-medium">${name}</span><span class="${colors.steamIds}">&lt;${userid}&gt;&lt;${steamid}&gt;</span><span class="${colors.teams}">&lt;${team}&gt;</span>`
    )
  }

  if (options.extractWeaponNames !== false) {
    // Highlight weapons
    highlighted = highlighted.replace(
      PATTERNS.weapon,
      `<span class="${colors.weapons} font-medium">$1</span>`
    )
  }

  // Highlight Steam IDs
  highlighted = highlighted.replace(
    PATTERNS.steamId,
    `<span class="${colors.steamIds} underline cursor-pointer hover:text-blue-200" title="Click to copy Steam ID">$1</span>`
  )

  // Highlight damage numbers
  highlighted = highlighted.replace(
    PATTERNS.damage,
    `<span class="${colors.damages} font-bold">$1</span>`
  )

  // Highlight health numbers
  highlighted = highlighted.replace(
    PATTERNS.health,
    `<span class="${colors.damages} font-medium">$1</span>`
  )

  // Highlight money
  highlighted = highlighted.replace(
    PATTERNS.money,
    `<span class="${colors.money} font-medium">$$$1</span>`
  )

  // Highlight positions
  if (options.extractPositions !== false) {
    highlighted = highlighted.replace(
      PATTERNS.position,
      `<span class="${colors.positions} font-mono">($1 $2 $3)</span>`
    )
  }

  // Highlight teams
  highlighted = highlighted.replace(
    PATTERNS.team,
    `<span class="${colors.teams} font-medium">$1</span>`
  )

  // Highlight sites
  highlighted = highlighted.replace(
    PATTERNS.site,
    `<span class="${colors.numbers} font-medium">$1</span>`
  )

  // Highlight round numbers
  highlighted = highlighted.replace(
    PATTERNS.round,
    `<span class="${colors.numbers} font-medium">$1</span>`
  )

  // Highlight special events
  highlighted = highlighted.replace(
    PATTERNS.headshot,
    '<span class="text-red-500 font-bold">$1</span>'
  )

  highlighted = highlighted.replace(
    PATTERNS.penetration,
    '<span class="text-orange-500 font-medium">$1</span>'
  )

  highlighted = highlighted.replace(
    PATTERNS.noScope,
    '<span class="text-purple-500 font-medium">$1</span>'
  )

  return highlighted
}

/**
 * Extract timestamp from log line
 */
export function extractTimestamp(line: string): string {
  const match = line.match(/(\d{2}:\d{2}:\d{2})/)
  if (match) return match[1]

  // Fallback: try to extract from beginning of line
  const timeMatch = line.match(/^(\d{2}:\d{2}:\d{2})/)
  if (timeMatch) return timeMatch[1]

  // Default timestamp
  return new Date().toLocaleTimeString('en-US', { 
    hour12: false, 
    hour: '2-digit', 
    minute: '2-digit', 
    second: '2-digit' 
  })
}

/**
 * Extract tick number from log line
 */
export function extractTick(line: string): number {
  const match = line.match(/tick:?\s*(\d+)/i)
  return match ? parseInt(match[1]) : Math.floor(Math.random() * 100000)
}

/**
 * Extract round number from log line
 */
export function extractRound(line: string): number {
  const match = line.match(/round\s*(\d+)/i)
  return match ? parseInt(match[1]) : 1
}

/**
 * Create mock event from log line (simplified for demo)
 */
function createMockEventFromLine(
  line: string, 
  tick: number, 
  round: number, 
  timestamp: string
): ISpecificGameEvent | null {
  const lowerLine = line.toLowerCase()

  // Determine event type based on keywords
  let eventType: TEventType = 'server_command'

  if (lowerLine.includes('killed')) {
    eventType = 'player_death'
  } else if (lowerLine.includes('hurt') || lowerLine.includes('damaged')) {
    eventType = 'player_hurt'
  } else if (lowerLine.includes('planted')) {
    eventType = 'bomb_planted'
  } else if (lowerLine.includes('defused')) {
    eventType = 'bomb_defused'
  } else if (lowerLine.includes('exploded')) {
    eventType = 'bomb_exploded'
  } else if (lowerLine.includes('purchased')) {
    eventType = 'item_purchase'
  } else if (lowerLine.includes('round') && lowerLine.includes('start')) {
    eventType = 'round_start'
  } else if (lowerLine.includes('round') && lowerLine.includes('end')) {
    eventType = 'round_end'
  } else if (lowerLine.includes('say') || lowerLine.includes('chat')) {
    eventType = 'chat_message'
  }

  // Create basic event structure
  const baseEvent = {
    timestamp,
    tick,
    round,
    type: eventType
  }

  return baseEvent as ISpecificGameEvent
}

/**
 * Format event for display
 */
function formatEventText(event: ISpecificGameEvent): string {
  // This would be more sophisticated in real implementation
  switch (event.type) {
    case 'player_death':
      return `Player killed another player`
    case 'player_hurt':
      return `Player hurt another player`
    case 'bomb_planted':
      return `Bomb planted`
    case 'bomb_defused':
      return `Bomb defused`
    case 'round_start':
      return `Round ${event.round} started`
    case 'round_end':
      return `Round ${event.round} ended`
    default:
      return `${event.type} event`
  }
}

/**
 * Determine log level for an event
 */
function determineLogLevel(event: ISpecificGameEvent): LogLevel {
  switch (event.type) {
    case 'player_death':
    case 'bomb_exploded':
      return 'error'
    
    case 'bomb_planted':
    case 'player_hurt':
      return 'warning'
    
    case 'round_start':
    case 'round_end':
    case 'bomb_defused':
      return 'success'
    
    case 'weapon_fire':
    case 'server_command':
      return 'debug'
    
    default:
      return 'info'
  }
}

/**
 * Extract searchable tags from event
 */
function extractTags(event: ISpecificGameEvent, rawLine: string): string[] {
  const tags: string[] = [event.type]

  // Add level as tag
  const level = determineLogLevel(event)
  tags.push(level)

  // Extract player names as tags
  const playerMatches = rawLine.match(PATTERNS.player)
  if (playerMatches) {
    playerMatches.forEach(match => {
      const playerName = match.split('<')[0]
      if (playerName) tags.push(`player:${playerName}`)
    })
  }

  // Extract weapon names as tags
  const weaponMatches = rawLine.match(PATTERNS.weapon)
  if (weaponMatches) {
    weaponMatches.forEach(weapon => {
      tags.push(`weapon:${weapon}`)
    })
  }

  // Add round and tick as searchable tags
  tags.push(`round:${event.round}`)
  tags.push(`tick:${event.tick}`)

  return tags
}

/**
 * Create searchable text for filtering
 */
function createSearchableText(event: ISpecificGameEvent, formattedText: string): string {
  const searchable = [
    formattedText,
    event.type,
    `round ${event.round}`,
    `tick ${event.tick}`
  ]

  return searchable.join(' ').toLowerCase()
}

/**
 * Filter parsed lines by search term
 */
export function filterLogLines(
  lines: IParsedLogLine[],
  searchTerm: string,
  includeEvents?: TEventType[],
  excludeEvents?: TEventType[]
): IParsedLogLine[] {
  if (!searchTerm && !includeEvents && !excludeEvents) {
    return lines
  }

  return lines.filter(line => {
    // Event type filtering
    if (includeEvents && !includeEvents.includes(line.event.type)) {
      return false
    }
    
    if (excludeEvents && excludeEvents.includes(line.event.type)) {
      return false
    }

    // Search term filtering
    if (searchTerm) {
      const term = searchTerm.toLowerCase()
      return (
        line.searchableText.includes(term) ||
        line.tags.some(tag => tag.toLowerCase().includes(term))
      )
    }

    return true
  })
}

/**
 * Export logs to different formats
 */
export function exportLogs(lines: IParsedLogLine[], format: 'text' | 'csv' | 'json' = 'text'): string {
  switch (format) {
    case 'text':
      return lines.map(line => 
        `[${line.timestamp}] [TICK:${line.tick}] [R${line.round}] ${line.formattedText}`
      ).join('\n')

    case 'csv':
      const csvHeader = 'Timestamp,Tick,Round,Event Type,Text,Level\n'
      const csvRows = lines.map(line =>
        `"${line.timestamp}",${line.tick},${line.round},"${line.event.type}","${line.formattedText}","${line.level}"`
      )
      return csvHeader + csvRows.join('\n')

    case 'json':
      return JSON.stringify(lines, null, 2)

    default:
      return exportLogs(lines, 'text')
  }
}

/**
 * Parse log statistics
 */
export interface ILogStats {
  totalLines: number
  eventCounts: Record<TEventType, number>
  playerStats: Record<string, { kills: number; deaths: number; damage: number }>
  roundCount: number
  timeSpan: { start: string; end: string }
}

export function calculateLogStats(lines: IParsedLogLine[]): ILogStats {
  const stats: ILogStats = {
    totalLines: lines.length,
    eventCounts: {} as Record<TEventType, number>,
    playerStats: {},
    roundCount: 0,
    timeSpan: { start: '', end: '' }
  }

  if (lines.length === 0) return stats

  // Calculate event counts
  lines.forEach(line => {
    const eventType = line.event.type
    stats.eventCounts[eventType] = (stats.eventCounts[eventType] || 0) + 1
  })

  // Calculate round count
  stats.roundCount = Math.max(...lines.map(line => line.round))

  // Calculate time span
  stats.timeSpan.start = lines[0].timestamp
  stats.timeSpan.end = lines[lines.length - 1].timestamp

  return stats
}