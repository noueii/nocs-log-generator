/**
 * TypeScript type definitions index
 * 
 * File naming conventions:
 * - Files: camelCase.ts (match.ts, player.ts, etc.)
 * - Interfaces: PascalCase with I prefix (IMatch, IPlayer, etc.)
 * - Types: PascalCase with T prefix (TPlayerRole, TMatchFormat, etc.)
 * - Enums: PascalCase (PlayerRole, MatchStatus, etc.)
 */

// Match-related types
export type {
  IMatch,
  IMatchConfig,
  IMatchOptions,
  IMatchState,
  IGenerateRequest,
  IGenerateResponse,
  IRoundData,
  ITeamEconomy as IMatchTeamEconomy,
  IPlayerState,
  IVector3,
  IWeapon,
  IGrenade,
  TWeaponType,
  TGrenadeType,
  TMatchStatus,
  TMatchFormat,
  TRoundEndReason,
  TMapName,
} from './match';

export {
  MAPS,
  DEFAULT_MATCH_CONFIG,
} from './match';

// Team-related types
export type {
  ITeam,
  ITeamEconomy,
  ITeamStats,
  ICreateTeamRequest,
  ITeamValidation,
  TSide,
} from './team';

export {
  DEFAULT_TEAM_ECONOMY,
  DEFAULT_TEAM_STATS,
  COMMON_TEAM_NAMES,
  createDefaultTeam,
  validateTeam,
  getTeamBySide,
  getOppositeSide,
} from './team';

// Player-related types
export type {
  IPlayer,
  IPlayerStats,
  IPlayerEconomy,
  IPlayerProfile,
  ICreatePlayerRequest,
  IPlayerValidation,
  IPurchase,
  TPlayerRole,
} from './player';

export {
  PLAYER_ROLES,
  DEFAULT_PLAYER_STATS,
  DEFAULT_PLAYER_ECONOMY,
  DEFAULT_PLAYER_PROFILE,
  DEFAULT_PLAYER_STATE,
  COMMON_PLAYER_NAMES,
  ROLE_PROFILES,
  validatePlayer,
  isValidSteamID,
  createDefaultPlayer,
  generateRandomSteamID,
  getPlayerByName,
  isPlayerAlive,
  calculateKDRatio,
  calculateHeadshotRate,
} from './player';

// Event-related types
export type {
  IGameEvent,
  IKillEvent,
  IRoundStartEvent,
  IRoundEndEvent,
  IBombPlantEvent,
  IBombDefuseEvent,
  IBombExplodeEvent,
  IPlayerHurtEvent,
  IPlayerConnectEvent,
  IPlayerDisconnectEvent,
  IItemPurchaseEvent,
  IGrenadeThrowEvent,
  IWeaponFireEvent,
  IFlashbangEvent,
  IChatEvent,
  ITeamSwitchEvent,
  IServerCommandEvent,
  ISpecificGameEvent,
  IEventFilter,
  IEventStreamMessage,
  IEventFactory,
  TEventType,
} from './events';

export {
  HIT_GROUPS,
  EVENT_PRIORITIES,
  EVENT_COLORS,
  isImportantEvent,
  formatEventForDisplay,
  getEventIcon,
  filterEvents,
  getEventPlayers,
} from './events';

// Re-export commonly used type combinations
import type {
  IMatch,
  IMatchConfig,
  IRoundData,
} from './match';
import type {
  ITeam,
  ICreateTeamRequest,
} from './team';
import type {
  IPlayer,
} from './player';
import type {
  ISpecificGameEvent,
} from './events';

export type IFullMatch = IMatch & {
  teams: ITeam[];
  rounds?: IRoundData[];
  events?: ISpecificGameEvent[];
};

export type IMatchWithPlayers = IMatch & {
  teams: (ITeam & { players: IPlayer[] })[];
};

// Utility types for forms and API
export type IMatchFormData = {
  config: IMatchConfig;
  teams: ICreateTeamRequest[];
};

export type IPlayerFormData = Omit<IPlayer, 'state' | 'stats' | 'economy' | 'team' | 'side'>;

export type ITeamFormData = Omit<ITeam, 'score' | 'rounds_won' | 'economy' | 'stats' | 'side'> & {
  players: IPlayerFormData[];
};