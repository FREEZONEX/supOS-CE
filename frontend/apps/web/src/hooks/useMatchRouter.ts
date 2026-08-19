import { useMemo } from 'react';
import { useLocation, useOutlet, type Location, matchRoutes } from 'react-router';
import { useTranslate } from '@/hooks';
import { useBaseStore } from '@/stores/base';
import { getRoutesDom } from '@/routers';
import { formatShowName } from '@/utils';

interface MatchRouteType {
  // 菜单名称
  title: string;
  // 要渲染的组件
  // 图标
  icon?: string;
  children: any;
  // tab对应的url
  pathname: string;
  // tab的key，目前和pathname一样
  routePath: string;
  tabKey: string;
  // 路由，和pathname区别是，详情页 path /:id，routePath是 /1
  path: string;
  // location对象，存储起来用来二次导航
  location: Location;
  // 模块联邦名称
  moduleName?: string;
  // 模块联邦父级菜单
  parentPath?: string;
  // 是否支持多实例tab（同一路由可以打开多个tab）
  multiInstance?: boolean;
  // 是否允许当前路由通过 location.state.tabName 更新页签标题
  dynamicTabName?: boolean;
  // 是否允许关闭页签
  closable?: boolean;
}

function getMultiInstanceTabKey(location: Location) {
  const searchParams = new URLSearchParams(location.search || '');
  const id = searchParams.get('id');
  const nodeKey = searchParams.get('nodeKey');
  const instanceKey = id || nodeKey;
  return instanceKey ? `${location.pathname}?${id ? 'id' : 'nodeKey'}=${instanceKey}` : location.pathname;
}

// 匹配路由，拿到信息
export function useMatchRoute(): MatchRouteType | undefined {
  const { menuGroup, systemInfo, currentUserInfo } = useBaseStore((state) => ({
    menuGroup: state.menuGroup,
    systemInfo: state.systemInfo,
    currentUserInfo: state.currentUserInfo,
  }));
  const children = useOutlet();
  const formatMessage = useTranslate();
  const location = useLocation();

  return useMemo(() => {
    const matches = matchRoutes(getRoutesDom({ menuGroup, systemInfo, currentUserInfo }), location.pathname) || [];
    const lastRoute = matches.at(-1)?.route;

    if (!lastRoute?.path && !lastRoute?.handle) return undefined;

    const normalizedRoutePath = lastRoute?.path?.startsWith('/') ? lastRoute.path : `/${lastRoute?.path || ''}`;
    const routeHandle = (lastRoute?.handle as any) || {};
    const multiInstance = routeHandle.multiInstance;
    const dynamicTabName = routeHandle.dynamicTabName === true || multiInstance;
    const routeTabKey = typeof routeHandle.tabKey === 'string' ? routeHandle.tabKey : undefined;
    const fallbackShowName =
      typeof lastRoute?.path === 'string' && lastRoute.path !== '*'
        ? lastRoute.path.replace(/^\//, '')
        : location.pathname;

    return {
      title: formatShowName({
        code: routeHandle.code,
        showName: routeHandle.showName,
        formatMessage,
        finallyShowName: fallbackShowName,
      }),
      icon: routeHandle.icon,
      path: routeHandle.path || normalizedRoutePath,
      pathname: location.pathname,
      children,
      routePath: normalizedRoutePath,
      tabKey: routeTabKey || (multiInstance ? getMultiInstanceTabKey(location) : normalizedRoutePath),
      moduleName: routeHandle.moduleName,
      parentPath: routeHandle.parentPath,
      multiInstance,
      dynamicTabName,
      closable: routeHandle.closable !== false,
      location,
    };
  }, [children, currentUserInfo, formatMessage, location, menuGroup, systemInfo]);
}
