"use client"

import React, { useState, useCallback, useEffect, useRef, useMemo } from 'react'
import { cn } from '@/lib/utils'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
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
import { LogLine } from './LogLine'
import type { ISpecificGameEvent } from '@/types/events'
import { filterEvents } from '@/types/events'

interface VirtualizedLogViewerProps {
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
  itemHeight?: number // Height of each log line in pixels
  visibleItems?: number // Number of visible items to render
}

export function VirtualizedLogViewer({
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
  itemHeight = 32,
  visibleItems = 50
}: VirtualizedLogViewerProps) {
  const [isPaused, setIsPaused] = useState(false)
  const [searchTerm, setSearchTerm] = useState('')
  const [scrollTop, setScrollTop] = useState(0)
  const [containerHeight, setContainerHeight] = useState(600)
  
  const scrollAreaRef = useRef<HTMLDivElement>(null)
  const endOfLogRef = useRef<HTMLDivElement>(null)

  // Filter events based on search term
  const filteredEvents = useMemo(() => {
    if (!searchTerm) return events
    
    return events.filter(event => {
      const searchText = `${event.type} ${JSON.stringify(event)}`.toLowerCase()
      return searchText.includes(searchTerm.toLowerCase())
    })
  }, [events, searchTerm])

  // Virtual scrolling calculations
  const totalItems = filteredEvents.length
  const totalHeight = totalItems * itemHeight
  const startIndex = Math.floor(scrollTop / itemHeight)
  const endIndex = Math.min(startIndex + visibleItems, totalItems)
  const visibleEvents = filteredEvents.slice(startIndex, endIndex)
  const offsetY = startIndex * itemHeight

  // Auto-scroll to bottom when new events arrive
  useEffect(() => {
    if (autoScroll && !isPaused && scrollAreaRef.current && filteredEvents.length > 0) {
      const container = scrollAreaRef.current
      const shouldScroll = container.scrollTop + container.clientHeight >= container.scrollHeight - 100
      
      if (shouldScroll) {
        container.scrollTop = container.scrollHeight
      }
    }
  }, [filteredEvents.length, autoScroll, isPaused])

  // Update container height on resize
  useEffect(() => {
    if (!scrollAreaRef.current) return

    const resizeObserver = new ResizeObserver((entries) => {
      for (const entry of entries) {
        setContainerHeight(entry.contentRect.height)
      }
    })

    resizeObserver.observe(scrollAreaRef.current)
    return () => resizeObserver.disconnect()
  }, [])

  const handleScroll = useCallback((e: React.UIEvent<HTMLDivElement>) => {
    setScrollTop(e.currentTarget.scrollTop)
  }, [])

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
    const logText = filteredEvents.map((event, index) => 
      `[${event.timestamp}] [TICK:${event.tick}] [R${event.round}] ${event.type}: ${JSON.stringify(event)}`
    ).join('\n')
    
    const blob = new Blob([logText], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `cs2-logs-${new Date().toISOString().slice(0, 19).replace(/:/g, '-')}.txt`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  }, [filteredEvents])

  return (
    <Card className={cn(
      'bg-gray-900 border-gray-700 text-green-400 font-mono text-sm flex flex-col',
      fullscreen && 'fixed inset-0 z-50 rounded-none',
      className
    )}>
      {/* Header with controls */}
      {showControls && (
        <div className="flex-shrink-0 flex items-center justify-between p-3 border-b border-gray-700 bg-gray-800">
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
              {filteredEvents.length} / {events.length} lines
            </Badge>
          </div>

          <div className="flex items-center gap-2">
            {/* Search */}
            <div className="relative">
              <Search className="absolute left-2 top-1/2 transform -translate-y-1/2 h-3 w-3 text-gray-400" />
              <Input
                type="text"
                placeholder="Search logs..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-7 pr-2 py-1 text-xs bg-gray-700 border border-gray-600 rounded focus:ring-1 focus:ring-blue-500 focus:border-blue-500 text-white placeholder-gray-400 w-40"
              />
            </div>

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
              disabled={filteredEvents.length === 0}
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

      {/* Virtualized log content */}
      <div className="flex-1 relative overflow-hidden">
        <div
          ref={scrollAreaRef}
          className="absolute inset-0 overflow-auto scrollbar-thin scrollbar-track-gray-800 scrollbar-thumb-gray-600"
          onScroll={handleScroll}
          style={{ height: fullscreen ? 'calc(100vh - 140px)' : '500px' }}
        >
          <div style={{ height: totalHeight, position: 'relative' }}>
            {filteredEvents.length === 0 ? (
              <div className="absolute inset-0 flex items-center justify-center text-gray-500">
                {events.length === 0 ? 'No log events yet...' : 'No events match your search'}
              </div>
            ) : (
              <>
                <div style={{ transform: `translateY(${offsetY}px)` }}>
                  {visibleEvents.map((event, index) => (
                    <LogLine
                      key={`${startIndex + index}-${event.tick}-${event.type}`}
                      event={event}
                      lineNumber={startIndex + index + 1}
                      highlighted={searchTerm && JSON.stringify(event).toLowerCase().includes(searchTerm.toLowerCase())}
                      searchTerm={searchTerm}
                      className="hover:bg-gray-800/50"
                    />
                  ))}
                </div>
                <div ref={endOfLogRef} />
              </>
            )}
          </div>
        </div>
      </div>

      {/* Footer with stats */}
      <div className="flex-shrink-0 px-3 py-2 bg-gray-800 border-t border-gray-700 text-xs text-gray-400">
        <div className="flex items-center justify-between">
          <span>
            Showing {startIndex + 1}-{Math.min(endIndex, totalItems)} of {totalItems} events
            {searchTerm && ` (filtered from ${events.length})`}
          </span>
          {isStreaming && (
            <span className="flex items-center gap-1">
              <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse" />
              Live
            </span>
          )}
        </div>
      </div>
    </Card>
  )
}

export default VirtualizedLogViewer