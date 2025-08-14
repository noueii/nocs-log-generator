/**
 * Match-related TypeScript interfaces mirroring the Go Match model
 * Maps to backend/pkg/models/match.go
 */

import type { ITeam } from './team';
import type { IGameEvent } from './events';
import type { IPlayer } from './player';

/**
 * Match configuration options
 */
export interface IMatchOptions {
  seed?: number;
  tick_rate?: number;
  overtime?: boolean;
  max_rounds?: number;
}

/**
 * Request body for match generation
 */
export interface IGenerateRequest {
  teams: ITeam[];
  map: string;
  format: 'mr12' | 'mr15';
  options?: IMatchOptions;
}

/**
 * Response from match generation
 */
export interface IGenerateResponse {
  match_id: string;
  status: string;
  log_url?: string;
  error?: string;
}

/**
 * Team economy state for round data
 */
export interface ITeamEconomy {
  total_money: number;
  average_money: number;
  equipment_value: number;
  consecutive_losses: number;
  loss_bonus: number;
  money_spent: number;
  money_earned: number;
  rifles: number;
  smgs: number;
  pistols: number;
  snipers: number;
  grenades: number;
  armor: number;
  helmets: number;
  defuse_kits: number;
}

/**
 * Round data containing state and events for a single round
 */
export interface IRoundData {
  round_number: number;
  start_time: string;
  end_time: string;
  winner: 'CT' | 'TERRORIST';
  reason: 'elimination' | 'bomb_defused' | 'bomb_exploded' | 'time';
  mvp?: IPlayer;
  events: IGameEvent[];
  economy: Record<string, ITeamEconomy>;
  scores: Record<string, number>;
}

/**
 * Match configuration
 */
export interface IMatchConfig {
  format: 'mr12' | 'mr15';
  map: string;
  overtime: boolean;
  max_rounds?: number;
  tick_rate: number;
  server_name?: string;
  seed?: number;
  duration?: number;
  rollback_enabled: boolean;
  rollback_probability: number;
  rollback_min_round: number;
  rollback_max_round: number;
  start_money: number;
  max_money: number;
  realistic_economy: boolean;
  network_issues: boolean;
  anti_cheat_events: boolean;
  chat_messages: boolean;
  skill_variance: number;
  log_format: string;
  timestamp_format: string;
  output_verbosity: 'minimal' | 'standard' | 'verbose';
  include_positions: boolean;
  include_weapon_fire: boolean;
}

/**
 * Main match interface
 */
export interface IMatch {
  id: string;
  title?: string;
  map: string;
  format: 'mr12' | 'mr15';
  status: 'pending' | 'generating' | 'completed' | 'error';
  start_time?: string;
  end_time?: string;
  log_url?: string;
  error?: string;
  config: IMatchConfig;
  teams: ITeam[];
  current_round: number;
  max_rounds: number;
  overtime: boolean;
  scores: Record<string, number>;
  rounds?: IRoundData[];
  events?: IGameEvent[];
  total_events: number;
  file_size?: number;
  duration?: number;
}

/**
 * Match state during generation
 */
export interface IMatchState {
  current_round: number;
  scores: Record<string, number>;
  team_economies: Record<string, ITeamEconomy>;
  player_states: Record<string, IPlayerState>;
  bomb_carrier?: IPlayer;
  is_live: boolean;
  is_freeze_time: boolean;
  round_start_time: string;
  current_tick: number;
}

/**
 * Player state interface from player model
 */
export interface IPlayerState {
  is_alive: boolean;
  health: number;
  armor: number;
  has_helmet: boolean;
  has_defuse_kit: boolean;
  position: IVector3;
  view_angle: IVector3;
  velocity: IVector3;
  primary_weapon?: IWeapon;
  secondary_weapon?: IWeapon;
  grenades: IGrenade[];
  money: number;
  is_flashed: boolean;
  is_smoked: boolean;
  is_defusing: boolean;
  is_planting: boolean;
  is_reloading: boolean;
  has_bomb: boolean;
  is_last_alive: boolean;
}

/**
 * 3D vector for positions and directions
 */
export interface IVector3 {
  x: number;
  y: number;
  z: number;
}

/**
 * Weapon interface
 */
export interface IWeapon {
  name: string;
  type: TWeaponType;
  damage: number;
  accuracy: number;
  range_modifier: number;
  penetration_power: number;
  price: number;
  ammo: number;
  ammo_reserve: number;
  max_ammo: number;
  skin?: string;
  stat_trak: boolean;
}

/**
 * Grenade interface
 */
export interface IGrenade {
  type: TGrenadeType;
  price: number;
  damage?: number;
  effect_radius?: number;
  duration?: number;
}

/**
 * Weapon type enum
 */
export type TWeaponType = 'rifle' | 'pistol' | 'sniper' | 'smg' | 'shotgun' | 'machinegun';

/**
 * Grenade type enum  
 */
export type TGrenadeType = 'he' | 'flash' | 'smoke' | 'incendiary' | 'molotov' | 'decoy';

/**
 * Match status type
 */
export type TMatchStatus = 'pending' | 'generating' | 'completed' | 'error';

/**
 * Match format type
 */
export type TMatchFormat = 'mr12' | 'mr15';

/**
 * Round end reason type
 */
export type TRoundEndReason = 'elimination' | 'bomb_defused' | 'bomb_exploded' | 'time';

/**
 * Available CS2 maps
 */
export const MAPS = [
  'de_mirage',
  'de_dust2', 
  'de_inferno',
  'de_cache',
  'de_overpass',
  'de_train',
  'de_nuke',
  'de_cbble',
  'de_vertigo',
  'de_ancient'
] as const;

export type TMapName = typeof MAPS[number];

/**
 * Default match configuration
 */
export const DEFAULT_MATCH_CONFIG: IMatchConfig = {
  format: 'mr12',
  map: 'de_mirage',
  overtime: false,
  tick_rate: 64,
  start_money: 800,
  max_money: 16000,
  realistic_economy: true,
  rollback_enabled: false,
  rollback_probability: 0.0,
  rollback_min_round: 0,
  rollback_max_round: 0,
  network_issues: false,
  anti_cheat_events: false,
  chat_messages: true,
  skill_variance: 0.15,
  log_format: 'standard',
  timestamp_format: '01/02/2006 - 15:04:05',
  output_verbosity: 'standard',
  include_positions: false,
  include_weapon_fire: false,
};