import { type Location, type RouteObject, useNavigate, useRoutes } from 'react-router';
import Layout from '@/layout';
import Uns from '@/pages/uns';
import Todo from '@/pages/todo';
import GrafanaDesign from '@/pages/grafana-design';
// import AppDisplay from '@/pages/app-management/AppDisplay';
// import AppSpace from '@/pages/app-management/AppSpace';
// import AppGUI from '@/pages/app-management/AppGUI';
// import AppPreview from '@/pages/app-management/AppPreview';
// import AppIframe from '@/pages/app-management/AppIframe';
import NotFoundPage from '@/pages/not-found-Page/NotFoundPage';
import NotPage from '@/pages/not-found-Page';
import CollectionFlow from '@/pages/collection-flow';
import FlowPreview from '@/pages/collection-flow/FlowPreview';
import Dashboards from '@/pages/dashboards';
import DashboardsPreview from '@/pages/dashboards/DashboardsPreview';
import Localization from '@/pages/localization';
import MenuConfiguration from '@/pages/menu-configuration';
import Home from '@/pages/home';
import AccountManagement from '@/pages/account-management';
import AboutUs from '@/pages/aboutus';
import AdvancedUse from '@/pages/advanced-use';
import DevPage from '@/pages/dev-page';
import NoPermission from '@/pages/not-found-Page/NoPermission';
import { LOGIN_URL, OMC_MODEL } from '@/common-types/constans';
import Share from '@/pages/share';
import EventFlow from '@/pages/event-flow';
import EventFlowPreview from '@/pages/event-flow/FlowPreview.tsx';
import PluginManagement from '@/pages/plugin-management';
import qs from 'qs';
import { useEffect } from 'react';
import { useBaseStore } from '@/stores/base';
import { setToken } from '@/utils';
import DynamicMFComponent from '../components/dynamic-mf-component';
import DynamicIframe from '@/pages/dynamic-iframe';
import type { ResourceProps, SystemInfoProps, UserInfoProps } from '@/stores/types.ts';
import Cookies from 'js-cookie';

// 根路径重定向到外部login页

const RootRedirect = () => {
  const { currentUserInfo, systemInfo } = useBaseStore((state) => ({
    currentUserInfo: state.currentUserInfo,
    systemInfo: state.systemInfo,
  }));
  const params = qs.parse(window.location.search, { ignoreQueryPrefix: true });
  useEffect(() => {
    if (params?.isLogin) {
      window.location.href = currentUserInfo?.homePage || '/uns';
    } else {
      if (Cookies.get(OMC_MODEL)) {
        console.warn('omc——cookie失效');
        window.location.href = '/403';
      } else {
        console.log('登录cookie不存在，要跳转到登录页');
        window.location.href = systemInfo?.loginPath || LOGIN_URL;
      }
    }
  }, [params?.isLogin]);
  return null;
};

const FreeLoginLoader = () => {
  const params = qs.parse(window.location.search, { ignoreQueryPrefix: true });
  useEffect(() => {
    if (params?.token) {
      setToken(params.token as string, {
        expires: 365,
      });
      if (params?.redirectUri) {
        window.location.href = (params?.redirectUri as string) || '/?isLogin=true';
      } else {
        window.location.href = '/?isLogin=true';
      }
    } else {
      window.location.href = '/403';
    }
  }, [params?.token]);
  return null;
};

export const childrenRoutes = [
  {
    path: '/home',
    Component: Home,
  },
  {
    path: '/uns',
    Component: Uns,
  },
  {
    path: '/todo',
    Component: Todo,
    handle: {
      parentPath: '/_common',
      code: 'common.taskCenter',
      type: 'all',
    },
  },
  {
    path: '/grafana-design',
    Component: GrafanaDesign,
    handle: {
      parentPath: '/_common',
      code: 'common.grafanaDesign',
      type: 'all',
    },
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
    path: '/collection-flow',
    Component: CollectionFlow,
  },
  {
    path: '/collection-flow/flow-editor',
    Component: FlowPreview,
    handle: {
      parentPath: '/collection-flow',
      code: 'route.flowEditor',
    },
  },
  {
    path: '/EventFlow',
    Component: EventFlow,
  },
  {
    path: '/EventFlow/Editor',
    Component: EventFlowPreview,
    handle: {
      parentPath: '/EventFlow',
      code: 'route.eventFlowEditor',
    },
  },
  {
    path: '/dashboards',
    Component: Dashboards,
  },
  {
    path: '/dashboards/preview',
    Component: DashboardsPreview,
    handle: {
      parentPath: '/dashboards',
      code: 'route.dashboardsPreview',
    },
  },
  {
    path: '/account-management',
    Component: AccountManagement,
  },
  {
    path: '/aboutus',
    Component: AboutUs,
  },
  {
    path: '/Localization',
    Component: Localization,
  },
  {
    path: '/MenuConfiguration',
    Component: MenuConfiguration,
  },
  {
    path: '/advanced-use',
    Component: AdvancedUse,
  },
  {
    path: '/dev-page',
    Component: DevPage,
    handle: {
      showName: 'devPage',
      type: 'dev',
    },
  },
  {
    path: '/plugin-management',
    Component: PluginManagement,
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
];

// 前端路由路径
export const frontendPathList = childrenRoutes?.map((item) => item.path);

const routes = [
  {
    path: '/',
    element: <RootRedirect />,
  },
  {
    path: '/',
    element: <Layout />,
    children: childrenRoutes,
  },
  {
    path: '/freeLogin',
    element: <FreeLoginLoader />,
    // 数据路由无法使用
    // loader: ({ request }: any) => {
    //   const url = new URL(request.url);
    //   const token = url.searchParams.get('token');
    //   if (token) {
    //     console.log('123');
    //     // 免登录逻辑
    //     // 21600秒 = 6小时 = 0.25天
    //     setToken(token, {
    //       expires: 0.25,
    //     });
    //     return redirect('/?isLogin=true');
    //   }
    //   return null;
    // },
  },
  {
    path: '/share',
    Component: Share,
  },
  {
    path: '*',
    element: <NotPage />,
  },
];

export const getRoutesDom = ({
  menuGroup,
  systemInfo,
  currentUserInfo,
}: {
  menuGroup: ResourceProps[];
  systemInfo?: SystemInfoProps;
  currentUserInfo?: UserInfoProps;
}) => {
  return routes.map((route, index) => {
    if (index === 1 && route.children) {
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
                    code: child.handle?.code ?? child.path,
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
