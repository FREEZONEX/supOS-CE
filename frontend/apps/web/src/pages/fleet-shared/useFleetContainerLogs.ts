import {
  useCallback,
  useEffect,
  useReducer,
  useRef,
  useState,
  type Dispatch,
  type SetStateAction,
  type UIEvent,
} from 'react';
import type { FleetContainerLogsReq, FleetContainerLogsResp, FleetContainerSummary } from '@/apis/core-api/fleet';

const POLL_INTERVAL_MS = 5000;
const MAX_LOG_LINES = 2000;
const MAX_LOG_TABS = 5;
const DOWNLOAD_LOG_LIMIT_BYTES = 2 * 1024 * 1024;

const mergeLogText = (current: string, incoming: string) => {
  const previousLines = current ? current.split('\n') : [];
  const incomingLines = incoming ? incoming.split('\n') : [];
  let overlap = Math.min(previousLines.length, incomingLines.length);
  while (overlap > 0) {
    if (previousLines.slice(-overlap).join('\n') === incomingLines.slice(0, overlap).join('\n')) break;
    overlap -= 1;
  }
  return [...previousLines, ...incomingLines.slice(overlap)].slice(-MAX_LOG_LINES).join('\n');
};

type FleetContainerLogsFetcher = (
  container: FleetContainerSummary,
  request: FleetContainerLogsReq
) => Promise<FleetContainerLogsResp>;

type UseFleetContainerLogsOptions = {
  fetchLogs: FleetContainerLogsFetcher;
  onError: () => void;
};

export type FleetContainerLogTab = {
  id: string;
  container: FleetContainerSummary;
  logs: string;
  loading: boolean;
  live: boolean;
  follow: boolean;
  filterText: string;
  observedAt: number;
};

type LogsState = {
  visible: boolean;
  activeID?: string;
  tabs: FleetContainerLogTab[];
};

type LogsAction =
  | { type: 'open'; container: FleetContainerSummary }
  | { type: 'hide' }
  | { type: 'activate'; id: string }
  | { type: 'close'; id: string }
  | { type: 'reorder'; fromID: string; toID: string }
  | { type: 'patch'; id: string; patch: Partial<FleetContainerLogTab> }
  | { type: 'loadSucceeded'; id: string; full: boolean; result: FleetContainerLogsResp }
  | { type: 'loadFailed'; id: string; full: boolean };

const updateTab = (
  tabs: FleetContainerLogTab[],
  id: string,
  update: (tab: FleetContainerLogTab) => FleetContainerLogTab
) => tabs.map((tab) => (tab.id === id ? update(tab) : tab));

const logsReducer = (state: LogsState, action: LogsAction): LogsState => {
  switch (action.type) {
    case 'open': {
      const id = action.container.containerID;
      const exists = state.tabs.some((tab) => tab.id === id);
      let tabs = state.tabs;
      if (!exists && tabs.length >= MAX_LOG_TABS) {
        const removableIndex = tabs.findIndex((tab) => tab.id !== state.activeID);
        tabs = tabs.filter((_, index) => index !== (removableIndex >= 0 ? removableIndex : 0));
      }
      return {
        visible: true,
        activeID: id,
        tabs: exists
          ? updateTab(tabs, id, (tab) => ({ ...tab, container: action.container }))
          : [
              ...tabs,
              {
                id,
                container: action.container,
                logs: '',
                loading: false,
                live: true,
                follow: true,
                filterText: '',
                observedAt: 0,
              },
            ],
      };
    }
    case 'hide':
      return { ...state, visible: false };
    case 'activate':
      return state.tabs.some((tab) => tab.id === action.id) ? { ...state, activeID: action.id } : state;
    case 'close': {
      const index = state.tabs.findIndex((tab) => tab.id === action.id);
      if (index < 0) return state;
      const tabs = state.tabs.filter((tab) => tab.id !== action.id);
      const activeID =
        state.activeID === action.id ? tabs[Math.min(index, Math.max(0, tabs.length - 1))]?.id : state.activeID;
      return { visible: tabs.length > 0 && state.visible, activeID, tabs };
    }
    case 'reorder': {
      const from = state.tabs.findIndex((tab) => tab.id === action.fromID);
      const to = state.tabs.findIndex((tab) => tab.id === action.toID);
      if (from < 0 || to < 0 || from === to) return state;
      const tabs = [...state.tabs];
      const [moved] = tabs.splice(from, 1);
      tabs.splice(to, 0, moved);
      return { ...state, tabs };
    }
    case 'patch':
      return { ...state, tabs: updateTab(state.tabs, action.id, (tab) => ({ ...tab, ...action.patch })) };
    case 'loadSucceeded':
      return {
        ...state,
        tabs: updateTab(state.tabs, action.id, (tab) => ({
          ...tab,
          loading: false,
          logs: action.full ? action.result.logs || '' : mergeLogText(tab.logs, action.result.logs || ''),
          observedAt: action.result.observedAt || Date.now(),
        })),
      };
    case 'loadFailed':
      return {
        ...state,
        tabs: updateTab(state.tabs, action.id, (tab) => ({
          ...tab,
          loading: false,
          logs: action.full ? '' : tab.logs,
        })),
      };
    default:
      return state;
  }
};

export const useFleetContainerLogs = ({ fetchLogs, onError }: UseFleetContainerLogsOptions) => {
  const [state, dispatch] = useReducer(logsReducer, { visible: false, tabs: [] });
  const [downloadingID, setDownloadingID] = useState<string>();
  const requestingRef = useRef(new Set<string>());
  const tabsRef = useRef(state.tabs);
  const logsRef = useRef<HTMLPreElement>(null);

  const activeTab = state.tabs.find((tab) => tab.id === state.activeID);

  useEffect(() => {
    tabsRef.current = state.tabs;
  }, [state.tabs]);

  const load = useCallback(
    async (target: FleetContainerSummary, full: boolean, observedAt = 0, silent = false) => {
      const id = target.containerID;
      if (requestingRef.current.has(id)) return;
      requestingRef.current.add(id);
      if (!silent) dispatch({ type: 'patch', id, patch: { loading: true } });
      try {
        const result = await fetchLogs(target, {
          tail: full ? 500 : 2000,
          since: full || observedAt <= 0 ? undefined : Math.max(1, observedAt - 1000),
          timestamps: true,
          limitBytes: full ? 1024 * 1024 : 512 * 1024,
        });
        dispatch({ type: 'loadSucceeded', id, full, result });
      } catch {
        dispatch({ type: 'loadFailed', id, full });
        if (!silent) onError();
      } finally {
        requestingRef.current.delete(id);
      }
    },
    [fetchLogs, onError]
  );

  const openLogs = useCallback(
    async (target: FleetContainerSummary) => {
      const existing = tabsRef.current.find((tab) => tab.id === target.containerID);
      dispatch({ type: 'open', container: target });
      if (!existing) {
        await load(target, true);
      } else if (existing.live) {
        await load(target, false, existing.observedAt, true);
      }
    },
    [load]
  );

  const closeLogs = useCallback(() => dispatch({ type: 'hide' }), []);
  const closeLogTab = useCallback((id: string) => dispatch({ type: 'close', id }), []);
  const reorderTabs = useCallback((fromID: string, toID: string) => dispatch({ type: 'reorder', fromID, toID }), []);

  const setActiveTab = useCallback(
    (id: string) => {
      const tab = tabsRef.current.find((item) => item.id === id);
      if (!tab) return;
      dispatch({ type: 'activate', id });
      if (tab.live) void load(tab.container, false, tab.observedAt, true);
    },
    [load]
  );

  const refreshLogs = useCallback(async () => {
    if (activeTab) await load(activeTab.container, true);
  }, [activeTab, load]);

  const downloadLogs = useCallback(async () => {
    if (!activeTab) return undefined;
    const target = activeTab.container;
    setDownloadingID(activeTab.id);
    try {
      return await fetchLogs(target, {
        tail: MAX_LOG_LINES,
        timestamps: true,
        limitBytes: DOWNLOAD_LOG_LIMIT_BYTES,
      });
    } catch {
      onError();
      return undefined;
    } finally {
      setDownloadingID((current) => (current === target.containerID ? undefined : current));
    }
  }, [activeTab, fetchLogs, onError]);

  const clearLogs = useCallback(() => {
    if (!activeTab) return;
    dispatch({ type: 'patch', id: activeTab.id, patch: { logs: '', observedAt: Date.now() + 1 } });
  }, [activeTab]);

  const setLive: Dispatch<SetStateAction<boolean>> = useCallback(
    (value) => {
      if (!activeTab) return;
      const live = typeof value === 'function' ? value(activeTab.live) : value;
      dispatch({ type: 'patch', id: activeTab.id, patch: { live } });
    },
    [activeTab]
  );

  const setFollow: Dispatch<SetStateAction<boolean>> = useCallback(
    (value) => {
      if (!activeTab) return;
      const follow = typeof value === 'function' ? value(activeTab.follow) : value;
      dispatch({ type: 'patch', id: activeTab.id, patch: { follow } });
    },
    [activeTab]
  );

  const setFilterText: Dispatch<SetStateAction<string>> = useCallback(
    (value) => {
      if (!activeTab) return;
      const filterText = typeof value === 'function' ? value(activeTab.filterText) : value;
      dispatch({ type: 'patch', id: activeTab.id, patch: { filterText } });
    },
    [activeTab]
  );

  const handleScroll = useCallback(
    (event: UIEvent<HTMLPreElement>) => {
      if (!activeTab) return;
      const target = event.currentTarget;
      dispatch({
        type: 'patch',
        id: activeTab.id,
        patch: { follow: target.scrollHeight - target.scrollTop - target.clientHeight < 32 },
      });
    },
    [activeTab]
  );

  useEffect(() => {
    if (!state.visible) return;
    const timer = window.setInterval(() => {
      if (document.visibilityState !== 'visible') return;
      tabsRef.current.forEach((tab) => {
        if (tab.live) void load(tab.container, false, tab.observedAt, true);
      });
    }, POLL_INTERVAL_MS);
    return () => window.clearInterval(timer);
  }, [load, state.visible]);

  useEffect(() => {
    if (!activeTab?.follow || !logsRef.current) return;
    const frame = window.requestAnimationFrame(() => {
      if (logsRef.current) logsRef.current.scrollTop = logsRef.current.scrollHeight;
    });
    return () => window.cancelAnimationFrame(frame);
  }, [activeTab?.filterText, activeTab?.follow, activeTab?.id, activeTab?.logs]);

  return {
    open: state.visible && state.tabs.length > 0,
    tabs: state.tabs,
    activeContainerID: state.activeID,
    loading: activeTab?.loading || false,
    downloading: downloadingID === activeTab?.id,
    logs: activeTab?.logs || '',
    container: activeTab?.container,
    live: activeTab?.live ?? true,
    follow: activeTab?.follow ?? true,
    filterText: activeTab?.filterText || '',
    logsRef,
    openLogs,
    closeLogs,
    closeLogTab,
    reorderTabs,
    setActiveTab,
    refreshLogs,
    downloadLogs,
    clearLogs,
    setLive,
    setFollow,
    setFilterText,
    handleScroll,
  };
};

export type FleetContainerLogsController = ReturnType<typeof useFleetContainerLogs>;
