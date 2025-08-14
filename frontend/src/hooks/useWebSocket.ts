"use client"

import { useEffect, useRef, useState, useCallback } from 'react'
import type { ISpecificGameEvent, IEventStreamMessage } from '@/types/events'

export type WebSocketStatus = 'connecting' | 'connected' | 'disconnected' | 'error' | 'reconnecting'

export interface UseWebSocketOptions {
  url: string
  protocols?: string | string[]
  reconnectAttempts?: number
  reconnectInterval?: number
  heartbeatInterval?: number
  onOpen?: () => void
  onClose?: (event: CloseEvent) => void
  onError?: (error: Event) => void
  onMessage?: (event: MessageEvent) => void
  autoConnect?: boolean
}

export interface UseWebSocketReturn {
  status: WebSocketStatus
  lastMessage: MessageEvent | null
  sendMessage: (message: string | object) => void
  sendJsonMessage: (message: object) => void
  connect: () => void
  disconnect: () => void
  reconnect: () => void
  isConnected: boolean
  connectionAttempts: number
}

export function useWebSocket(options: UseWebSocketOptions): UseWebSocketReturn {
  const {
    url,
    protocols,
    reconnectAttempts = 5,
    reconnectInterval = 3000,
    heartbeatInterval = 30000,
    onOpen,
    onClose,
    onError,
    onMessage,
    autoConnect = true
  } = options

  const [status, setStatus] = useState<WebSocketStatus>('disconnected')
  const [lastMessage, setLastMessage] = useState<MessageEvent | null>(null)
  const [connectionAttempts, setConnectionAttempts] = useState(0)

  const websocketRef = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null)
  const heartbeatTimeoutRef = useRef<NodeJS.Timeout | null>(null)
  const reconnectAttemptsRef = useRef(0)

  const clearTimeouts = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
      reconnectTimeoutRef.current = null
    }
    if (heartbeatTimeoutRef.current) {
      clearTimeout(heartbeatTimeoutRef.current)
      heartbeatTimeoutRef.current = null
    }
  }, [])

  const startHeartbeat = useCallback(() => {
    if (heartbeatInterval > 0) {
      heartbeatTimeoutRef.current = setInterval(() => {
        if (websocketRef.current?.readyState === WebSocket.OPEN) {
          websocketRef.current.send(JSON.stringify({ type: 'ping' }))
        }
      }, heartbeatInterval)
    }
  }, [heartbeatInterval])

  const connect = useCallback(() => {
    if (websocketRef.current?.readyState === WebSocket.OPEN) {
      return // Already connected
    }

    clearTimeouts()
    setStatus('connecting')

    try {
      websocketRef.current = new WebSocket(url, protocols)

      websocketRef.current.onopen = (event) => {
        setStatus('connected')
        reconnectAttemptsRef.current = 0
        setConnectionAttempts(0)
        startHeartbeat()
        onOpen?.()
      }

      websocketRef.current.onclose = (event) => {
        setStatus('disconnected')
        clearTimeouts()
        
        onClose?.(event)

        // Attempt reconnection if not manually closed and attempts remain
        if (!event.wasClean && reconnectAttemptsRef.current < reconnectAttempts) {
          reconnectAttemptsRef.current++
          setConnectionAttempts(reconnectAttemptsRef.current)
          setStatus('reconnecting')
          
          reconnectTimeoutRef.current = setTimeout(() => {
            connect()
          }, reconnectInterval)
        } else if (reconnectAttemptsRef.current >= reconnectAttempts) {
          setStatus('error')
        }
      }

      websocketRef.current.onerror = (event) => {
        setStatus('error')
        onError?.(event)
      }

      websocketRef.current.onmessage = (event) => {
        setLastMessage(event)
        onMessage?.(event)
      }

    } catch (error) {
      setStatus('error')
      console.error('WebSocket connection error:', error)
    }
  }, [url, protocols, reconnectAttempts, reconnectInterval, onOpen, onClose, onError, onMessage, clearTimeouts, startHeartbeat])

  const disconnect = useCallback(() => {
    clearTimeouts()
    reconnectAttemptsRef.current = reconnectAttempts // Prevent reconnection
    
    if (websocketRef.current) {
      websocketRef.current.close(1000, 'Manual disconnect')
      websocketRef.current = null
    }
    
    setStatus('disconnected')
  }, [reconnectAttempts, clearTimeouts])

  const reconnect = useCallback(() => {
    reconnectAttemptsRef.current = 0
    setConnectionAttempts(0)
    disconnect()
    setTimeout(() => connect(), 100)
  }, [connect, disconnect])

  const sendMessage = useCallback((message: string | object) => {
    if (websocketRef.current?.readyState === WebSocket.OPEN) {
      const messageStr = typeof message === 'string' ? message : JSON.stringify(message)
      websocketRef.current.send(messageStr)
    } else {
      console.warn('WebSocket is not connected')
    }
  }, [])

  const sendJsonMessage = useCallback((message: object) => {
    sendMessage(JSON.stringify(message))
  }, [sendMessage])

  // Auto-connect on mount if enabled
  useEffect(() => {
    if (autoConnect) {
      connect()
    }

    // Cleanup on unmount
    return () => {
      clearTimeouts()
      if (websocketRef.current) {
        websocketRef.current.close(1000, 'Component unmounting')
      }
    }
  }, [autoConnect, connect, clearTimeouts])

  return {
    status,
    lastMessage,
    sendMessage,
    sendJsonMessage,
    connect,
    disconnect,
    reconnect,
    isConnected: status === 'connected',
    connectionAttempts
  }
}

// Specialized hook for match event streaming
export interface UseMatchStreamOptions {
  matchId?: string
  autoConnect?: boolean
  onEvent?: (event: ISpecificGameEvent) => void
  onStatusChange?: (status: string) => void
  onError?: (error: string) => void
}

export interface UseMatchStreamReturn {
  events: ISpecificGameEvent[]
  status: WebSocketStatus
  isConnected: boolean
  connect: () => void
  disconnect: () => void
  pause: () => void
  resume: () => void
  clear: () => void
  connectionAttempts: number
  isPaused: boolean
}

export function useMatchStream(options: UseMatchStreamOptions): UseMatchStreamReturn {
  const {
    matchId,
    autoConnect = false,
    onEvent,
    onStatusChange,
    onError
  } = options

  const [events, setEvents] = useState<ISpecificGameEvent[]>([])
  const [isPaused, setIsPaused] = useState(false)

  const websocketUrl = matchId 
    ? `${getWebSocketUrl()}/stream/${matchId}`
    : `${getWebSocketUrl()}/stream`

  const handleMessage = useCallback((event: MessageEvent) => {
    try {
      const message: IEventStreamMessage = JSON.parse(event.data)
      
      switch (message.type) {
        case 'event':
          if (!isPaused) {
            const gameEvent = message.data as ISpecificGameEvent
            setEvents(prev => [...prev, gameEvent])
            onEvent?.(gameEvent)
          }
          break
          
        case 'status':
          onStatusChange?.(message.data as string)
          break
          
        case 'error':
          onError?.(message.data as string)
          break
      }
    } catch (error) {
      console.error('Failed to parse WebSocket message:', error)
      onError?.('Failed to parse message')
    }
  }, [isPaused, onEvent, onStatusChange, onError])

  const {
    status,
    connect: wsConnect,
    disconnect: wsDisconnect,
    sendJsonMessage,
    isConnected,
    connectionAttempts
  } = useWebSocket({
    url: websocketUrl,
    autoConnect: autoConnect && !!matchId,
    onMessage: handleMessage,
    onError: (error) => {
      console.error('WebSocket error:', error)
      onError?.('WebSocket connection error')
    }
  })

  const connect = useCallback(() => {
    if (matchId) {
      wsConnect()
    } else {
      console.warn('Cannot connect without matchId')
    }
  }, [matchId, wsConnect])

  const disconnect = useCallback(() => {
    wsDisconnect()
  }, [wsDisconnect])

  const pause = useCallback(() => {
    setIsPaused(true)
    sendJsonMessage({ type: 'pause' })
  }, [sendJsonMessage])

  const resume = useCallback(() => {
    setIsPaused(false)
    sendJsonMessage({ type: 'resume' })
  }, [sendJsonMessage])

  const clear = useCallback(() => {
    setEvents([])
  }, [])

  // Subscribe to match when connected
  useEffect(() => {
    if (isConnected && matchId) {
      sendJsonMessage({ 
        type: 'subscribe',
        matchId: matchId
      })
    }
  }, [isConnected, matchId, sendJsonMessage])

  return {
    events,
    status,
    isConnected,
    connect,
    disconnect,
    pause,
    resume,
    clear,
    connectionAttempts,
    isPaused
  }
}

// Helper to get WebSocket URL
function getWebSocketUrl(): string {
  if (typeof window === 'undefined') {
    return 'ws://localhost:8080/ws'
  }

  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  
  // In development, use localhost:8080
  if (host.includes('localhost') || host.includes('127.0.0.1')) {
    return 'ws://localhost:8080/ws'
  }
  
  return `${protocol}//${host}/ws`
}