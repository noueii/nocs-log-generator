/**
 * Team-related TypeScript interfaces mirroring the Go Team model
 * Maps to backend/pkg/models/team.go
 */

import type { IPlayer } from './player';

/**
 * Team side enum
 */
export type TSide = 'CT' | 'TERRORIST';

/**
 * Team economy state
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
 * Team statistics
 */
export interface ITeamStats {
  kills: number;
  deaths: number;
  assists: number;
  damage: number;
  headshots: number;
  headshot_rate: number;
  first_kills: number;
  first_deaths: number;
  bomb_plants: number;
  bomb_defuses: number;
  rounds_played: number;
  rounds_won_ct: number;
  rounds_won_t: number;
  money_per_round: number;
  economy_rating: number;
}

/**
 * Main team interface
 */
export interface ITeam {
  name: string;
  tag?: string;
  country?: string;
  ranking?: number;
  players: IPlayer[];
  side: TSide;
  score: number;
  rounds_won: number;
  economy: ITeamEconomy;
  stats: ITeamStats;
}

/**
 * Team creation request
 */
export interface ICreateTeamRequest {
  name: string;
  tag?: string;
  country?: string;
  ranking?: number;
  players: Omit<IPlayer, 'team' | 'side' | 'state' | 'stats' | 'economy'>[];
}

/**
 * Team validation result
 */
export interface ITeamValidation {
  isValid: boolean;
  errors: string[];
}

/**
 * Default team economy
 */
export const DEFAULT_TEAM_ECONOMY: ITeamEconomy = {
  total_money: 4000, // 5 players * 800 starting money
  average_money: 800,
  equipment_value: 0,
  consecutive_losses: 0,
  loss_bonus: 1400,
  money_spent: 0,
  money_earned: 0,
  rifles: 0,
  smgs: 0,
  pistols: 5, // Starting pistols
  snipers: 0,
  grenades: 0,
  armor: 0,
  helmets: 0,
  defuse_kits: 0,
};

/**
 * Default team stats
 */
export const DEFAULT_TEAM_STATS: ITeamStats = {
  kills: 0,
  deaths: 0,
  assists: 0,
  damage: 0,
  headshots: 0,
  headshot_rate: 0,
  first_kills: 0,
  first_deaths: 0,
  bomb_plants: 0,
  bomb_defuses: 0,
  rounds_played: 0,
  rounds_won_ct: 0,
  rounds_won_t: 0,
  money_per_round: 0,
  economy_rating: 0,
};

/**
 * Creates a default team structure
 */
export const createDefaultTeam = (name: string, side: TSide): Omit<ITeam, 'players'> => ({
  name,
  side,
  score: 0,
  rounds_won: 0,
  economy: { ...DEFAULT_TEAM_ECONOMY },
  stats: { ...DEFAULT_TEAM_STATS },
});

/**
 * Common team names for quick setup
 */
export const COMMON_TEAM_NAMES = [
  { name: 'Team Alpha', tag: 'ALPH' },
  { name: 'Team Beta', tag: 'BETA' },
  { name: 'Counter-Strike Champions', tag: 'CSC' },
  { name: 'Elite Squad', tag: 'ELIT' },
  { name: 'Phoenix Gaming', tag: 'PHX' },
  { name: 'Vanguard', tag: 'VAN' },
  { name: 'Titan Esports', tag: 'TTN' },
  { name: 'Storm Riders', tag: 'STRM' },
] as const;

/**
 * Team validation utility
 */
export const validateTeam = (team: Partial<ITeam>): ITeamValidation => {
  const errors: string[] = [];
  
  if (!team.name || team.name.trim() === '') {
    errors.push('Team name is required');
  }
  
  if (!team.players || team.players.length !== 5) {
    errors.push('Team must have exactly 5 players');
  }
  
  if (team.players) {
    // Check for duplicate player names
    const playerNames = team.players.map(p => p.name.toLowerCase());
    const duplicateNames = playerNames.filter((name, index) => 
      playerNames.indexOf(name) !== index
    );
    
    if (duplicateNames.length > 0) {
      errors.push(`Duplicate player names: ${duplicateNames.join(', ')}`);
    }
    
    // Validate each player
    team.players.forEach((player, index) => {
      if (!player.name || player.name.trim() === '') {
        errors.push(`Player ${index + 1} name is required`);
      }
    });
  }
  
  if (team.side && !['CT', 'TERRORIST'].includes(team.side)) {
    errors.push('Invalid team side');
  }
  
  return {
    isValid: errors.length === 0,
    errors,
  };
};

/**
 * Helper to get team by side
 */
export const getTeamBySide = (teams: ITeam[], side: TSide): ITeam | undefined => {
  return teams.find(team => team.side === side);
};

/**
 * Helper to get opposing side
 */
export const getOppositeSide = (side: TSide): TSide => {
  return side === 'CT' ? 'TERRORIST' : 'CT';
};