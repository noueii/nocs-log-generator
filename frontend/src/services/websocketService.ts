/**
 * WebSocket service for real-time match generation streaming
 * Handles WebSocket connections, events, and error recovery
 */

import { IGenerationProgress } from './matchService';

/**
 * WebSocket connection states
 */
export type TWebSocketState = 'connecting' | 'connected' | 'disconnected' | 'error' | 'reconnecting';

/**
 * WebSocket event types
 */
export interface IWebSocketEvents {
  onStateChange: (state: TWebSocketState) => void;
  onProgress: (progress: IGenerationProgress) => void;
  onError: (error: string) => void;
  onMessage: (message: any) => void;
}

/**
 * WebSocket service configuration
 */
interface IWebSocketConfig {
  baseUrl: string;
  reconnectAttempts: number;
  reconnectInterval: number;
  heartbeatInterval: number;
  timeout: number;
}

/**
 * Default WebSocket configuration
 */
const defaultConfig: IWebSocketConfig = {
  baseUrl: import.meta.env.VITE_WS_BASE_URL || 'ws://localhost:8080',
  reconnectAttempts: 5,
  reconnectInterval: 3000,
  heartbeatInterval: 30000,
  timeout: 10000,
};

/**
 * WebSocket Service Class
 */
export class WebSocketService {
  private ws: WebSocket | null = null;
  private config: IWebSocketConfig;
  private state: TWebSocketState = 'disconnected';
  private listeners: Partial<IWebSocketEvents> = {};
  private reconnectCount = 0;
  private reconnectTimer?: NodeJS.Timeout;
  private heartbeatTimer?: NodeJS.Timeout;
  private isManualClose = false;
  private subscriptions = new Set<string>();

  constructor(config: Partial<IWebSocketConfig> = {}) {
    this.config = { ...defaultConfig, ...config };
    
    // Check if WebSocket is enabled
    if (import.meta.env.VITE_ENABLE_WEBSOCKET !== 'true') {
      console.log('📡 WebSocket disabled by configuration');
      return;
    }
  }

  /**
   * Connect to WebSocket server
   */
  async connect(): Promise<void> {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      return;
    }

    if (import.meta.env.VITE_ENABLE_WEBSOCKET !== 'true') {
      console.log('📡 WebSocket disabled, skipping connection');
      return;
    }

    try {
      this.setState('connecting');
      this.isManualClose = false;

      this.ws = new WebSocket(`${this.config.baseUrl}/ws`);
      
      // Set up event handlers
      this.ws.onopen = this.onOpen.bind(this);
      this.ws.onmessage = this.onMessage.bind(this);
      this.ws.onerror = this.onError.bind(this);
      this.ws.onclose = this.onClose.bind(this);

      // Connection timeout
      setTimeout(() => {
        if (this.ws && this.ws.readyState === WebSocket.CONNECTING) {
          this.ws.close();
          this.onError({ error: 'Connection timeout' } as ErrorEvent);
        }
      }, this.config.timeout);

    } catch (error) {
      console.error('📡 WebSocket connection failed:', error);
      this.setState('error');
      this.listeners.onError?.(`Connection failed: ${error}`);
      this.scheduleReconnect();
    }
  }

  /**
   * Disconnect from WebSocket server
   */
  disconnect(): void {
    this.isManualClose = true;
    this.clearReconnectTimer();
    this.clearHeartbeatTimer();
    
    if (this.ws) {
      this.ws.close(1000, 'Manual disconnect');
      this.ws = null;
    }
    
    this.setState('disconnected');
    this.subscriptions.clear();
  }

  /**
   * Subscribe to match generation progress
   */
  subscribeToMatch(matchId: string): void {
    if (!this.isConnected()) {
      console.warn('📡 Cannot subscribe: WebSocket not connected');
      return;
    }

    const message = {
      type: 'subscribe',
      match_id: matchId,
      timestamp: new Date().toISOString(),
    };

    this.send(message);
    this.subscriptions.add(matchId);
  }

  /**
   * Unsubscribe from match generation progress
   */
  unsubscribeFromMatch(matchId: string): void {
    if (!this.isConnected()) {
      return;
    }

    const message = {
      type: 'unsubscribe',
      match_id: matchId,
      timestamp: new Date().toISOString(),
    };

    this.send(message);
    this.subscriptions.delete(matchId);
  }

  /**
   * Send message to server
   */
  private send(message: any): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
    } else {
      console.warn('📡 Cannot send message: WebSocket not open');
    }
  }

  /**
   * Set event listeners
   */
  setListeners(listeners: Partial<IWebSocketEvents>): void {
    this.listeners = { ...this.listeners, ...listeners };
  }

  /**
   * Check if WebSocket is connected
   */
  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }

  /**
   * Get current connection state
   */
  getState(): TWebSocketState {
    return this.state;
  }

  // Private event handlers

  private onOpen(): void {
    console.log('📡 WebSocket connected');
    this.setState('connected');
    this.reconnectCount = 0;
    this.startHeartbeat();
    
    // Re-subscribe to all active subscriptions
    this.subscriptions.forEach(matchId => {
      this.subscribeToMatch(matchId);
    });
  }

  private onMessage(event: MessageEvent): void {
    try {
      const data = JSON.parse(event.data);
      
      if (import.meta.env.DEV || import.meta.env.VITE_ENABLE_DEBUG_LOGS === 'true') {
        console.log('📡 WebSocket message:', data);
      }

      // Handle different message types
      switch (data.type) {
        case 'progress':
          this.listeners.onProgress?.(data);
          break;
          
        case 'error':
          this.listeners.onError?.(data.message || 'Unknown error');
          break;
          
        case 'heartbeat':
          // Respond to heartbeat
          this.send({ type: 'heartbeat_response' });
          break;
          
        default:
          this.listeners.onMessage?.(data);
      }
    } catch (error) {
      console.error('📡 Failed to parse WebSocket message:', error);
      this.listeners.onError?.('Failed to parse server message');
    }
  }

  private onError(event: ErrorEvent): void {
    console.error('📡 WebSocket error:', event.error || event.message);
    this.setState('error');
    this.listeners.onError?.(event.error?.toString() || 'WebSocket error');
    
    if (!this.isManualClose) {
      this.scheduleReconnect();
    }
  }

  private onClose(event: CloseEvent): void {
    console.log('📡 WebSocket closed:', event.code, event.reason);
    this.clearHeartbeatTimer();
    this.ws = null;
    
    if (!this.isManualClose && event.code !== 1000) {
      this.setState('disconnected');
      this.scheduleReconnect();
    } else {
      this.setState('disconnected');
    }
  }

  // Private utility methods

  private setState(newState: TWebSocketState): void {
    if (this.state !== newState) {
      this.state = newState;
      this.listeners.onStateChange?.(newState);
    }
  }

  private scheduleReconnect(): void {
    if (this.isManualClose || this.reconnectCount >= this.config.reconnectAttempts) {
      console.log('📡 Max reconnection attempts reached');
      this.setState('error');
      return;
    }

    this.setState('reconnecting');
    this.reconnectCount++;
    
    const delay = this.config.reconnectInterval * Math.pow(1.5, this.reconnectCount - 1);
    
    console.log(`📡 Scheduling reconnection attempt ${this.reconnectCount}/${this.config.reconnectAttempts} in ${delay}ms`);
    
    this.reconnectTimer = setTimeout(() => {
      this.connect();
    }, delay);
  }

  private clearReconnectTimer(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = undefined;
    }
  }

  private startHeartbeat(): void {
    this.clearHeartbeatTimer();
    
    this.heartbeatTimer = setInterval(() => {
      if (this.isConnected()) {
        this.send({ type: 'heartbeat', timestamp: new Date().toISOString() });
      }
    }, this.config.heartbeatInterval);
  }

  private clearHeartbeatTimer(): void {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer);
      this.heartbeatTimer = undefined;
    }
  }
}

/**
 * Create and export singleton WebSocket service instance
 */
export const websocketService = new WebSocketService();

/**
 * Hook for React components to use WebSocket
 */
export const useWebSocket = (
  listeners: Partial<IWebSocketEvents> = {},
  autoConnect = true
) => {
  const [state, setState] = React.useState<TWebSocketState>(websocketService.getState());
  
  React.useEffect(() => {
    // Set up listeners
    websocketService.setListeners({
      ...listeners,
      onStateChange: (newState) => {
        setState(newState);
        listeners.onStateChange?.(newState);
      },
    });

    // Auto connect if enabled
    if (autoConnect && import.meta.env.VITE_ENABLE_WEBSOCKET === 'true') {
      websocketService.connect();
    }

    return () => {
      // Clean up on unmount
      if (!autoConnect) {
        websocketService.disconnect();
      }
    };
  }, []);

  return {
    state,
    connect: () => websocketService.connect(),
    disconnect: () => websocketService.disconnect(),
    subscribeToMatch: (matchId: string) => websocketService.subscribeToMatch(matchId),
    unsubscribeFromMatch: (matchId: string) => websocketService.unsubscribeFromMatch(matchId),
    isConnected: websocketService.isConnected(),
  };
};

// Import React for the hook
import React from 'react';

export default websocketService;