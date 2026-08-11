import {
  getLaunchpadAppByNameApi,
  getLaunchpadProjectByNameApi,
  recordLaunchpadAppViewApi,
} from '@/apis/core-api/launchpad';
import ComDetailHeader from '@/components/com-detail-header';
import { useActivate } from '@/contexts/tabs-lifecycle-context';
import { useTranslate } from '@/hooks';
import type { App } from '@/pages/launchpad/type';
import { useLocationNavigate } from '@/routers';
import { normalizeLocalAppUrl } from '@/utils/app-url';
import { getToken } from '@/utils/auth';
import { useBaseStore } from '@/stores/base';
import { Maximize, Minimize } from '@/components/lucide-icon/carbon';
import { Button, Result, Spin } from 'antd';
import { type FC, useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useLocation, useNavigate, type Location } from 'react-router';
import AppSwitcher from './components/AppSwitcher';
import styles from './app-detail.module.scss';

interface AppDetailPageProps {
  projectName?: string;
  appName?: string;
  location?: Location;
  onBack?: () => void;
  onAppChange?: (appName: string, app?: App) => void;
}

const decodePathSegment = (value?: string) => {
  if (!value) return '';
  try {
    return decodeURIComponent(value);
  } catch {
    return value;
  }
};

const getLaunchpadAppRouteParams = (pathname: string) => {
  const parts = pathname.split('/').filter(Boolean);
  if (parts[0] !== 'launchpad') {
    return { projectName: '', appName: '' };
  }
  return {
    projectName: decodePathSegment(parts[1]),
    appName: decodePathSegment(parts[2]),
  };
};

const AppDetailPage: FC<AppDetailPageProps> = ({
  projectName: propProjectName,
  appName: propAppName,
  location: keepAliveLocation,
  onBack,
  onAppChange,
}) => {
  const navigate = useLocationNavigate();
  const routerNavigate = useNavigate();
  const routerLocation = useLocation();
  const location = keepAliveLocation || routerLocation;
  const routeParams = useMemo(() => getLaunchpadAppRouteParams(location.pathname), [location.pathname]);
  const projectName = propProjectName || routeParams.projectName;
  const appName = propAppName || routeParams.appName;
  const iframeRef = useRef<HTMLIFrameElement>(null);
  const isMountedRef = useRef(true);
  const appDetailRequestIdRef = useRef(0);
  const [app, setApp] = useState<App | null>(null);
  const [appNotFound, setAppNotFound] = useState(false);
  const [noPermission, setNoPermission] = useState(false);
  const [apps, setApps] = useState<App[]>([]);
  const [loading, setLoading] = useState(!!appName);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const formatMessage = useTranslate();
  // Builder（auth/me.roleCode=admin/builder → superAdmin）版异常页多一个 Go to Project；Operator 仅 Back to Launchpad
  const isBuilder = useBaseStore((state) => state.currentUserInfo?.superAdmin === true);

  const fetchAppDetail = useCallback(() => {
    if (!projectName || !appName) return;
    const requestId = ++appDetailRequestIdRef.current;
    setLoading(true);
    setAppNotFound(false);
    setNoPermission(false);
    getLaunchpadAppByNameApi(projectName, appName)
      .then((data) => {
        if (!isMountedRef.current || requestId !== appDetailRequestIdRef.current) return;
        setApp(data?.app || null);
        setAppNotFound(false);
        setNoPermission(false);
      })
      .catch((error) => {
        if (!isMountedRef.current || requestId !== appDetailRequestIdRef.current) return;
        setApp(null);
        // 403 = App 存在但成员缺 App Role → App access required；
        // 其余（含 404 不存在/已删除/非成员，不披露存在性）→ App not found
        if (error?.code === 403) {
          setNoPermission(true);
          setAppNotFound(false);
        } else {
          setAppNotFound(true);
          setNoPermission(false);
        }
      })
      .finally(() => {
        if (isMountedRef.current && requestId === appDetailRequestIdRef.current) setLoading(false);
      });
  }, [projectName, appName]);

  const fetchProjectApps = useCallback(() => {
    if (!projectName) return;
    getLaunchpadProjectByNameApi(projectName)
      .then((data) => {
        if (isMountedRef.current) setApps(data?.apps || []);
      })
      .catch(() => {
        if (!isMountedRef.current) return;
        setApps([]);
      });
  }, [projectName]);

  const recordAppView = useCallback((appId?: string | number) => {
    if (!appId) return;
    void recordLaunchpadAppViewApi(appId);
  }, []);

  useEffect(() => {
    isMountedRef.current = true;
    return () => {
      isMountedRef.current = false;
    };
  }, []);

  useEffect(() => {
    fetchAppDetail();
  }, [fetchAppDetail]);

  useEffect(() => {
    recordAppView(app?.appId);
  }, [app?.appId, recordAppView]);

  useEffect(() => {
    if (!app) return;
    if (routerLocation.pathname !== location.pathname || routerLocation.search !== location.search) return;
    const tabName = app.displayName || app.appName || appName;
    if (!tabName || location.state?.tabName === tabName) return;
    routerNavigate(location.pathname + location.search, {
      replace: true,
      state: {
        ...(location.state || {}),
        tabName,
      },
    });
  }, [
    app,
    appName,
    location.pathname,
    location.search,
    location.state,
    routerLocation.pathname,
    routerLocation.search,
    routerNavigate,
  ]);

  // 只在 projectName 变化时重新获取 apps 列表，切换 app 不重新请求
  useEffect(() => {
    fetchProjectApps();
  }, [fetchProjectApps]);

  useActivate(() => {
    fetchAppDetail();
    fetchProjectApps();
  });

  useEffect(() => {
    const handleFullscreenChange = () => {
      setIsFullscreen(!!document.fullscreenElement);
    };
    document.addEventListener('fullscreenchange', handleFullscreenChange);
    return () => document.removeEventListener('fullscreenchange', handleFullscreenChange);
  }, []);

  const handleBack = () => {
    if (onBack) {
      onBack();
      return;
    }
    if (appNotFound) {
      routerNavigate('/launchpad', {
        replace: true,
        state: { tabName: formatMessage('Launchpad.title') },
      });
      return;
    }
    navigate({
      pathname: `/launchpad/${encodeURIComponent(projectName || '')}`,
      state: { tabName: projectName || formatMessage('Launchpad.title') },
    });
  };

  const handleAppChange = (newAppName: string, nextApp?: App) => {
    if (onAppChange) {
      onAppChange(newAppName, nextApp);
    } else {
      routerNavigate(`/launchpad/${encodeURIComponent(projectName || '')}/${encodeURIComponent(newAppName)}`, {
        replace: true,
        state: {
          tabName: nextApp?.displayName || nextApp?.appName || newAppName,
          replaceCurrentTab: true,
        },
      });
    }
  };

  const toggleFullscreen = () => {
    if (!isFullscreen) {
      const elem = iframeRef.current?.parentElement;
      if (elem?.requestFullscreen) elem.requestFullscreen();
      setIsFullscreen(true);
    } else {
      document.exitFullscreen?.();
      setIsFullscreen(false);
    }
  };

  // 头部始终渲染，title 优先用已有 app state，loading 时用 appName 作为占位
  const title = app?.displayName || app?.appName || appName || '';
  const rawAppUrl = app?.url;

  const appUrl = useMemo(() => {
    if (!rawAppUrl) return '';
    const shouldRewritePrivateHost = app?.siteType === 'dynamic' && !app?.manual && app?.appType !== 'manual';
    const normalizedUrl = normalizeLocalAppUrl(rawAppUrl, { allowPrivateHostRewrite: shouldRewritePrivateHost });
    const token = getToken();
    const separator = normalizedUrl.includes('?') ? '&' : '?';
    return token ? `${normalizedUrl}${separator}_token=${encodeURIComponent(token)}` : normalizedUrl;
  }, [app?.appType, app?.manual, app?.siteType, rawAppUrl]);

  // App not found：不存在 / 已删除 / 非 Workspace Member（不披露存在性）
  if (appNotFound) {
    return (
      <div className={styles.appDetailPage}>
        <Result
          status="404"
          title={formatMessage('Launchpad.appNotFound', {}, 'App not found')}
          subTitle={formatMessage(
            'Launchpad.appNotFoundDesc',
            {},
            'This app is unavailable. It may have been removed, or the link is no longer valid.'
          )}
          extra={
            <Button type="primary" onClick={() => routerNavigate('/launchpad', { replace: true })}>
              {formatMessage('Launchpad.backToLaunchpad', {}, 'Back to Launchpad')}
            </Button>
          }
        />
      </div>
    );
  }

  // App access required：成员但缺 App 绑定的 Business Role（Builder 版多 Go to Project，Operator 仅 Back）
  if (noPermission) {
    return (
      <div className={styles.appDetailPage}>
        <Result
          status="403"
          title={formatMessage('Launchpad.appAccessTitle', {}, 'App access required')}
          subTitle={formatMessage(
            isBuilder ? 'Launchpad.appAccessDescBuilder' : 'Launchpad.appAccessDescOperator',
            {},
            isBuilder
              ? 'You do not have the Business Role required to open this app. Configure role permissions in the project, or contact the project owner.'
              : 'You have not been assigned the role required to access this app. Please contact the project owner or administrator to configure role permissions.'
          )}
          extra={
            <>
              {isBuilder && (
                <Button
                  type="primary"
                  onClick={() =>
                    navigate({
                      pathname: `/launchpad/${encodeURIComponent(projectName || '')}`,
                      state: { tabName: projectName || formatMessage('Launchpad.title') },
                    })
                  }
                >
                  {formatMessage('Launchpad.goToProject', {}, 'Go to Project')}
                </Button>
              )}
              <Button
                type={isBuilder ? 'default' : 'primary'}
                onClick={() => routerNavigate('/launchpad', { replace: true })}
              >
                {formatMessage('Launchpad.backToLaunchpad', {}, 'Back to Launchpad')}
              </Button>
            </>
          }
        />
      </div>
    );
  }

  return (
    <div className={styles.appDetailPage}>
      <ComDetailHeader
        title={title}
        onBack={handleBack}
        showBack
        showDesc={false}
        rightExtra={
          <>
            <AppSwitcher apps={apps} currentAppName={appName} onAppChange={handleAppChange} />
            <Button
              type="text"
              icon={isFullscreen ? <Minimize size={16} /> : <Maximize size={16} />}
              onClick={toggleFullscreen}
              className={styles.fullscreenButton}
            />
          </>
        }
        style={{
          border: '1px solid var(--ui-menu-hover-color)',
          alignItems: 'center',
        }}
      />
      <div className={styles.iframeContainer}>
        {loading ? (
          <div className={styles.loading}>
            <Spin size="large" />
          </div>
        ) : !app ? (
          <div className={styles.error}>{formatMessage('Launchpad.appNotFound')}</div>
        ) : !appUrl ? (
          <div className={styles.loading}>{formatMessage('project.appNotReady')}</div>
        ) : (
          <iframe ref={iframeRef} src={appUrl} className={styles.iframe} title={app.displayName || app.appName} />
        )}
      </div>
    </div>
  );
};

export default AppDetailPage;
