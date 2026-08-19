export type FlowTabType = 'all' | 'source' | 'event';

export type FlowKind = 'source' | 'event';

export interface FlowListPanelHandle {
  refreshRequest: () => void;
  getCreateContext: () => { groupId?: number };
}
