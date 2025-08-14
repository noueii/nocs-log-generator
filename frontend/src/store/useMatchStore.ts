/**
 * Match Store for managing match history and data
 * Handles local storage of matches and provides offline capabilities
 */

import { create } from 'zustand';
import { devtools, persist } from 'zustand/middleware';
import { immer } from 'zustand/middleware/immer';
import type { IMatch, IGenerateResponse } from '@/types/match';

/**
 * Match Store state interface
 */
export interface IMatchStore {
  // State
  matches: IMatch[];
  isLoading: boolean;
  error: string | null;
  
  // Actions
  addMatch: (match: IMatch) => void;
  removeMatch: (matchId: string) => void;
  updateMatch: (matchId: string, updates: Partial<IMatch>) => void;
  getMatch: (matchId: string) => IMatch | undefined;
  getMatchHistory: () => IMatch[];
  clearHistory: () => void;
  
  // Generation result handling
  saveGenerationResult: (result: IGenerateResponse, matchData?: Partial<IMatch>) => IMatch;
  
  // Utility functions
  getMatchesByStatus: (status: string) => IMatch[];
  getMatchesByFormat: (format: string) => IMatch[];
  searchMatches: (query: string) => IMatch[];
}

/**
 * Create a match from generation result
 */
const createMatchFromResult = (
  result: IGenerateResponse, 
  matchData?: Partial<IMatch>
): IMatch => {
  const now = new Date();
  
  return {
    id: result.match_id,
    status: result.status as any,
    map: matchData?.map || 'de_mirage',
    format: matchData?.format || 'mr12',
    teams: matchData?.teams || [
      {
        name: 'Team 1',
        tag: 'T1',
        side: 'CT',
        score: 0,
        rounds_won: 0,
        players: [],
        economy: {
          total_money: 16000,
          average_money: 800,
          equipment_value: 0,
          consecutive_losses: 0,
          loss_bonus: 1400,
          money_spent: 0,
          money_earned: 16000,
          rifles: 0,
          smgs: 0,
          pistols: 5,
          snipers: 0,
          grenades: 0,
          armor: 0,
          helmets: 0,
          defuse_kits: 0
        },
        stats: {
          kills: 0,
          deaths: 0,
          assists: 0,
          score: 0,
          mvps: 0,
          adr: 0,
          first_kills: 0,
          first_deaths: 0,
          clutch_wins: 0,
          total_damage: 0
        }
      },
      {
        name: 'Team 2',
        tag: 'T2',
        side: 'TERRORIST',
        score: 0,
        rounds_won: 0,
        players: [],
        economy: {
          total_money: 16000,
          average_money: 800,
          equipment_value: 0,
          consecutive_losses: 0,
          loss_bonus: 1400,
          money_spent: 0,
          money_earned: 16000,
          rifles: 0,
          smgs: 0,
          pistols: 5,
          snipers: 0,
          grenades: 0,
          armor: 0,
          helmets: 0,
          defuse_kits: 0
        },
        stats: {
          kills: 0,
          deaths: 0,
          assists: 0,
          score: 0,
          mvps: 0,
          adr: 0,
          first_kills: 0,
          first_deaths: 0,
          clutch_wins: 0,
          total_damage: 0
        }
      }
    ],
    current_round: 0,
    max_rounds: matchData?.format === 'mr15' ? 30 : 24,
    overtime: false,
    start_time: now.toISOString(),
    end_time: result.status === 'completed' ? now.toISOString() : undefined,
    duration: result.status === 'completed' ? 2700 : undefined, // 45 minutes default
    total_events: 0,
    file_size: 0,
    log_url: result.log_url || '',
    server_info: {
      hostname: 'CS2 Log Generator',
      ip: '127.0.0.1',
      port: 27015,
      version: '1.0.0',
      tick_rate: 64,
      map_group: 'mg_active',
      game_mode: 'competitive'
    },
    game_state: {
      phase: result.status === 'completed' ? 'finished' : 'warmup',
      round_start_time: now.toISOString(),
      round_end_time: undefined,
      bomb_planted: false,
      bomb_defused: false,
      round_winner: undefined,
      round_end_reason: undefined,
      money_restart: false,
      freeze_time: true
    },
    scores: {
      [matchData?.teams?.[0]?.name || 'Team 1']: 0,
      [matchData?.teams?.[1]?.name || 'Team 2']: 0,
    },
    rounds: [],
    chat_messages: [],
    ...matchData,
  };
};

/**
 * Create the match store
 */
export const useMatchStore = create<IMatchStore>()(
  devtools(
    persist(
      immer((set, get) => ({
        // Initial state
        matches: [],
        isLoading: false,
        error: null,
        
        // Actions
        addMatch: (match) => {
          set((state) => {
            // Check if match already exists
            const existingIndex = state.matches.findIndex(m => m.id === match.id);
            
            if (existingIndex >= 0) {
              // Update existing match
              state.matches[existingIndex] = match;
            } else {
              // Add new match to the beginning
              state.matches.unshift(match);
            }
            
            // Keep only last 100 matches to prevent storage bloat
            if (state.matches.length > 100) {
              state.matches = state.matches.slice(0, 100);
            }
          });
        },
        
        removeMatch: (matchId) => {
          set((state) => {
            state.matches = state.matches.filter(m => m.id !== matchId);
          });
        },
        
        updateMatch: (matchId, updates) => {
          set((state) => {
            const matchIndex = state.matches.findIndex(m => m.id === matchId);
            if (matchIndex >= 0) {
              Object.assign(state.matches[matchIndex], updates);
            }
          });
        },
        
        getMatch: (matchId) => {
          return get().matches.find(m => m.id === matchId);
        },
        
        getMatchHistory: () => {
          return get().matches;
        },
        
        clearHistory: () => {
          set((state) => {
            state.matches = [];
          });
        },
        
        saveGenerationResult: (result, matchData) => {
          const match = createMatchFromResult(result, matchData);
          get().addMatch(match);
          return match;
        },
        
        // Utility functions
        getMatchesByStatus: (status) => {
          return get().matches.filter(m => m.status === status);
        },
        
        getMatchesByFormat: (format) => {
          return get().matches.filter(m => m.format === format);
        },
        
        searchMatches: (query) => {
          const lowercaseQuery = query.toLowerCase();
          return get().matches.filter(match => 
            match.id.toLowerCase().includes(lowercaseQuery) ||
            match.map.toLowerCase().includes(lowercaseQuery) ||
            match.teams.some(team => 
              team.name.toLowerCase().includes(lowercaseQuery) ||
              team.players.some(player => 
                player.name.toLowerCase().includes(lowercaseQuery)
              )
            )
          );
        },
      }))
      ,
      {
        name: 'match-store-persist',
        // Only persist matches data
        partialize: (state) => ({
          matches: state.matches,
        }),
      }
    ),
    {
      name: 'match-store',
    }
  )
);

/**
 * Hook for match statistics
 */
export const useMatchStatistics = () => {
  const matches = useMatchStore(state => state.matches);
  
  return React.useMemo(() => {
    const totalMatches = matches.length;
    const completedMatches = matches.filter(m => m.status === 'completed').length;
    const errorMatches = matches.filter(m => m.status === 'error').length;
    
    const formatStats = matches.reduce((acc, match) => {
      acc[match.format] = (acc[match.format] || 0) + 1;
      return acc;
    }, {} as Record<string, number>);
    
    const mapStats = matches.reduce((acc, match) => {
      acc[match.map] = (acc[match.map] || 0) + 1;
      return acc;
    }, {} as Record<string, number>);
    
    return {
      totalMatches,
      completedMatches,
      errorMatches,
      successRate: totalMatches > 0 ? (completedMatches / totalMatches) * 100 : 0,
      formatStats,
      mapStats,
      recentMatches: matches.slice(0, 5),
    };
  }, [matches]);
};

// Import React for the hook
import React from 'react';

/**
 * Hook for match history (convenience hook)
 */
export const useMatchHistory = () => {
  return useMatchStore(state => ({
    matches: state.matches,
    isLoading: state.isLoading,
    error: state.error,
    addMatch: state.addMatch,
    removeMatch: state.removeMatch,
    clearHistory: state.clearHistory,
    searchMatches: state.searchMatches,
  }));
};

/**
 * Hook for match filters (for filter UI)
 */
export const useMatchFilters = () => {
  const [filters, setFilters] = React.useState<IMatchFilters>({
    status: 'all',
    format: 'all',
    map: 'all',
    dateRange: 'all',
    searchQuery: '',
  });

  return {
    filters,
    setFilters,
    resetFilters: () => setFilters({
      status: 'all',
      format: 'all',
      map: 'all',
      dateRange: 'all',
      searchQuery: '',
    }),
  };
};

/**
 * Hook for match selection (for bulk operations)
 */
export const useMatchSelection = () => {
  const [selectedIds, setSelectedIds] = React.useState<Set<string>>(new Set());

  return {
    selectedIds,
    selectMatch: (id: string) => {
      setSelectedIds(prev => new Set([...prev, id]));
    },
    deselectMatch: (id: string) => {
      setSelectedIds(prev => {
        const next = new Set(prev);
        next.delete(id);
        return next;
      });
    },
    toggleMatch: (id: string) => {
      setSelectedIds(prev => {
        const next = new Set(prev);
        if (next.has(id)) {
          next.delete(id);
        } else {
          next.add(id);
        }
        return next;
      });
    },
    selectAll: (ids: string[]) => {
      setSelectedIds(new Set(ids));
    },
    clearSelection: () => {
      setSelectedIds(new Set());
    },
    isSelected: (id: string) => selectedIds.has(id),
  };
};

// Type exports
export type IMatchHistoryItem = IMatch;
export type IMatchFilters = {
  status: 'all' | 'completed' | 'in_progress' | 'error';
  format: 'all' | 'mr12' | 'mr15';
  map: string;
  dateRange: 'all' | 'today' | 'week' | 'month';
  searchQuery: string;
};
export type TMatchSortBy = 'date' | 'map' | 'format' | 'status' | 'duration';
export type TMatchSortOrder = 'asc' | 'desc';

export default useMatchStore;