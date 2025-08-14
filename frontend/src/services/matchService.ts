/**
 * Match service for type-safe API calls
 * Handles all match-related API operations with proper error handling and TypeScript types
 */

import { get, post, apiCallWithOfflineHandling } from './api';
import type { 
  IMatch, 
  IGenerateRequest, 
  IGenerateResponse, 
  IMatchConfig,
  TMapName,
  TMatchFormat
} from '../types/match';
import {
  DEFAULT_MATCH_CONFIG,
  MAPS,
} from '../types/match';
import type { ITeam } from '../types/team';

/**
 * Match list response
 */
export interface IMatchListResponse {
  matches: IMatch[];
  total: number;
  page: number;
  pageSize: number;
}

/**
 * Match search/filter parameters
 */
export interface IMatchSearchParams {
  page?: number;
  pageSize?: number;
  status?: string[];
  format?: TMatchFormat[];
  maps?: TMapName[];
  startDate?: string;
  endDate?: string;
  search?: string;
}

/**
 * Match templates response
 */
export interface IMatchTemplate {
  id: string;
  name: string;
  description: string;
  config: IMatchConfig;
  teams: ITeam[];
}

/**
 * Available maps response
 */
export interface IMapInfo {
  name: TMapName;
  display_name: string;
  thumbnail?: string;
  description?: string;
  bomb_sites: string[];
  spawn_points: {
    ct: number;
    t: number;
  };
}

/**
 * Server status response
 */
export interface IServerStatus {
  status: 'healthy' | 'unhealthy';
  version: string;
  uptime: number;
  active_matches: number;
  queue_length: number;
  last_updated: string;
}

/**
 * Match generation progress (for WebSocket)
 */
export interface IGenerationProgress {
  match_id: string;
  status: 'generating' | 'completed' | 'error';
  progress: number; // 0-100
  current_round?: number;
  total_rounds?: number;
  events_generated?: number;
  message?: string;
  error?: string;
}

/**
 * Match Service Class
 */
class MatchService {
  /**
   * Generate a new match
   */
  async generateMatch(request: IGenerateRequest): Promise<IGenerateResponse> {
    return apiCallWithOfflineHandling(
      () => post<IGenerateResponse>('/api/v1/generate', request),
      () => this.generateMatchOffline(request)
    );
  }

  /**
   * Get match by ID
   */
  async getMatch(matchId: string): Promise<IMatch> {
    return apiCallWithOfflineHandling(
      () => get<IMatch>(`/api/v1/matches/${matchId}`),
      () => this.getMatchOffline(matchId)
    );
  }

  /**
   * Get list of matches with pagination and filtering
   */
  async getMatches(params: IMatchSearchParams = {}): Promise<IMatchListResponse> {
    const searchParams = {
      page: params.page || 1,
      pageSize: params.pageSize || 20,
      ...params,
    };
    
    return apiCallWithOfflineHandling(
      () => get<IMatchListResponse>('/api/v1/matches', searchParams),
      () => this.getMatchesOffline(params)
    );
  }

  /**
   * Delete a match
   */
  async deleteMatch(matchId: string): Promise<void> {
    try {
      await post<void>(`/api/v1/matches/${matchId}/delete`);
    } catch (error) {
      console.error(`Failed to delete match ${matchId}:`, error);
      throw error;
    }
  }

  /**
   * Download match log
   */
  async downloadMatchLog(matchId: string): Promise<Blob> {
    try {
      // Note: This would typically use a different method to download binary data
      const response = await get<Blob>(`/api/v1/matches/${matchId}/download`);
      return response;
    } catch (error) {
      console.error(`Failed to download match log ${matchId}:`, error);
      throw error;
    }
  }

  /**
   * Get match templates
   */
  async getMatchTemplates(): Promise<IMatchTemplate[]> {
    return apiCallWithOfflineHandling(
      async () => {
        const response = await get<{templates: Record<string, any>}>('/api/v1/config/templates');
        // Transform backend response to frontend format
        return Object.entries(response.templates).map(([id, config]) => ({
          id,
          name: id.charAt(0).toUpperCase() + id.slice(1),
          description: `${id.charAt(0).toUpperCase() + id.slice(1)} match template`,
          config,
          teams: [] // Empty teams, filled by user
        }));
      },
      () => Promise.resolve(this.getDefaultTemplates())
    );
  }

  /**
   * Get available maps
   */
  async getMaps(): Promise<IMapInfo[]> {
    return apiCallWithOfflineHandling(
      async () => {
        const response = await get<{maps: Array<{name: string, display_name: string, type: string}>}>('/api/v1/config/maps');
        return response.maps.map(map => ({
          name: map.name as TMapName,
          display_name: map.display_name,
          bomb_sites: ['A', 'B'],
          spawn_points: {
            ct: 5,
            t: 5
          }
        }));
      },
      () => Promise.resolve(this.getDefaultMaps())
    );
  }

  /**
   * Get server status
   */
  async getServerStatus(): Promise<IServerStatus> {
    return apiCallWithOfflineHandling(
      async () => {
        const response = await get<{status: string, checks: any}>('/ready');
        return {
          status: response.status === 'ready' ? 'healthy' : 'unhealthy',
          version: '0.1.0',
          uptime: 0,
          active_matches: 0,
          queue_length: 0,
          last_updated: new Date().toISOString()
        };
      },
      async () => ({
        status: 'unhealthy' as const,
        version: '0.1.0',
        uptime: 0,
        active_matches: 0,
        queue_length: 0,
        last_updated: new Date().toISOString()
      })
    );
  }

  /**
   * Validate match configuration
   */
  async validateMatchConfig(config: IMatchConfig): Promise<{ isValid: boolean; errors: string[] }> {
    try {
      const response = await post<{ isValid: boolean; errors: string[] }>(
        '/api/v1/validate/config', 
        config
      );
      return response;
    } catch (error) {
      console.error('Failed to validate match config:', error);
      
      // Return client-side validation as fallback
      return this.validateConfigClientSide(config);
    }
  }

  /**
   * Validate teams
   */
  async validateTeams(teams: ITeam[]): Promise<{ isValid: boolean; errors: string[] }> {
    try {
      const response = await post<{ isValid: boolean; errors: string[] }>(
        '/api/v1/validate/teams', 
        teams
      );
      return response;
    } catch (error) {
      console.error('Failed to validate teams:', error);
      
      // Return client-side validation as fallback
      return this.validateTeamsClientSide(teams);
    }
  }

  /**
   * Get match statistics
   */
  async getMatchStatistics(matchId: string): Promise<any> {
    try {
      const response = await get<any>(`/api/v1/matches/${matchId}/stats`);
      return response;
    } catch (error) {
      console.error(`Failed to get match statistics ${matchId}:`, error);
      throw error;
    }
  }

  // Private helper methods

  /**
   * Generate match offline (fallback)
   */
  private async generateMatchOffline(request: IGenerateRequest): Promise<IGenerateResponse> {
    // Create mock response for offline mode
    const matchId = `offline_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    
    // Simulate processing delay
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    return {
      match_id: matchId,
      status: 'completed',
      log_url: '',
      message: 'Generated in offline mode - no actual log file available'
    };
  }

  /**
   * Get match offline (fallback)
   */
  private async getMatchOffline(matchId: string): Promise<IMatch> {
    throw new Error(`Cannot retrieve match ${matchId} - offline mode`);
  }

  /**
   * Get matches offline (fallback)
   */
  private async getMatchesOffline(params: IMatchSearchParams): Promise<IMatchListResponse> {
    return {
      matches: [],
      total: 0,
      page: params.page || 1,
      pageSize: params.pageSize || 20
    };
  }

  /**
   * Get default templates (fallback)
   */
  private getDefaultTemplates(): IMatchTemplate[] {
    return [
      {
        id: 'competitive-mr12',
        name: 'Competitive MR12',
        description: 'Standard competitive match with MR12 format',
        config: {
          ...DEFAULT_MATCH_CONFIG,
          format: 'mr12',
          realistic_economy: true,
          skill_variance: 0.15,
        },
        teams: [], // Would be populated from user input
      },
      {
        id: 'casual-mr15',
        name: 'Casual MR15',
        description: 'Casual match with MR15 format and relaxed settings',
        config: {
          ...DEFAULT_MATCH_CONFIG,
          format: 'mr15',
          realistic_economy: false,
          skill_variance: 0.25,
          chat_messages: true,
        },
        teams: [],
      },
      {
        id: 'testing-verbose',
        name: 'Testing Match',
        description: 'Verbose logging for testing and debugging',
        config: {
          ...DEFAULT_MATCH_CONFIG,
          output_verbosity: 'verbose',
          include_positions: true,
          include_weapon_fire: true,
          rollback_enabled: true,
          rollback_probability: 0.1,
        },
        teams: [],
      },
    ];
  }

  /**
   * Get default maps (fallback)
   */
  private getDefaultMaps(): IMapInfo[] {
    return MAPS.map(mapName => ({
      name: mapName,
      display_name: mapName.replace('de_', '').replace(/^\w/, c => c.toUpperCase()),
      bomb_sites: ['A', 'B'],
      spawn_points: {
        ct: 5,
        t: 5,
      },
    }));
  }

  /**
   * Client-side config validation (fallback)
   */
  private validateConfigClientSide(config: IMatchConfig): { isValid: boolean; errors: string[] } {
    const errors: string[] = [];

    if (!config.format || !['mr12', 'mr15'].includes(config.format)) {
      errors.push('Format must be either mr12 or mr15');
    }

    if (!config.map || config.map.trim() === '') {
      errors.push('Map is required');
    }

    if (config.tick_rate && (config.tick_rate < 64 || config.tick_rate > 128)) {
      errors.push('Tick rate must be between 64 and 128');
    }

    if (config.skill_variance < 0 || config.skill_variance > 1) {
      errors.push('Skill variance must be between 0 and 1');
    }

    if (config.rollback_probability < 0 || config.rollback_probability > 1) {
      errors.push('Rollback probability must be between 0 and 1');
    }

    if (config.start_money < 0 || config.start_money > config.max_money) {
      errors.push('Start money must be between 0 and max money');
    }

    return {
      isValid: errors.length === 0,
      errors,
    };
  }

  /**
   * Client-side teams validation (fallback)
   */
  private validateTeamsClientSide(teams: ITeam[]): { isValid: boolean; errors: string[] } {
    const errors: string[] = [];

    if (!teams || teams.length !== 2) {
      errors.push('Exactly 2 teams are required');
      return { isValid: false, errors };
    }

    teams.forEach((team, index) => {
      if (!team.name || team.name.trim() === '') {
        errors.push(`Team ${index + 1} name is required`);
      }

      if (!team.players || team.players.length !== 5) {
        errors.push(`Team ${index + 1} must have exactly 5 players`);
      }

      if (team.players) {
        const playerNames = team.players.map(p => p.name.toLowerCase());
        const duplicateNames = playerNames.filter((name, i) => 
          playerNames.indexOf(name) !== i
        );
        
        if (duplicateNames.length > 0) {
          errors.push(`Team ${index + 1} has duplicate player names`);
        }

        team.players.forEach((player, playerIndex) => {
          if (!player.name || player.name.trim() === '') {
            errors.push(`Team ${index + 1}, Player ${playerIndex + 1} name is required`);
          }
        });
      }
    });

    // Check for duplicate team names
    if (teams.length === 2 && teams[0].name.toLowerCase() === teams[1].name.toLowerCase()) {
      errors.push('Teams must have different names');
    }

    return {
      isValid: errors.length === 0,
      errors,
    };
  }
}

/**
 * Export singleton instance
 */
export const matchService = new MatchService();

/**
 * Export individual functions for tree-shaking
 */
export const {
  generateMatch,
  getMatch,
  getMatches,
  deleteMatch,
  downloadMatchLog,
  getMatchTemplates,
  getMaps,
  getServerStatus,
  validateMatchConfig,
  validateTeams,
  getMatchStatistics,
} = matchService;