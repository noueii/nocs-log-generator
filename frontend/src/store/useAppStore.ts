/**
 * Main application store using Zustand
 * Handles match configuration, generation status, and UI state
 */

import { create } from 'zustand';
import { devtools, subscribeWithSelector, persist } from 'zustand/middleware';
import { immer } from 'zustand/middleware/immer';
import { matchService } from '@/services/matchService';
import type { 
  IMatchConfig, 
  IGenerateRequest, 
  IGenerateResponse,
  ITeamFormData,
  TPlayerRole,
  TSide
} from '@/types';
import { DEFAULT_MATCH_CONFIG } from '@/types';

/**
 * Generation status type
 */
export type TGenerationStatus = 
  | 'idle' 
  | 'validating' 
  | 'generating' 
  | 'completed' 
  | 'error';

/**
 * Configuration template type
 */
export interface IConfigTemplate {
  id: string;
  name: string;
  description: string;
  config: IMatchConfig;
}

/**
 * Application store state
 */
export interface IAppStore {
  // Match Configuration
  matchConfig: IMatchConfig;
  teams: ITeamFormData[];
  
  // Generation State
  generationStatus: TGenerationStatus;
  generationProgress: number;
  currentGeneration: IGenerateResponse | null;
  generationError: string | null;
  
  // UI State
  activeTab: string;
  isLoading: boolean;
  toastMessage: string | null;
  toastType: 'success' | 'error' | 'info' | null;
  
  // Templates and Settings
  configTemplates: IConfigTemplate[];
  isTemplatesLoading: boolean;
  
  // Actions - Configuration
  setMatchConfig: (config: Partial<IMatchConfig>) => void;
  resetMatchConfig: () => void;
  loadConfigTemplate: (templateId: string) => void;
  
  // Actions - Teams
  setTeams: (teams: ITeamFormData[]) => void;
  updateTeam: (teamIndex: number, team: Partial<ITeamFormData>) => void;
  resetTeams: () => void;
  
  // Actions - Generation
  generateMatch: () => Promise<void>;
  cancelGeneration: () => void;
  clearGenerationResult: () => void;
  retryGeneration: () => Promise<void>;
  
  // Actions - UI
  setActiveTab: (tab: string) => void;
  setLoading: (loading: boolean) => void;
  showToast: (message: string, type: 'success' | 'error' | 'info') => void;
  hideToast: () => void;
  
  // Actions - Templates
  loadConfigTemplates: () => Promise<void>;
  
  // Computed properties
  canProceedToSettings: () => boolean;
  canGenerate: () => boolean;
  isGenerating: () => boolean;
}

/**
 * Default team structure
 */
const createDefaultTeams = (): ITeamFormData[] => [
  {
    name: "",
    tag: "",
    country: "US",
    players: Array(5).fill(null).map((_, i) => ({
      name: "",
      role: (i === 0 ? "entry" : i === 1 ? "awp" : i === 2 ? "support" : i === 3 ? "lurker" : "igl") as TPlayerRole,
      rating: 1.0,
      steam_id: "",
      country: "US"
    }))
  },
  {
    name: "",
    tag: "",
    country: "US",
    players: Array(5).fill(null).map((_, i) => ({
      name: "",
      role: (i === 0 ? "entry" : i === 1 ? "awp" : i === 2 ? "support" : i === 3 ? "lurker" : "igl") as const,
      rating: 1.0,
      steam_id: "",
      country: "US"
    }))
  }
];

/**
 * Default configuration templates
 */
const getDefaultTemplates = (): IConfigTemplate[] => [
  {
    id: 'competitive',
    name: 'Competitive Match',
    description: 'Standard competitive match (MR12)',
    config: {
      ...DEFAULT_MATCH_CONFIG,
      format: 'mr12',
      realistic_economy: true,
      skill_variance: 0.15,
    }
  },
  {
    id: 'casual',
    name: 'Casual Match', 
    description: 'Relaxed match settings (MR15)',
    config: {
      ...DEFAULT_MATCH_CONFIG,
      format: 'mr15',
      realistic_economy: false,
      skill_variance: 0.25,
      chat_messages: true,
    }
  },
  {
    id: 'testing',
    name: 'Testing Match',
    description: 'Verbose logging for testing',
    config: {
      ...DEFAULT_MATCH_CONFIG,
      output_verbosity: 'verbose',
      include_positions: true,
      include_weapon_fire: true,
    }
  }
];

/**
 * Create the app store
 */
export const useAppStore = create<IAppStore>()(
  devtools(
    persist(
      subscribeWithSelector(
        immer((set, get) => ({
        // Initial State
        matchConfig: DEFAULT_MATCH_CONFIG,
        teams: createDefaultTeams(),
        
        generationStatus: 'idle',
        generationProgress: 0,
        currentGeneration: null,
        generationError: null,
        
        activeTab: 'teams',
        isLoading: false,
        toastMessage: null,
        toastType: null,
        
        configTemplates: getDefaultTemplates(),
        isTemplatesLoading: false,
        
        // Configuration Actions
        setMatchConfig: (config) => {
          set((state) => {
            Object.assign(state.matchConfig, config);
          });
        },
        
        resetMatchConfig: () => {
          set((state) => {
            state.matchConfig = { ...DEFAULT_MATCH_CONFIG };
          });
        },
        
        loadConfigTemplate: (templateId) => {
          const template = get().configTemplates.find(t => t.id === templateId);
          if (template) {
            set((state) => {
              state.matchConfig = { ...template.config };
            });
            get().showToast(`Loaded template: ${template.name}`, 'success');
          }
        },
        
        // Teams Actions
        setTeams: (teams) => {
          set((state) => {
            state.teams = teams;
          });
        },
        
        updateTeam: (teamIndex, teamUpdate) => {
          set((state) => {
            if (state.teams[teamIndex]) {
              Object.assign(state.teams[teamIndex], teamUpdate);
            }
          });
        },
        
        resetTeams: () => {
          set((state) => {
            state.teams = createDefaultTeams();
          });
        },
        
        // Generation Actions
        generateMatch: async () => {
          const state = get();
          
          if (!state.canGenerate()) {
            state.showToast('Please complete team setup and configuration', 'error');
            return;
          }
          
          set((state) => {
            state.generationStatus = 'validating';
            state.generationProgress = 0;
            state.generationError = null;
            state.currentGeneration = null;
          });
          
          try {
            // Build generation request
            const generateRequest: IGenerateRequest = {
              teams: state.teams.map((teamData, index) => ({
                name: teamData.name,
                tag: teamData.tag,
                country: teamData.country,
                side: (index === 0 ? "CT" : "TERRORIST") as TSide,
                score: 0,
                rounds_won: 0,
                economy: {
                  total_money: state.matchConfig.start_money * 5,
                  average_money: state.matchConfig.start_money,
                  equipment_value: 0,
                  consecutive_losses: 0,
                  loss_bonus: 1400,
                  money_spent: 0,
                  money_earned: 0,
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
                },
                players: teamData.players.map(playerData => ({
                  name: playerData.name,
                  steam_id: playerData.steam_id || `STEAM_1:0:${Math.floor(Math.random() * 1000000)}`,
                  role: playerData.role,
                  team: teamData.name,
                  side: (index === 0 ? "CT" : "TERRORIST") as const,
                  state: {
                    is_alive: true,
                    health: 100,
                    armor: 0,
                    has_helmet: false,
                    has_defuse_kit: false,
                    position: { x: 0, y: 0, z: 0 },
                    view_angle: { x: 0, y: 0, z: 0 },
                    velocity: { x: 0, y: 0, z: 0 },
                    grenades: [],
                    money: state.matchConfig.start_money,
                    is_flashed: false,
                    is_smoked: false,
                    is_defusing: false,
                    is_planting: false,
                    is_reloading: false,
                    has_bomb: false,
                    is_last_alive: false
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
                    total_damage: 0,
                    headshot_kills: 0,
                    utility_damage: 0,
                    enemies_flashed: 0,
                    damage: 0,
                    headshots: 0,
                    headshot_rate: 0,
                    accuracy: 0,
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
                    money_spent: 0,
                    grenades_thrown: {},
                    flash_assists: 0,
                    team_kills: 0,
                    team_damage: 0,
                    kd_ratio: 0,
                    rating: playerData.rating,
                    kast: 0
                  },
                  economy: {
                    money: state.matchConfig.start_money,
                    money_spent: 0,
                    money_earned: state.matchConfig.start_money,
                    equipment_value: 0,
                    purchases: [],
                    eco_rounds: 0,
                    force_buy_rounds: 0,
                    full_buy_rounds: 0,
                    economy_rating: 0
                  }
                }))
              })),
              map: state.matchConfig.map,
              format: state.matchConfig.format,
              options: {
                seed: state.matchConfig.seed,
                tick_rate: state.matchConfig.tick_rate,
                overtime: state.matchConfig.overtime,
                max_rounds: state.matchConfig.max_rounds
              }
            };
            
            set((state) => {
              state.generationStatus = 'generating';
              state.generationProgress = 25;
            });
            
            // Call the API
            const result = await matchService.generateMatch(generateRequest);
            
            set((state) => {
              state.generationStatus = 'completed';
              state.generationProgress = 100;
              state.currentGeneration = result;
              state.activeTab = 'results';
            });
            
            get().showToast('Match generated successfully!', 'success');
            
          } catch (error) {
            console.error('Match generation failed:', error);
            
            const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred';
            
            set((state) => {
              state.generationStatus = 'error';
              state.generationError = errorMessage;
              state.currentGeneration = {
                match_id: "",
                status: "error",
                error: errorMessage
              };
              state.activeTab = 'results';
            });
            
            get().showToast(`Generation failed: ${errorMessage}`, 'error');
          }
        },
        
        cancelGeneration: () => {
          set((state) => {
            state.generationStatus = 'idle';
            state.generationProgress = 0;
            state.generationError = null;
          });
        },
        
        clearGenerationResult: () => {
          set((state) => {
            state.generationStatus = 'idle';
            state.generationProgress = 0;
            state.currentGeneration = null;
            state.generationError = null;
          });
        },
        
        retryGeneration: async () => {
          get().clearGenerationResult();
          await get().generateMatch();
        },
        
        // UI Actions
        setActiveTab: (tab) => {
          set((state) => {
            state.activeTab = tab;
          });
        },
        
        setLoading: (loading) => {
          set((state) => {
            state.isLoading = loading;
          });
        },
        
        showToast: (message, type) => {
          set((state) => {
            state.toastMessage = message;
            state.toastType = type;
          });
          
          // Auto-hide toast after 5 seconds
          setTimeout(() => {
            get().hideToast();
          }, 5000);
        },
        
        hideToast: () => {
          set((state) => {
            state.toastMessage = null;
            state.toastType = null;
          });
        },
        
        // Templates Actions
        loadConfigTemplates: async () => {
          set((state) => {
            state.isTemplatesLoading = true;
          });
          
          try {
            const templates = await matchService.getMatchTemplates();
            const configTemplates = templates.map(t => ({
              id: t.id,
              name: t.name,
              description: t.description,
              config: t.config
            }));
            
            set((state) => {
              state.configTemplates = configTemplates;
              state.isTemplatesLoading = false;
            });
          } catch (error) {
            console.error('Failed to load templates:', error);
            
            // Keep default templates on error
            set((state) => {
              state.isTemplatesLoading = false;
            });
          }
        },
        
        // Computed Properties
        canProceedToSettings: () => {
          const { teams } = get();
          return teams.every(team => 
            team.name.trim() !== "" && 
            team.players.every(player => player.name.trim() !== "")
          );
        },
        
        canGenerate: () => {
          const state = get();
          return state.canProceedToSettings() && !!state.matchConfig;
        },
        
        isGenerating: () => {
          const { generationStatus } = get();
          return generationStatus === 'validating' || generationStatus === 'generating';
        },
      }))
      ),
      {
        name: 'app-store-persist',
        // Only persist user preferences and configuration, not UI state
        partialize: (state) => ({
          matchConfig: state.matchConfig,
          teams: state.teams,
          configTemplates: state.configTemplates,
        }),
      }
    ),
    {
      name: 'app-store',
    }
  )
);

/**
 * Hook for reactive subscriptions to specific state slices
 */
export const useAppStoreSelector = <T>(selector: (state: IAppStore) => T) => 
  useAppStore(selector);

/**
 * Hook for generation status
 */
export const useGenerationStatus = () => useAppStore((state) => ({
  status: state.generationStatus,
  progress: state.generationProgress,
  result: state.currentGeneration,
  error: state.generationError,
  isGenerating: state.isGenerating()
}));

/**
 * Hook for match configuration
 */
export const useMatchConfig = () => useAppStore((state) => ({
  config: state.matchConfig,
  setConfig: state.setMatchConfig,
  resetConfig: state.resetMatchConfig,
  templates: state.configTemplates,
  loadTemplate: state.loadConfigTemplate
}));

/**
 * Hook for teams management
 */
export const useTeamsManagement = () => useAppStore((state) => ({
  teams: state.teams,
  setTeams: state.setTeams,
  updateTeam: state.updateTeam,
  resetTeams: state.resetTeams
}));

export default useAppStore;