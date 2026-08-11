import { LOGIN_URL, OMC_MODEL } from '@/common-types/constans';
import Layout from '@/layout';
import AccountManagement from '@/pages/account-management';
import AuditLogPage from '@/pages/audit-log';
import FlowPreview from '@/pages/collection-flow/FlowPreview';
import FlowPage, { EventFlowRedirect, SourceFlowRedirect } from '@/pages/flow';
import DevPage from '@/pages/dev-page';
import AnchorFrame from '@/components/anchor-frame';
import DynamicIframe from '@/pages/dynamic-iframe';
import EventFlowPreview from '@/pages/event-flow/FlowPreview.tsx';
import Home from '@/pages/home';
import MenuConfiguration from '@/pages/menu-configuration';
import MqttAuthPage from '@/pages/mqtt-auth';
import FleetNodeDetailPage from '@/pages/fleet-management/FleetNodeDetailPage';
import FleetRuntimePage from '@/pages/fleet-runtime';
import NotPage from '@/pages/not-found-Page';
import NoPermission from '@/pages/not-found-Page/NoPermission';
import NotFoundPage from '@/pages/not-found-Page/NotFoundPage';
import NotebookPage from '@/pages/notebook';
import NotebookEditorPage from '@/pages/notebook/editor';
import OAuthClients from '@/pages/oauth-clients';
import OpenData from '@/pages/open-data';
import PermissionManagement from '@/pages/permission-management';
import RoutingManagement from '@/pages/routing-management';
import LoginPage from '@/pages/login';
import LicenseActivation from '@/pages/license-activation';
import CliAuthPage from '@/pages/cli-auth';
import ClusterPage from '@/pages/cluster';
import Share from '@/pages/share';
import SettingsPage from '@/pages/settings';
import Uns from '@/pages/uns';
import VisionCameraPage from '@/pages/vision';
import { useBaseStore } from '@/stores/base';
import type { ResourceProps, SystemInfoProps, UserInfoProps } from '@/stores/types.ts';
import { isLaunchpadStandaloneAllowedPath, isLaunchpadStandalonePort } from '@/utils/launchpad-site';
import Cookies from 'js-cookie';
import qs from 'qs';
import { lazy, type ReactNode, useEffect } from 'react';
import { type Location, type RouteObject, useNavigate, useRoutes } from 'react-router';
import DynamicMFComponent from '../components/dynamic-mf-component';

// 根路径重定向到外部login页

const defaultPostLoginPath = '/uns';
const disabledOpenSourceRoutePrefixes = ["/home","/account-management","/audit-log","/oauth-clients","/routing-management","/edge-connection","/mqtt-auth","/permission-management","/menu-configuration","/PermissionManagement","/MenuConfiguration","/runtime","/vision","/cluster","/project","/notebook","/launchpad","/license-activation"];
const isDisabledOpenSourceRoute = (path?: string) =>
  Boolean(path && disabledOpenSourceRoutePrefixes.some((prefix) => path === prefix || path.startsWith(prefix + '/')));

const LaunchpadPage = lazy(() => import('@/pages/launchpad'));
const LaunchpadProjectDetail = lazy(() => import('@/pages/launchpad/project-detail'));
const LaunchpadAppDetail = lazy(() => import('@/pages/launchpad/app-detail'));

const AdminOnlyPage = ({ children }: { children: ReactNode }) => {
  const isAdmin = useBaseStore((state) => state.currentUserInfo?.superAdmin === true);
  return isAdmin ? children : <NoPermission />;
};

const FleetRuntimeAdminPage = () => (
  <AdminOnlyPage>
    <FleetRuntimePage />
  </AdminOnlyPage>
);

const FleetNodeDetailAdminPage = () => (
  <AdminOnlyPage>
    <FleetNodeDetailPage />
  </AdminOnlyPage>
);

const ClusterAdminPage = () => (
  <AdminOnlyPage>
    <ClusterPage />
  </AdminOnlyPage>
);

const normalizeRedirectUri = (redirectUri?: string | string[] | null) => {
  if (typeof redirectUri !== 'string') {
    return '';
  }
  const next = redirectUri.trim();
  if (!next || !next.startsWith('/')) {
    return '';
  }
  const [pathname] = next.split('?');
  if (
    next === '/' ||
    next === '/?isLogin=true' ||
    next === '/login' ||
    next === LOGIN_URL ||
    (isLaunchpadStandalonePort() && !isLaunchpadStandaloneAllowedPath(pathname))
  ) {
    return '';
  }
  return next;
};

const RootRedirect = () => {
  const { currentUserInfo, systemInfo } = useBaseStore((state) => ({
    currentUserInfo: state.currentUserInfo,
    systemInfo: state.systemInfo,
  }));
  const params = qs.parse(window.location.search, { ignoreQueryPrefix: true });
  useEffect(() => {
    const handleRedirect = async () => {
      if (params?.isLogin) {
        window.location.replace(normalizeRedirectUri(currentUserInfo?.homePage) || defaultPostLoginPath);
        return;
      }

      if (currentUserInfo?.homePage) {
        window.location.replace(normalizeRedirectUri(currentUserInfo.homePage) || defaultPostLoginPath);
        return;
      }

      if (Cookies.get(OMC_MODEL)) {
        console.warn('omc——cookie失效');
        window.location.replace('/403');
        return;
      }


      window.location.replace(systemInfo?.loginPath || LOGIN_URL);
    };

    handleRedirect();
  }, [params?.isLogin, currentUserInfo?.homePage, systemInfo?.loginPath]);
  return null;
};

export const childrenRoutes = [
  {
    path: '/home',
    Component: Home,
  },
  {
    path: '/runtime',
    Component: FleetRuntimeAdminPage,
    handle: {
      parentPath: '/home',
      code: 'fleet.runtime.title',
    },
  },
  {
    path: '/uns',
    Component: Uns,
  },
  {
    path: '/vision',
    Component: VisionCameraPage,
  },
  // {
  //   path: '/app-display',
  //   Component: AppDisplay,
  // },
  // {
  //   path: '/app-iframe',
  //   Component: AppIframe,
  //   handle: {
  //     parentPath: '/app-display',
  //     code: 'route.appIframe',
  //   },
  // },
  // {
  //   path: '/app-space',
  //   Component: AppSpace,
  // },
  // {
  //   path: '/app-gui',
  //   Component: AppGUI,
  //   handle: {
  //     parentPath: '/app-space',
  //     code: 'route.appGUI',
  //   },
  // },
  // {
  //   path: '/app-preview',
  //   Component: AppPreview,
  //   handle: {
  //     parentPath: '/app-space',
  //     code: 'route.appPreview',
  //   },
  // },
  {
    path: '/flow',
    Component: FlowPage,
  },
  {
    path: '/collection-flow',
    Component: SourceFlowRedirect,
  },
  {
    path: '/collection-flow/flow-editor',
    Component: FlowPreview,
    handle: {
      parentPath: '/flow',
      code: 'route.flowEditor',
      multiInstance: true,
    },
  },
  {
    path: '/EventFlow',
    Component: EventFlowRedirect,
  },
  {
    path: '/event-flow',
    Component: EventFlowRedirect,
  },
  {
    path: '/EventFlow/Editor',
    Component: EventFlowPreview,
    handle: {
      parentPath: '/flow',
      code: 'route.eventFlowEditor',
      multiInstance: true,
    },
  },
  {
    path: '/event-flow/flow-editor',
    Component: EventFlowPreview,
    handle: {
      parentPath: '/flow',
      code: 'route.eventFlowEditor',
      multiInstance: true,
    },
  },
  {
    path: '/event-flow/editor',
    Component: EventFlowPreview,
    handle: {
      parentPath: '/flow',
      code: 'route.eventFlowEditor',
      multiInstance: true,
    },
  },
  {
    path: '/account-management',
    Component: AccountManagement,
  },
  {
    path: '/settings',
    Component: SettingsPage,
    handle: {
      parentPath: '/_common',
      code: 'common.settings',
      type: 'all',
    },
  },
  {
    path: '/settings/:section',
    Component: SettingsPage,
    handle: {
      parentPath: '/settings',
      code: 'common.settings',
      type: 'all',
    },
  },
  {
    path: '/audit-log',
    Component: AuditLogPage,
    handle: {
      parentPath: '/_common',
      code: 'route.auditLog',
      type: 'all',
    },
  },
  {
    path: '/MenuConfiguration',
    Component: MenuConfiguration,
  },
  {
    path: '/PermissionManagement',
    Component: PermissionManagement,
  },
  {
    path: '/dev-page',
    Component: DevPage,
    handle: {
      showName: 'devPage',
      type: 'all',
    },
  },
  // 插件移到主项目 数据开放
  {
    path: '/OpenData',
    Component: OpenData,
  },
  {
    path: '/oauth-clients',
    Component: OAuthClients,
  },
  {
    path: '/routing-management',
    Component: RoutingManagement,
    handle: {
      parentPath: '/_common',
      code: 'route.routingManagement',
    },
  },
  {
    path: '/cluster',
    Component: ClusterAdminPage,
    handle: {
      parentPath: '/_common',
      code: 'menu.cloudSync',
    },
  },
  {
    path: '/403',
    Component: NoPermission,
    handle: {
      parentPath: '/_common',
      showName: '403',
      type: 'all',
    },
  },
  {
    path: '/404',
    element: <NotFoundPage />,
    handle: {
      parentPath: '/_common',
      showName: '404',
      type: 'all',
    },
  },
  {
    path: '/project',
    Component: lazy(() => import('@/pages/project')),
  },
  {
    path: '/project/:projectId',
    Component: lazy(() => import('@/pages/project/Detail')),
    handle: {
      parentPath: '/project',
      code: 'route.projectDetail',
    },
  },
  {
    path: '/notebook',
    Component: NotebookPage,
    handle: {
      parentPath: '/_common',
      showName: 'Notebook',
      code: 'Notebook.title',
      type: 'all',
    },
  },
  {
    path: '/notebook/editor/:id',
    Component: NotebookEditorPage,
    handle: {
      parentPath: '/notebook',
      code: 'Notebook.editor',
      dynamicTabName: true,
    },
  },
  {
    path: '/edge-connection/nodes/:nodeKey',
    Component: FleetNodeDetailAdminPage,
    handle: {
      parentPath: '/edge-connection',
      code: 'menu.mqttAuth',
      multiInstance: true,
    },
  },
  {
    path: '/edge-connection',
    Component: MqttAuthPage,
    handle: {
      parentPath: '/_common',
      code: 'menu.mqttAuth',
      type: 'all',
      multiInstance: true,
    },
  },
  {
    path: '/mqtt-auth',
    Component: MqttAuthPage,
    handle: {
      parentPath: '/_common',
      code: 'menu.mqttAuth',
      type: 'all',
      multiInstance: true,
    },
  },
].filter((route) => !isDisabledOpenSourceRoute(route.path));

export const launchpadChildRoutes = [
  {
    path: '/launchpad',
    Component: LaunchpadPage,
    handle: {
      parentPath: '/_common',
      code: 'Launchpad.title',
      type: 'all',
    },
  },
  {
    path: '/launchpad/:projectName',
    Component: LaunchpadProjectDetail,
    handle: {
      parentPath: '/launchpad',
      code: 'Launchpad.title',
      showName: 'Launchpad.title',
      type: 'all',
      multiInstance: true,
    },
  },
  {
    path: '/launchpad/:projectName/:appName',
    Component: LaunchpadAppDetail,
    handle: {
      parentPath: '/launchpad',
      code: 'Launchpad.title',
      showName: 'Launchpad.title',
      type: 'all',
      multiInstance: true,
    },
  },
];

const launchpadSettingsChildRoutes = [
  {
    path: '/settings',
    Component: SettingsPage,
    handle: {
      parentPath: '/_common',
      code: 'common.settings',
      type: 'all',
    },
  },
  {
    path: '/settings/:section',
    Component: SettingsPage,
    handle: {
      parentPath: '/settings',
      code: 'common.settings',
      type: 'all',
    },
  },
];

const getLaunchpadRoutes = () => [];

const getMainChildrenRoutes = () => childrenRoutes;

// 前端路由路径
export const frontendPathList = [
  ...childrenRoutes.map((item) => item.path),
  ...launchpadChildRoutes.map((item) => item.path),
  ...launchpadSettingsChildRoutes.map((item) => item.path),
];

const getRoutes = () => [
  {
    path: '/',
    element: <RootRedirect />,
  },
  ...getLaunchpadRoutes(),
  {
    path: '/',
    element: <Layout />,
    handle: {
      layout: 'main',
    },
    children: getMainChildrenRoutes(),
  },
  {
    path: LOGIN_URL,
    Component: LoginPage,
  },
  {
    path: '/license-activation',
    Component: LicenseActivation,
  },
  {
    path: '/cli-auth',
    Component: CliAuthPage,
  },
  {
    path: '/share',
    Component: Share,
  },
  {
    path: '*',
    element: <NotPage />,
  },
].filter((route) => !isDisabledOpenSourceRoute(route.path));

export const getRoutesDom = ({
  menuGroup,
  systemInfo,
  currentUserInfo,
}: {
  menuGroup: ResourceProps[];
  systemInfo?: SystemInfoProps;
  currentUserInfo?: UserInfoProps;
}) => {
  return getRoutes().map((route) => {
    if ((route.handle as any)?.layout === 'launchpad') {
      return route;
    }
    if (route.children) {
      return {
        ...route,
        children: [
          // 前端路由
          ...((route.children ?? [])
            ?.map((child) => {
              const info = menuGroup?.find((f) => f.isFrontend && child.path === f?.url);
              if (info) {
                return {
                  ...child,
                  handle: {
                    ...child.handle,
                    path: child.path,
                    // code: info?.code,
                    showName: info?.showName,
                    icon: info?.icon,
                  },
                };
              } else if (child.handle?.parentPath === '/_common') {
                // 开发环境打开方便调试
                if (import.meta.env.DEV) return child;
                if (child.handle?.type === 'all') {
                  return child;
                }
                // 没有正真父级菜单情况
                if (
                  systemInfo?.authEnable &&
                  !currentUserInfo?.pageList?.some((s: any) => s.uri?.toLowerCase?.() === child.path?.toLowerCase?.())
                ) {
                  return null;
                }
                return {
                  ...child,
                  handle: {
                    ...child.handle,
                    path: child.path,
                    code: child.handle?.code ?? child.path,
                  },
                };
              } else if (child.handle?.parentPath) {
                // 没有暴露出去的路由
                return {
                  ...child,
                  handle: {
                    ...child.handle,
                    path: child.path,
                    // multiInstance 路由不设置 code，因为它们的标题来自 location.state.tabName
                    code: (child.handle as any)?.multiInstance ? undefined : (child.handle?.code ?? child.path),
                  },
                };
              } else {
                // 开发环境打开方便调试
                if (import.meta.env.DEV) return child;
                return null;
              }
            })
            ?.filter((f) => f) || []),
          // 后端路由（及前端模块联邦路由）
          ...(menuGroup
            ?.filter((item) => !item.isFrontend)
            ?.map((d) => {
              if (!d) return null;
              // 模块联邦-插件及插件子路由
              if (d.isRemote) {
                const path = d?.remoteModelName ? `/${d?.parentCode}/${d?.remoteModelName}` : '/' + d?.code;
                return {
                  path,
                  Component: DynamicMFComponent,
                  handle: {
                    key: d?.code,
                    code: d?.code,
                    showName: d?.showName,
                    icon: d?.icon,
                    path,
                    // 模块联邦子模块
                    moduleName: d?.remoteModelName,
                    parentPath: '/' + d?.parentCode,
                  },
                };
              }
              // Anchor 3D 子应用走专用宿主（bridge 上下文/生命周期同步、加载态与同源白名单）
              if (d?.code?.startsWith('anchor.') || d?.url?.startsWith('/anchor')) {
                return {
                  path: '/' + d?.code,
                  element: <AnchorFrame url={d?.url} name={d?.showName} />,
                  handle: {
                    openType: d?.openType,
                    key: d?.code,
                    code: d?.code,
                    showName: d?.showName,
                    icon: d?.icon,
                    path: '/' + d?.code,
                  },
                };
              }
              return {
                path: '/' + d?.code,
                element: <DynamicIframe url={d?.url} name={d?.showName} code={d?.code} />,
                handle: {
                  openType: d?.openType,
                  key: d?.code,
                  code: d?.code,
                  showName: d?.showName,
                  icon: d?.icon,
                  path: '/' + d?.code,
                },
              };
            })
            ?.filter((f) => f) || []),
        ],
      };
    } else {
      return route;
    }
  }) as RouteObject[];
};

export const RoutesElement = ({ routeDom }: { routeDom: RouteObject[] }) => {
  return useRoutes(routeDom);
};

export const useLocationNavigate = () => {
  const navigate = useNavigate();
  return (location: Partial<Location>) => {
    const { pathname, search, state } = location;
    navigate(pathname + (search ?? ''), { state });
  };
};
