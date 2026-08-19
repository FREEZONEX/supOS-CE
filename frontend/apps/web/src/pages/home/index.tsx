import { type ReactNode, createElement, useCallback, useEffect, useRef, useState } from 'react';
import ComEmpty from '@/components/com-empty';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import styles from './index.module.scss';
import { useActivate } from '@/contexts/tabs-lifecycle-context.ts';
import ComLayout from '@/components/com-layout';
import ComContent from '@/components/com-layout/ComContent';
import { fetchBaseStore, useBaseStore } from '@/stores/base';
import { useTranslate } from '@/hooks';
import { getHomeRecents, type HomeRecentItem } from '@/apis/core-api';
import useMenuNavigate from '@/hooks/useMenuNavigate';
import type { ResourceProps } from '@/stores/types';
import { resolveMenuLucideIcon } from '@/components/lucide-icon';
import { useNavigate } from 'react-router';

dayjs.extend(relativeTime);

type RecentWorkItem = {
  key: string;
  type: string;
  title: string;
  time: string;
  icon: ReactNode;
  item?: ResourceProps;
  path?: string;
  disabled?: boolean;
};


const RECENT_TYPE_KEY_BY_RESOURCE_TYPE: Record<string, string> = {
  uns: 'home.recentTypeUns',
  namespace: 'home.recentTypeUns',
  sourceflow: 'home.recentTypeFlow',
  eventflow: 'home.recentTypeFlow',
  flow: 'home.recentTypeFlow',
  notebook: 'Notebook.title',
  project: 'Project',
  projectapp: 'project.apps',
  app: 'project.apps',
  launchpad: 'Launchpad.title',
  oauthclient: 'menu.oauthClient',
  apikey: 'menu.apiKey',
  edgeconnection: 'menu.mqttAuth',
  mqttauth: 'menu.mqttAuth',
  routingmanagement: 'route.routingManagement',
  gateway: 'route.routingManagement',
  cluster: 'home.actionCluster',
};

const normalizeRecentType = (value?: string) =>
  String(value || '')
    .replace(/[\s_-]+/g, '')
    .toLowerCase();

const getRecentTypeKey = (item: HomeRecentItem) => {
  const resourceKey = String(item.resourceKey || '').trim();
  if (resourceKey.startsWith('uns.')) return 'home.recentTypeUns';
  if (resourceKey.startsWith('flow.')) return 'home.recentTypeFlow';
  if (resourceKey.startsWith('notebook.')) return 'Notebook.title';
  if (resourceKey.startsWith('project.')) return 'Project';
  if (resourceKey.startsWith('launchpad.')) {
    const resourceType = normalizeRecentType(item.resourceType);
    return resourceType === 'app' || resourceType === 'projectapp' ? 'project.apps' : 'Launchpad.title';
  }
  if (resourceKey.startsWith('oauth.client.')) return 'menu.oauthClient';
  if (resourceKey.startsWith('apikey.')) return 'menu.apiKey';
  if (resourceKey.startsWith('mqtt.auth.')) return 'menu.mqttAuth';
  if (resourceKey.startsWith('gateway.route.')) return 'route.routingManagement';
  if (resourceKey.startsWith('cluster.')) return 'home.actionCluster';
  return RECENT_TYPE_KEY_BY_RESOURCE_TYPE[normalizeRecentType(item.resourceType)];
};


const formatTime = (value?: number) => (value ? dayjs(value).fromNow() : '-');


const SectionHeader = ({ title, description }: { title: string; description?: string }) => (
  <div className={styles['section-header']}>
    <h2>{title}</h2>
    {description ? <p>{description}</p> : null}
  </div>
);

const RecentWorkRow = ({ item, onNavigate }: { item: RecentWorkItem; onNavigate: (item?: ResourceProps) => void }) => {
  const navigate = useNavigate();
  const canNavigate = !item.disabled && Boolean(item.item || item.path);
  const openItem = () => {
    if (!canNavigate) {
      return;
    }
    if (item.path) {
      navigate(item.path, {
        state: {
          tabName: item.title,
        },
      });
      return;
    }
    onNavigate(item.item);
  };

  return (
    <button type="button" className={styles['recent-row']} onClick={openItem} disabled={!canNavigate}>
      <span className={styles['recent-icon']}>{item.icon}</span>
      <span className={styles['recent-name']}>
        <strong>{item.title}</strong>
      </span>
      <span className={styles['recent-type']}>{item.type}</span>
      <span className={styles['recent-time']}>{item.time}</span>
    </button>
  );
};

const Index = () => {
  const formatMessage = useTranslate();
  const navigateMenu = useMenuNavigate();
  const appTitle = useBaseStore((state) => state.systemInfo?.appTitle || 'Tier0 Edge');

  const [initialLoading, setInitialLoading] = useState(false);
  const [homeRecents, setHomeRecents] = useState<HomeRecentItem[]>([]);
  const hasLoadedHomeRef = useRef(false);

  const loadHome = useCallback(
    async (options: { silent?: boolean } = {}) => {
      const showLoading = !options.silent && !hasLoadedHomeRef.current;
      if (showLoading) {
        setInitialLoading(true);
      }
      try {
        setHomeRecents(await getHomeRecents());
        hasLoadedHomeRef.current = true;
      } finally {
        if (showLoading) {
          setInitialLoading(false);
        }
      }
    },
    []
  );

  useEffect(() => {
    fetchBaseStore?.();
  }, []);

  useEffect(() => {
    void loadHome();
  }, [loadHome]);


  useActivate(() => {
    fetchBaseStore?.();
    void loadHome({ silent: true });
  });

  const recentItems: RecentWorkItem[] = homeRecents.map((item) => {
    const typeKey = getRecentTypeKey(item);
    const actionType = normalizeRecentType(item.businessType);
    const disabled = actionType === 'delete';
    return {
      key: item.key,
      type: typeKey ? formatMessage(typeKey, undefined, item.resourceType) : item.resourceType,
      title: item.resourceName,
      time: formatTime(item.operatedAt),
      icon: createElement(
        resolveMenuLucideIcon({
          resourceKey: item.resourceKey,
          url: item.routePath,
          icon: item.icon,
        } as ResourceProps),
        { size: 12, strokeWidth: 1.75 }
      ),
      path: disabled ? undefined : item.routePath,
      disabled,
    };
  });

  return (
    <ComLayout className={styles['home-page']} loading={initialLoading}>
      <ComContent title={<div />} hasBack={false} mustShowTitle={false}>
        <div className={styles['home-shell']}>
          <header className={styles['welcome-block']}>
            <h1 className={styles['welcome-title']}>
              {formatMessage('common.welcome', { appTitle }, `Welcome to ${appTitle}!`)}
            </h1>
          </header>

          <div className={styles['home-content']}>

            <section className={styles['recent-panel']}>
              <SectionHeader
                title={formatMessage('home.recentWork')}
                description={formatMessage('home.recentWorkDesc')}
              />
              <div className={styles['recent-list']}>
                {recentItems.length ? (
                  recentItems.map((item) => <RecentWorkRow key={item.key} item={item} onNavigate={navigateMenu} />)
                ) : (
                  <div className={styles['empty-state']}>
                    <ComEmpty description={formatMessage('home.noRecentWork')} />
                  </div>
                )}
              </div>
            </section>
          </div>
        </div>
      </ComContent>
    </ComLayout>
  );
};

export default Index;
