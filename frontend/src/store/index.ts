// State management
// File naming: camelCase.ts  
// Export all stores from this index file

// Main application store
export { 
  useAppStore,
  useAppStoreSelector,
  useGenerationStatus,
  useMatchConfig,
  useTeamsManagement,
  type IAppStore,
  type TGenerationStatus,
  type IConfigTemplate
} from './useAppStore';

// Match history store
export {
  useMatchStore,
  useMatchHistory,
  useMatchFilters, 
  useMatchSelection,
  type IMatchStore,
  type IMatchHistoryItem,
  type IMatchFilters,
  type TMatchSortBy,
  type TMatchSortOrder
} from './useMatchStore';

// Main stores (already exported above as named exports)
// useAppStore and useMatchStore are available as named exports