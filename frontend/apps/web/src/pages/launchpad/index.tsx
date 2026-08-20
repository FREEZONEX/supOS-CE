import { getLaunchpadProjectsApi, getLaunchpadRecentsApi, recordLaunchpadAppViewApi } from '@/apis/core-api/launchpad';
import empty from '@/assets/common/empty.svg';
import { useActivate } from '@/contexts/tabs-lifecycle-context';
import { useDebouncedRequest, useTranslate } from '@/hooks';
import type { App, Project } from '@/pages/launchpad/type';
import { Run } from '@carbon/icons-react';
import { Empty, Spin, Typography } from 'antd';
import { type FC, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useNavigate } from 'react-router';
import AllProjectsSection from './components/AllProjectsSection';
import FavoriteRecentSection from './components/FavoriteRecentSection';
import { useLaunchpadFavorites } from './useLaunchpadFavorites';
import styles from './index.module.scss';

const SEARCH_DEBOUNCE_MS = 300;

const LaunchpadPage: FC = () => {
  const formatMessage = useTranslate('Launchpad');
  const navigate = useNavigate();
  const [searchValue, setSearchValue] = useState('');
  const [appliedSearchKeyword, setAppliedSearchKeyword] = useState('');
  const [viewMode, setViewMode] = useState('card');
  const [recentApps, setRecentApps] = useState<App[]>([]);
  const latestSearchValueRef = useRef('');
  const hasSearchKeyword = Boolean(appliedSearchKeyword);
  const { isAppFavorite, isProjectFavorite, toggleAppFavorite, toggleProjectFavorite } = useLaunchpadFavorites();

  const {
    loading,
    data: projects,
    runImmediate,
    runDebounced,
  } = useDebouncedRequest(
    async (keyword?: string) => {
      try {
        return await getLaunchpadProjectsApi(keyword);
      } catch (error: any) {
        if (error?.code === 403) {
          return [] as Project[];
        }
        throw error;
      }
    },
    {
      debounceMs: SEARCH_DEBOUNCE_MS,
      initialData: [] as Project[],
    }
  );

  const loadRecents = useCallback(async () => {
    try {
      const res = await getLaunchpadRecentsApi();
      setRecentApps(res?.list || []);
    } catch {
      setRecentApps([]);
    }
  }, []);

  useEffect(() => {
    void runImmediate(undefined);
    void loadRecents();
  }, [runImmediate, loadRecents]);

  useActivate(() => {
    runDebounced.cancel();
    void runImmediate(appliedSearchKeyword || undefined);
    void loadRecents();
  });

  const handleSearchChange = useCallback(
    (value: string) => {
      setSearchValue(value);
      latestSearchValueRef.current = value;
      const keyword = value.trim();
      if (!keyword) {
        runDebounced.cancel();
        void runImmediate(undefined).finally(() => {
          if (!latestSearchValueRef.current.trim()) {
            setAppliedSearchKeyword('');
          }
        });
        return;
      }
      setAppliedSearchKeyword(keyword);
      runDebounced(keyword);
    },
    [runDebounced, runImmediate]
  );

  const handleProjectClick = (project: Project) => {
    navigate(`/launchpad/${encodeURIComponent(project.projectName)}`, {
      state: { tabName: project.displayName || project.projectName },
    });
  };

  const handleAppClick = (app: App) => {
    const openInNewTab = app.manual && app.openInPlatform === false && app.url;
    if (openInNewTab) {
      void recordLaunchpadAppViewApi(app.appId);
      window.open(app.url, '_blank', 'noopener,noreferrer');
      return;
    }
    navigate(
      `/launchpad/${encodeURIComponent(app.projectName || '')}/${encodeURIComponent(String(app.appId || app.appName))}`,
      { state: { tabName: app.displayName || app.appName } }
    );
  };

  const isInitialEmpty = !loading && projects.length === 0 && !hasSearchKeyword;

  // 对齐云端：App 数量多的项目排前面；数量相同时按名称升序，保证顺序稳定
  const sortedProjects = useMemo(
    () =>
      [...projects].sort((a, b) => {
        const countDiff = (b.apps?.length || 0) - (a.apps?.length || 0);
        if (countDiff !== 0) return countDiff;
        const nameA = a.displayName || a.projectName;
        const nameB = b.displayName || b.projectName;
        return nameA.localeCompare(nameB);
      }),
    [projects]
  );

  return (
    <div className={styles.launchpadPage}>
      <div className={styles.pageHeader}>
        <div className={styles.pageTitle}>
          <Run size={20} />
          <h1>{formatMessage('title')}</h1>
        </div>
      </div>
      <div className={styles.mainContent}>
        {loading ? (
          <div className={styles.loading}>
            <Spin size="large" />
          </div>
        ) : isInitialEmpty ? (
          <div className={styles.emptyState}>
            <Empty
              image={<img src={empty} alt="" />}
              description={
                <div className={styles['empty-desc']}>
                  <Typography.Title level={3} className={styles['empty-title']}>
                    {formatMessage('noProjectsTitle')}
                  </Typography.Title>
                  <Typography.Text type="secondary" className={styles['empty-text']}>
                    {formatMessage('noProjectsDescription')}
                  </Typography.Text>
                </div>
              }
            />
          </div>
        ) : (
          <>
            {!hasSearchKeyword && (
              <FavoriteRecentSection
                apps={recentApps}
                isAppFavorite={isAppFavorite}
                onToggleFavorite={(app) => toggleAppFavorite(app.appId)}
                onAppClick={handleAppClick}
              />
            )}
            <AllProjectsSection
              projects={sortedProjects}
              searchValue={searchValue}
              onSearchChange={handleSearchChange}
              appliedSearchKeyword={appliedSearchKeyword}
              hasSearchKeyword={hasSearchKeyword}
              viewMode={viewMode}
              onViewModeChange={setViewMode}
              isProjectFavorite={isProjectFavorite}
              onToggleFavorite={(project) => toggleProjectFavorite(project.projectId)}
              onProjectClick={handleProjectClick}
            />
          </>
        )}
      </div>
    </div>
  );
};

export default LaunchpadPage;
