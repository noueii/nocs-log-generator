/**
 * Event-related TypeScript interfaces mirroring the Go GameEvent model
 * Maps to backend/pkg/models/events.go
 */

import type { IPlayer } from './player';
import type { IVector3 } from './match';
import type { ITeamEconomy } from './team';

/**
 * Event type enum
 */
export type TEventType = 
  | 'player_death'
  | 'round_start'
  | 'round_end'
  | 'bomb_planted'
  | 'bomb_defused'
  | 'bomb_exploded'
  | 'player_hurt'
  | 'player_connect'
  | 'player_disconnect'
  | 'item_purchase'
  | 'grenade_thrown'
  | 'weapon_fire'
  | 'flashbang_detonate'
  | 'chat_message'
  | 'team_switch'
  | 'server_command';

/**
 * Base game event interface
 */
export interface IGameEvent {
  timestamp: string;
  type: TEventType;
  tick: number;
  round: number;
}

/**
 * Player kill event
 */
export interface IKillEvent extends IGameEvent {
  type: 'player_death';
  attacker: IPlayer;
  victim: IPlayer;
  assister?: IPlayer;
  weapon: string;
  headshot: boolean;
  penetrated: number;
  no_scope: boolean;
  attacker_blind: boolean;
  distance: number;
  attacker_pos: IVector3;
  victim_pos: IVector3;
}

/**
 * Round start event
 */
export interface IRoundStartEvent extends IGameEvent {
  type: 'round_start';
  ct_score: number;
  t_score: number;
  ct_players: number;
  t_players: number;
  team_economies: Record<string, ITeamEconomy>;
}

/**
 * Round end event
 */
export interface IRoundEndEvent extends IGameEvent {
  type: 'round_end';
  winner: 'CT' | 'TERRORIST';
  reason: 'elimination' | 'bomb_defused' | 'bomb_exploded' | 'time';
  ct_score: number;
  t_score: number;
  mvp?: IPlayer;
}

/**
 * Bomb plant event
 */
export interface IBombPlantEvent extends IGameEvent {
  type: 'bomb_planted';
  player: IPlayer;
  site: 'A' | 'B';
  position: IVector3;
}

/**
 * Bomb defuse event
 */
export interface IBombDefuseEvent extends IGameEvent {
  type: 'bomb_defused';
  player: IPlayer;
  site: 'A' | 'B';
  with_kit: boolean;
  position: IVector3;
}

/**
 * Bomb explosion event
 */
export interface IBombExplodeEvent extends IGameEvent {
  type: 'bomb_exploded';
  site: 'A' | 'B';
  position: IVector3;
}

/**
 * Player hurt event
 */
export interface IPlayerHurtEvent extends IGameEvent {
  type: 'player_hurt';
  attacker: IPlayer;
  victim: IPlayer;
  weapon: string;
  damage: number;
  damage_armor: number;
  health: number;
  armor: number;
  hitgroup: number; // 0=generic, 1=head, 2=chest, 3=stomach, 4=leftarm, 5=rightarm, 6=leftleg, 7=rightleg
}

/**
 * Player connect event
 */
export interface IPlayerConnectEvent extends IGameEvent {
  type: 'player_connect';
  player: IPlayer;
  address: string;
}

/**
 * Player disconnect event
 */
export interface IPlayerDisconnectEvent extends IGameEvent {
  type: 'player_disconnect';
  player: IPlayer;
  reason: string;
}

/**
 * Equipment purchase event
 */
export interface IItemPurchaseEvent extends IGameEvent {
  type: 'item_purchase';
  player: IPlayer;
  item: string;
  cost: number;
}

/**
 * Grenade throw event
 */
export interface IGrenadeThrowEvent extends IGameEvent {
  type: 'grenade_thrown';
  player: IPlayer;
  grenade_type: string;
  position: IVector3;
  velocity: IVector3;
}

/**
 * Weapon fire event
 */
export interface IWeaponFireEvent extends IGameEvent {
  type: 'weapon_fire';
  player: IPlayer;
  weapon: string;
  position: IVector3;
  angle: IVector3;
  silenced: boolean;
}

/**
 * Flashbang detonation event
 */
export interface IFlashbangEvent extends IGameEvent {
  type: 'flashbang_detonate';
  player: IPlayer;
  position: IVector3;
  flashed: IPlayer[];
  duration: number;
}

/**
 * Chat message event
 */
export interface IChatEvent extends IGameEvent {
  type: 'chat_message';
  player?: IPlayer;
  message: string;
  team: boolean;
  dead: boolean;
}

/**
 * Team switch event
 */
export interface ITeamSwitchEvent extends IGameEvent {
  type: 'team_switch';
  player: IPlayer;
  from_team: string;
  to_team: string;
}

/**
 * Server command event
 */
export interface IServerCommandEvent extends IGameEvent {
  type: 'server_command';
  command: string;
  args: string;
  result?: string;
}

/**
 * Union type for all event types
 */
export type ISpecificGameEvent = 
  | IKillEvent
  | IRoundStartEvent
  | IRoundEndEvent
  | IBombPlantEvent
  | IBombDefuseEvent
  | IBombExplodeEvent
  | IPlayerHurtEvent
  | IPlayerConnectEvent
  | IPlayerDisconnectEvent
  | IItemPurchaseEvent
  | IGrenadeThrowEvent
  | IWeaponFireEvent
  | IFlashbangEvent
  | IChatEvent
  | ITeamSwitchEvent
  | IServerCommandEvent;

/**
 * Event filter options
 */
export interface IEventFilter {
  types?: TEventType[];
  players?: string[];
  rounds?: number[];
  startTick?: number;
  endTick?: number;
}

/**
 * Event stream message (for WebSocket)
 */
export interface IEventStreamMessage {
  type: 'event' | 'status' | 'error';
  data: ISpecificGameEvent | string;
  timestamp: string;
}

/**
 * Event factory for creating events
 */
export interface IEventFactory {
  currentTick: number;
  currentRound: number;
}

/**
 * Hit group mapping
 */
export const HIT_GROUPS = {
  0: 'generic',
  1: 'head',
  2: 'chest',
  3: 'stomach',
  4: 'left_arm',
  5: 'right_arm',
  6: 'left_leg',
  7: 'right_leg',
} as const;

/**
 * Event priority for display ordering
 */
export const EVENT_PRIORITIES: Record<TEventType, number> = {
  'round_start': 10,
  'round_end': 10,
  'bomb_planted': 9,
  'bomb_defused': 9,
  'bomb_exploded': 9,
  'player_death': 8,
  'player_hurt': 6,
  'item_purchase': 4,
  'grenade_thrown': 5,
  'flashbang_detonate': 5,
  'weapon_fire': 2,
  'chat_message': 3,
  'player_connect': 7,
  'player_disconnect': 7,
  'team_switch': 7,
  'server_command': 1,
};

/**
 * Event colors for UI display
 */
export const EVENT_COLORS: Record<TEventType, string> = {
  'player_death': '#ef4444', // red-500
  'round_start': '#22c55e', // green-500
  'round_end': '#3b82f6', // blue-500
  'bomb_planted': '#f59e0b', // amber-500
  'bomb_defused': '#22c55e', // green-500
  'bomb_exploded': '#dc2626', // red-600
  'player_hurt': '#f97316', // orange-500
  'player_connect': '#10b981', // emerald-500
  'player_disconnect': '#6b7280', // gray-500
  'item_purchase': '#8b5cf6', // violet-500
  'grenade_thrown': '#f59e0b', // amber-500
  'weapon_fire': '#64748b', // slate-500
  'flashbang_detonate': '#fbbf24', // amber-400
  'chat_message': '#06b6d4', // cyan-500
  'team_switch': '#8b5cf6', // violet-500
  'server_command': '#6b7280', // gray-500
};

/**
 * Helper to determine if event is important
 */
export const isImportantEvent = (event: IGameEvent): boolean => {
  const importantTypes: TEventType[] = [
    'player_death',
    'round_start',
    'round_end',
    'bomb_planted',
    'bomb_defused',
    'bomb_exploded',
  ];
  
  return importantTypes.includes(event.type);
};

/**
 * Helper to format event for display
 */
export const formatEventForDisplay = (event: ISpecificGameEvent): string => {
  switch (event.type) {
    case 'player_death':
      const killEvent = event as IKillEvent;
      return `${killEvent.attacker.name} killed ${killEvent.victim.name} with ${killEvent.weapon}${killEvent.headshot ? ' (headshot)' : ''}`;
    
    case 'round_start':
      const roundStart = event as IRoundStartEvent;
      return `Round ${event.round} started (CT: ${roundStart.ct_score}, T: ${roundStart.t_score})`;
    
    case 'round_end':
      const roundEnd = event as IRoundEndEvent;
      return `Round ${event.round} ended - ${roundEnd.winner} wins (${roundEnd.reason})`;
    
    case 'bomb_planted':
      const plantEvent = event as IBombPlantEvent;
      return `${plantEvent.player.name} planted the bomb at site ${plantEvent.site}`;
    
    case 'bomb_defused':
      const defuseEvent = event as IBombDefuseEvent;
      return `${defuseEvent.player.name} defused the bomb${defuseEvent.with_kit ? ' (with kit)' : ''}`;
    
    case 'bomb_exploded':
      return 'The bomb has exploded';
    
    case 'player_hurt':
      const hurtEvent = event as IPlayerHurtEvent;
      return `${hurtEvent.attacker.name} damaged ${hurtEvent.victim.name} for ${hurtEvent.damage} HP with ${hurtEvent.weapon}`;
    
    case 'item_purchase':
      const purchaseEvent = event as IItemPurchaseEvent;
      return `${purchaseEvent.player.name} purchased ${purchaseEvent.item} ($${purchaseEvent.cost})`;
    
    case 'chat_message':
      const chatEvent = event as IChatEvent;
      const chatPrefix = chatEvent.team ? '[TEAM]' : '[ALL]';
      return `${chatPrefix} ${chatEvent.player?.name || 'Server'}: ${chatEvent.message}`;
    
    default:
      return `${event.type} event`;
  }
};

/**
 * Helper to get event icon
 */
export const getEventIcon = (eventType: TEventType): string => {
  const icons: Record<TEventType, string> = {
    'player_death': '☠️',
    'round_start': '🏁',
    'round_end': '🏆',
    'bomb_planted': '💣',
    'bomb_defused': '🛡️',
    'bomb_exploded': '💥',
    'player_hurt': '🩹',
    'player_connect': '🔗',
    'player_disconnect': '🔌',
    'item_purchase': '💰',
    'grenade_thrown': '🎯',
    'weapon_fire': '🔫',
    'flashbang_detonate': '💡',
    'chat_message': '💬',
    'team_switch': '🔄',
    'server_command': '⚙️',
  };
  
  return icons[eventType] || '📝';
};

/**
 * Helper to filter events
 */
export const filterEvents = (events: ISpecificGameEvent[], filter: IEventFilter): ISpecificGameEvent[] => {
  return events.filter(event => {
    if (filter.types && !filter.types.includes(event.type)) {
      return false;
    }
    
    if (filter.rounds && !filter.rounds.includes(event.round)) {
      return false;
    }
    
    if (filter.startTick !== undefined && event.tick < filter.startTick) {
      return false;
    }
    
    if (filter.endTick !== undefined && event.tick > filter.endTick) {
      return false;
    }
    
    if (filter.players && filter.players.length > 0) {
      // Check if any player in the event matches the filter
      const eventPlayers = getEventPlayers(event);
      const hasMatchingPlayer = eventPlayers.some(player => 
        filter.players!.includes(player.name)
      );
      
      if (!hasMatchingPlayer) {
        return false;
      }
    }
    
    return true;
  });
};

/**
 * Helper to extract players from an event
 */
export const getEventPlayers = (event: ISpecificGameEvent): IPlayer[] => {
  const players: IPlayer[] = [];
  
  switch (event.type) {
    case 'player_death':
      const killEvent = event as IKillEvent;
      players.push(killEvent.attacker, killEvent.victim);
      if (killEvent.assister) players.push(killEvent.assister);
      break;
      
    case 'player_hurt':
      const hurtEvent = event as IPlayerHurtEvent;
      players.push(hurtEvent.attacker, hurtEvent.victim);
      break;
      
    case 'bomb_planted':
    case 'bomb_defused':
      const bombEvent = event as IBombPlantEvent | IBombDefuseEvent;
      players.push(bombEvent.player);
      break;
      
    case 'item_purchase':
      const purchaseEvent = event as IItemPurchaseEvent;
      players.push(purchaseEvent.player);
      break;
      
    case 'chat_message':
      const chatEvent = event as IChatEvent;
      if (chatEvent.player) players.push(chatEvent.player);
      break;
      
    // Add more cases as needed
  }
  
  return players;
};