import { useState } from 'react';
import { useLocation, useOutlet, type Location, matchRoutes } from 'react-router';
import { useDeepCompareEffect } from 'ahooks';
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
  // 路由，和pathname区别是，详情页 path /:id，routePath是 /1
  path: string;
  // location对象，存储起来用来二次导航
  location: Location;
  // 模块联邦名称
  moduleName?: string;
  // 模块联邦父级菜单
  parentPath?: string;
}

// 匹配路由，拿到信息
export function useMatchRoute(): MatchRouteType | undefined {
  const { menuGroup, systemInfo, currentUserInfo } = useBaseStore((state) => ({
    menuGroup: state.menuGroup,
    systemInfo: state.systemInfo,
    currentUserInfo: state.currentUserInfo,
  }));
  // 获取路由组件实例
  const children = useOutlet();
  // 获取嵌套路由信息
  const formatMessage = useTranslate();
  // 获取当前url
  const location = useLocation();

  const [matchRoute, setMatchRoute] = useState<MatchRouteType | undefined>();

  // 监听pathname变了，说明路由有变化，重新匹配，返回新路由信息
  useDeepCompareEffect(() => {
    // 获取当前匹配的路由
    const matches = matchRoutes(getRoutesDom({ menuGroup, systemInfo, currentUserInfo }), location.pathname) || [];
    const lastRoute = matches.at(-1)?.route;

    if (!lastRoute?.handle) return;

    setMatchRoute({
      title: formatShowName({
        code: (lastRoute?.handle as any)?.code,
        showName: (lastRoute?.handle as any)?.showName,
        formatMessage,
      }),
      icon: (lastRoute?.handle as any)?.icon,
      path: (lastRoute?.handle as any)?.path,
      pathname: location.pathname,
      children,
      routePath: lastRoute?.path || '',
      moduleName: (lastRoute?.handle as any)?.moduleName,
      parentPath: (lastRoute?.handle as any)?.parentPath,
      location,
    });
  }, [location, menuGroup, systemInfo, currentUserInfo]);

  return matchRoute;
}
