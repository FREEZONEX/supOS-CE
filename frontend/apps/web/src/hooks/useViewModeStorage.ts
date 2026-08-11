import { useLocalStorageState } from 'ahooks';

export const VIEW_MODE_LIST = 'list';
export const VIEW_MODE_CARD = 'card';

export const VIEW_MODE_STORAGE_KEYS = {
  collection: 'APP_COLLECTION_VIEW_V2',
  eventFlow: 'APP_EVENTFLOW_VIEW_V2',
  flowAll: 'APP_FLOW_ALL_VIEW_V2',
  flow: 'APP_FLOW_VIEW_V2',
  app: 'APP_APP_VIEW_V2',
  plugin: 'APP_PLUGIN_VIEW_V2',
  visionAlgorithm: 'APP_VISION_ALGO_VIEW_V2',
} as const;

export function useViewModeStorage(key: string) {
  return useLocalStorageState<string>(key, {
    defaultValue: VIEW_MODE_LIST,
  });
}
