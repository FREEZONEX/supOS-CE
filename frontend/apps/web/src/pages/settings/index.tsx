import { useMemo, type ReactNode } from 'react';
import { Api, CloudServices, User } from '@carbon/icons-react';
import AccountManagement from '@/pages/account-management';
import OpenData from '@/pages/open-data';
import AIProviderConfigPanel from '@/pages/open-data/AIProviderConfigPanel';
import { useTranslate } from '@/hooks';
import { useBaseStore } from '@/stores/base';
import { isLaunchpadStandalonePort } from '@/utils/launchpad-site';
import { useLocation, useNavigate } from 'react-router';
import AccountSettingsPanel from './AccountSettingsPanel';
import styles from './index.module.scss';

type SettingsItem = {
  key: string;
  labelKey: string;
  path: string;
  resourceKey?: string;
  accessUris?: string[];
  icon?: ReactNode;
  render: (title: string) => ReactNode;
};

const hiddenSettingsItemKeys = new Set(['permission', 'routing', 'menu']);

const canAccess = (item: SettingsItem, pageUris: string[] = [], superAdmin?: boolean, authEnable?: boolean) => {
  if (!item.resourceKey || superAdmin || authEnable === false) return true;
  const allowed = new Set([item.resourceKey, item.path, ...(item.accessUris || [])].map((uri) => uri.toLowerCase()));
  return pageUris.some((uri) => allowed.has(uri.toLowerCase()));
};

const SettingsPage = () => {
  const formatMessage = useTranslate();
  const navigate = useNavigate();
  const location = useLocation();
  const { currentUserInfo, systemInfo } = useBaseStore((state) => ({
    currentUserInfo: state.currentUserInfo,
    systemInfo: state.systemInfo,
  }));
  const pageUris = currentUserInfo?.pageList?.map((item) => (item as any).uri || item.url || '') || [];
  const groups = useMemo(
    () => [
      {
        key: 'user',
        titleKey: 'settings.user',
        icon: <User size={16} />,
        items: [
          {
            key: 'profile',
            labelKey: 'account.profile',
            path: '/settings/profile',
            render: () => <AccountSettingsPanel activeTab="profile" />,
          },
          {
            key: 'preferences',
            labelKey: 'settings.preferences',
            path: '/settings/preferences',
            render: () => <AccountSettingsPanel activeTab="preferences" />,
          },
          {
            key: 'security',
            labelKey: 'settings.security',
            path: '/settings/security',
            render: () => <AccountSettingsPanel activeTab="security" />,
          },
        ] as SettingsItem[],
      },
      {
        key: 'service',
        titleKey: 'settings.service',
        icon: <Api size={16} />,
        items: [
          {
            key: 'license',
            labelKey: 'common.license',
            path: '/settings/license',
            render: () => <AccountSettingsPanel activeTab="license" />,
          },
          {
            key: 'api-keys',
            labelKey: 'menu.apiKey',
            path: '/settings/api-keys',
            resourceKey: 'apikey.manage',
            accessUris: ['/OpenData'],
            render: () => <OpenData />,
          },
          {
            key: 'ai-settings',
            labelKey: 'settings.aiSettings',
            path: '/settings/ai-settings',
            resourceKey: 'ai.gateway.config.manage',
            accessUris: ['/OpenData'],
            render: () => (
              <div className={styles['settings-service-page']}>
                <AIProviderConfigPanel />
              </div>
            ),
          },
        ] as SettingsItem[],
      },
      {
        key: 'platform',
        titleKey: 'settings.platform',
        icon: <CloudServices size={16} />,
        items: [
          {
            key: 'users',
            labelKey: 'UserManagement',
            path: '/settings/users',
            resourceKey: 'iam.user.view',
            accessUris: ['/account-management'],
            render: (title) => <AccountManagement title={title} />,
          },
        ] as SettingsItem[],
      },
    ],
    []
  );
  const isLaunchpadStandalone = isLaunchpadStandalonePort();
  const disabledOpenSourceSettings = new Set(['license', 'ai-settings']);
  const sourceGroups = (isLaunchpadStandalone ? groups.filter((group) => group.key === 'user') : groups).map((group) => ({
    ...group,
    items: group.items.filter((item) => !disabledOpenSourceSettings.has(item.key)),
  }));
  const visibleGroups = sourceGroups
    .map((group) => ({
      ...group,
      items: group.items.filter(
        (item) =>
          !hiddenSettingsItemKeys.has(item.key) &&
          canAccess(item, pageUris, currentUserInfo?.superAdmin, systemInfo?.authEnable)
      ),
    }))
    .filter((group) => group.items.length);
  const allItems = visibleGroups.flatMap((group) => group.items);
  const activeItem =
    allItems.find((item) => location.pathname.toLowerCase() === item.path.toLowerCase()) ||
    allItems.find((item) => location.pathname.toLowerCase().startsWith(`${item.path.toLowerCase()}/`)) ||
    allItems[0];
  const platformItemKeys = new Set(groups.find((group) => group.key === 'platform')?.items.map((item) => item.key));
  const isPlatformItem = activeItem?.key ? platformItemKeys.has(activeItem.key) : false;
  const pageClassName = isLaunchpadStandalone
    ? `${styles['settings-page']} ${styles['settings-page-standalone']}`
    : styles['settings-page'];

  return (
    <div className={pageClassName}>
      <aside className={styles['settings-nav']}>
        {visibleGroups.map((group) => (
          <div className={styles['nav-group']} key={group.key}>
            <div className={styles['group-title']}>
              {group.icon}
              <span>{formatMessage(group.titleKey)}</span>
            </div>
            {group.items.map((item) => (
              <button
                type="button"
                key={item.key}
                className={item.key === activeItem?.key ? `${styles['nav-item']} ${styles.active}` : styles['nav-item']}
                onClick={() => navigate(item.path)}
              >
                {formatMessage(item.labelKey)}
              </button>
            ))}
          </div>
        ))}
      </aside>
      <main className={styles['settings-content']}>
        <div className={styles['content-inner']}>
          <div
            className={
              activeItem?.resourceKey
                ? `${styles['embedded-page']} ${isPlatformItem ? styles['platform-embedded-page'] : ''}`
                : `${styles['embedded-page']} ${styles['settings-account-page']}`
            }
          >
            {activeItem?.render(formatMessage(activeItem.labelKey))}
          </div>
        </div>
      </main>
    </div>
  );
};

export default SettingsPage;
