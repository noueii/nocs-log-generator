/**
 * Player-related TypeScript interfaces mirroring the Go Player model
 * Maps to backend/pkg/models/player.go
 */

import type { IVector3, IWeapon, IGrenade } from './match';
import type { TSide } from './team';

/**
 * Player role enum
 */
export type TPlayerRole = 'entry' | 'awp' | 'support' | 'igl' | 'lurker' | 'rifler';

/**
 * Valid player roles array for validation
 */
export const PLAYER_ROLES: TPlayerRole[] = ['entry', 'awp', 'support', 'igl', 'lurker', 'rifler'];

/**
 * Player statistics
 */
export interface IPlayerStats {
  kills: number;
  deaths: number;
  assists: number;
  score: number;
  damage: number;
  utility_damage: number;
  enemies_flashed: number;
  headshots: number;
  headshot_rate: number;
  accuracy: number;
  first_kills: number;
  first_deaths: number;
  trade_kills: number;
  entry_kills: number;
  '2k_rounds': number;
  '3k_rounds': number;
  '4k_rounds': number;
  '5k_rounds': number;
  bomb_plants: number;
  bomb_defuses: number;
  bomb_defuse_attempts: number;
  hostages_rescued: number;
  mvps: number;
  money_spent: number;
  grenades_thrown: Record<string, number>;
  flash_assists: number;
  team_kills: number;
  team_damage: number;
  adr: number; // Average damage per round
  kd_ratio: number;
  rating: number;
  kast: number; // Kills, Assists, Survival, Trades percentage
}

/**
 * Equipment purchase record
 */
export interface IPurchase {
  round: number;
  item: string;
  cost: number;
  timestamp?: string;
}

/**
 * Player economy state
 */
export interface IPlayerEconomy {
  money: number;
  money_spent: number;
  money_earned: number;
  equipment_value: number;
  purchases: IPurchase[];
  eco_rounds: number;
  force_buy_rounds: number;
  full_buy_rounds: number;
  economy_rating: number;
}

/**
 * Player skill and behavioral profile
 */
export interface IPlayerProfile {
  aim_skill: number;
  reflex_speed: number;
  game_sense: number;
  positioning: number;
  teamwork: number;
  utility_usage: number;
  aggression: number;
  economy_discipline: number;
  clutch_factor: number;
  rifle_skill: number;
  awp_skill: number;
  pistol_skill: number;
  entry_fragging: number;
  support_play: number;
  igl_skill: number;
  consistency_factor: number;
}

/**
 * Player current state during match
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
 * Main player interface
 */
export interface IPlayer {
  name: string;
  steam_id?: string;
  user_id?: number;
  team: string;
  side: TSide;
  role: TPlayerRole;
  state: IPlayerState;
  stats: IPlayerStats;
  economy: IPlayerEconomy;
  profile?: IPlayerProfile;
}

/**
 * Player creation request (minimal data needed)
 */
export interface ICreatePlayerRequest {
  name: string;
  steam_id?: string;
  role: TPlayerRole;
  profile?: Partial<IPlayerProfile>;
}

/**
 * Player validation result
 */
export interface IPlayerValidation {
  isValid: boolean;
  errors: string[];
}

/**
 * Default player stats
 */
export const DEFAULT_PLAYER_STATS: IPlayerStats = {
  kills: 0,
  deaths: 0,
  assists: 0,
  score: 0,
  damage: 0,
  utility_damage: 0,
  enemies_flashed: 0,
  headshots: 0,
  headshot_rate: 0,
  accuracy: 0,
  first_kills: 0,
  first_deaths: 0,
  trade_kills: 0,
  entry_kills: 0,
  '2k_rounds': 0,
  '3k_rounds': 0,
  '4k_rounds': 0,
  '5k_rounds': 0,
  bomb_plants: 0,
  bomb_defuses: 0,
  bomb_defuse_attempts: 0,
  hostages_rescued: 0,
  mvps: 0,
  money_spent: 0,
  grenades_thrown: {},
  flash_assists: 0,
  team_kills: 0,
  team_damage: 0,
  adr: 0,
  kd_ratio: 0,
  rating: 0,
  kast: 0,
};

/**
 * Default player economy
 */
export const DEFAULT_PLAYER_ECONOMY: IPlayerEconomy = {
  money: 800, // Starting money
  money_spent: 0,
  money_earned: 0,
  equipment_value: 0,
  purchases: [],
  eco_rounds: 0,
  force_buy_rounds: 0,
  full_buy_rounds: 0,
  economy_rating: 0,
};

/**
 * Default player profile (average skills)
 */
export const DEFAULT_PLAYER_PROFILE: IPlayerProfile = {
  aim_skill: 0.5,
  reflex_speed: 0.5,
  game_sense: 0.5,
  positioning: 0.5,
  teamwork: 0.5,
  utility_usage: 0.5,
  aggression: 0.5,
  economy_discipline: 0.5,
  clutch_factor: 0.5,
  rifle_skill: 0.5,
  awp_skill: 0.3,
  pistol_skill: 0.5,
  entry_fragging: 0.5,
  support_play: 0.5,
  igl_skill: 0.3,
  consistency_factor: 0.5,
};

/**
 * Default player state (spawn state)
 */
export const DEFAULT_PLAYER_STATE: IPlayerState = {
  is_alive: true,
  health: 100,
  armor: 0,
  has_helmet: false,
  has_defuse_kit: false,
  position: { x: 0, y: 0, z: 0 },
  view_angle: { x: 0, y: 0, z: 0 },
  velocity: { x: 0, y: 0, z: 0 },
  grenades: [],
  money: 800,
  is_flashed: false,
  is_smoked: false,
  is_defusing: false,
  is_planting: false,
  is_reloading: false,
  has_bomb: false,
  is_last_alive: false,
};

/**
 * Common player names for quick setup
 */
export const COMMON_PLAYER_NAMES = [
  'Player1', 'Player2', 'Player3', 'Player4', 'Player5',
  'Alpha', 'Bravo', 'Charlie', 'Delta', 'Echo',
  's1mple', 'ZywOo', 'device', 'NiKo', 'sh1ro',
  'Axle', 'Bullet', 'Crash', 'Dash', 'Edge'
] as const;

/**
 * Role-based skill profiles for realistic generation
 */
export const ROLE_PROFILES: Record<TPlayerRole, Partial<IPlayerProfile>> = {
  entry: {
    aim_skill: 0.8,
    reflex_speed: 0.9,
    aggression: 0.9,
    entry_fragging: 0.9,
    rifle_skill: 0.8,
    positioning: 0.6,
  },
  awp: {
    aim_skill: 0.9,
    awp_skill: 0.9,
    positioning: 0.8,
    aggression: 0.4,
    rifle_skill: 0.6,
    reflex_speed: 0.8,
  },
  support: {
    teamwork: 0.9,
    utility_usage: 0.9,
    support_play: 0.9,
    game_sense: 0.8,
    aggression: 0.3,
    consistency_factor: 0.8,
  },
  igl: {
    game_sense: 0.9,
    teamwork: 0.9,
    igl_skill: 0.9,
    positioning: 0.8,
    economy_discipline: 0.9,
    consistency_factor: 0.8,
  },
  lurker: {
    game_sense: 0.8,
    positioning: 0.9,
    clutch_factor: 0.8,
    aggression: 0.6,
    consistency_factor: 0.7,
    rifle_skill: 0.7,
  },
  rifler: {
    aim_skill: 0.8,
    rifle_skill: 0.9,
    consistency_factor: 0.8,
    teamwork: 0.7,
    positioning: 0.7,
    reflex_speed: 0.7,
  },
};

/**
 * Validates player data
 */
export const validatePlayer = (player: Partial<IPlayer>): IPlayerValidation => {
  const errors: string[] = [];
  
  if (!player.name || player.name.trim() === '') {
    errors.push('Player name is required');
  }
  
  if (player.role && !PLAYER_ROLES.includes(player.role)) {
    errors.push(`Invalid role: ${player.role}. Must be one of: ${PLAYER_ROLES.join(', ')}`);
  }
  
  if (player.steam_id && !isValidSteamID(player.steam_id)) {
    errors.push('Invalid SteamID format');
  }
  
  if (player.side && !['CT', 'TERRORIST'].includes(player.side)) {
    errors.push('Invalid side');
  }
  
  return {
    isValid: errors.length === 0,
    errors,
  };
};

/**
 * Validates SteamID format (STEAM_X:Y:Z)
 */
export const isValidSteamID = (steamID: string): boolean => {
  const steamIDRegex = /^STEAM_[0-1]:[0-1]:\d+$/;
  return steamIDRegex.test(steamID);
};

/**
 * Creates a player with default values
 */
export const createDefaultPlayer = (
  name: string, 
  role: TPlayerRole, 
  team: string = '', 
  side: TSide = 'CT'
): IPlayer => {
  // Apply role-based profile
  const roleProfile = ROLE_PROFILES[role] || {};
  const profile = { ...DEFAULT_PLAYER_PROFILE, ...roleProfile };
  
  return {
    name,
    role,
    team,
    side,
    state: { ...DEFAULT_PLAYER_STATE },
    stats: { ...DEFAULT_PLAYER_STATS },
    economy: { ...DEFAULT_PLAYER_ECONOMY },
    profile,
  };
};

/**
 * Generates a random SteamID for testing
 */
export const generateRandomSteamID = (): string => {
  const universe = Math.floor(Math.random() * 2);
  const type = Math.floor(Math.random() * 2);
  const instance = Math.floor(Math.random() * 1000000000);
  
  return `STEAM_${universe}:${type}:${instance}`;
};

/**
 * Helper to get player by name
 */
export const getPlayerByName = (players: IPlayer[], name: string): IPlayer | undefined => {
  return players.find(player => 
    player.name.toLowerCase() === name.toLowerCase()
  );
};

/**
 * Helper to check if player is alive
 */
export const isPlayerAlive = (player: IPlayer): boolean => {
  return player.state.is_alive && player.state.health > 0;
};

/**
 * Helper to calculate K/D ratio
 */
export const calculateKDRatio = (kills: number, deaths: number): number => {
  return deaths === 0 ? kills : Number((kills / deaths).toFixed(2));
};

/**
 * Helper to calculate headshot percentage
 */
export const calculateHeadshotRate = (headshots: number, kills: number): number => {
  return kills === 0 ? 0 : Number(((headshots / kills) * 100).toFixed(1));
};