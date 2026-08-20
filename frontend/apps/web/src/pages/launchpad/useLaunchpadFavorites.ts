import {
  addLaunchpadFavoritesApi,
  getLaunchpadFavoritesApi,
  removeLaunchpadFavoritesApi,
} from '@/apis/core-api/launchpad';
import { useCallback, useEffect, useState } from 'react';

/**
 * Launchpad 收藏状态。
 * 应用收藏优先同步后端；后端不可用或请求失败时自动降级为浏览器本地存储。
 * 项目收藏暂由本地存储维护。
 */
const STORAGE_KEY = 'launchpad:favorites';

interface FavoriteState {
  apps: string[];
  projects: string[];
}

const readFavorites = (): FavoriteState => {
  try {
    const raw = window.localStorage.getItem(STORAGE_KEY);
    if (!raw) return { apps: [], projects: [] };
    const parsed = JSON.parse(raw);
    return {
      apps: Array.isArray(parsed?.apps) ? parsed.apps.map(String) : [],
      projects: Array.isArray(parsed?.projects) ? parsed.projects.map(String) : [],
    };
  } catch {
    return { apps: [], projects: [] };
  }
};

const writeFavorites = (state: FavoriteState) => {
  try {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
  } catch {
    /* localStorage 不可用时静默降级为内存态 */
  }
};

const toggleId = (list: string[], id: string) =>
  list.includes(id) ? list.filter((item) => item !== id) : [...list, id];

export interface LaunchpadFavorites {
  isAppFavorite: (id: string | number) => boolean;
  isProjectFavorite: (id: string | number) => boolean;
  favoriteAppIds: string[];
  toggleAppFavorite: (id: string | number) => void;
  toggleProjectFavorite: (id: string | number) => void;
}

export function useLaunchpadFavorites(): LaunchpadFavorites {
  const [state, setState] = useState<FavoriteState>(readFavorites);
  const [backendEnabled, setBackendEnabled] = useState(false);

  // 初始化时拉取后端收藏；后端不可用时保留本地存储作为降级。
  useEffect(() => {
    let mounted = true;
    void getLaunchpadFavoritesApi('app')
      .then((res) => {
        if (!mounted) return;
        const appIds = (res?.list || []).map((item) => String(item.resId));
        setState((prev) => {
          const next = { ...prev, apps: appIds };
          writeFavorites(next);
          return next;
        });
        setBackendEnabled(true);
      })
      .catch(() => {
        if (!mounted) return;
        setBackendEnabled(false);
      });
    return () => {
      mounted = false;
    };
  }, []);

  const isAppFavorite = useCallback((id: string | number) => state.apps.includes(String(id)), [state.apps]);
  const isProjectFavorite = useCallback(
    (id: string | number) => state.projects.includes(String(id)),
    [state.projects]
  );

  const toggleAppFavorite = useCallback(
    async (id: string | number) => {
      const idStr = String(id);
      let previousApps: string[] = [];
      let isAdding = false;

      setState((prev) => {
        previousApps = prev.apps;
        isAdding = !previousApps.includes(idStr);
        const nextApps = isAdding ? [...previousApps, idStr] : previousApps.filter((item) => item !== idStr);
        const next = { ...prev, apps: nextApps };
        writeFavorites(next);
        return next;
      });

      if (!backendEnabled) return;

      try {
        const item = { resType: 'app', resId: idStr };
        if (isAdding) {
          await addLaunchpadFavoritesApi([item]);
        } else {
          await removeLaunchpadFavoritesApi([item]);
        }
      } catch {
        // 后端同步失败时回滚本地状态
        setState((prev) => {
          const next = { ...prev, apps: previousApps };
          writeFavorites(next);
          return next;
        });
      }
    },
    [backendEnabled]
  );

  const toggleProjectFavorite = useCallback((id: string | number) => {
    setState((prev) => {
      const next = { ...prev, projects: toggleId(prev.projects, String(id)) };
      writeFavorites(next);
      return next;
    });
  }, []);

  return {
    isAppFavorite,
    isProjectFavorite,
    favoriteAppIds: state.apps,
    toggleAppFavorite,
    toggleProjectFavorite,
  };
}
