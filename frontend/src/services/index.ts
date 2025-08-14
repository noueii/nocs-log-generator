/**
 * API clients and external service integrations
 * 
 * File naming conventions:
 * - Files: camelCase.ts (api.ts, matchService.ts, etc.)
 * - Classes: PascalCase (MatchService, ApiClient, etc.)
 * - Functions: camelCase (generateMatch, getMatches, etc.)
 */

// Main API client
export {
  apiClient as default,
  apiClient,
  checkApiHealth,
  getApiStatus,
  handleApiCall,
  createCustomApiClient,
  withTimeout,
  get,
  post,
  put,
  del,
} from './api';

export type {
  IApiResponse,
  IApiError,
} from './api';

// Match service
export {
  matchService,
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
} from './matchService';

export type {
  IMatchListResponse,
  IMatchSearchParams,
  IMatchTemplate,
  IMapInfo,
  IServerStatus,
  IGenerationProgress,
} from './matchService';