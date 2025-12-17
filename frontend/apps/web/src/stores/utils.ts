import type { ContainerItemProps, ResourceProps } from './types.ts';
import { frontendPathList } from '@/routers';
import { storageOpt } from '@/utils';
import {
  SUPOS_USER_GUIDE_ROUTES,
  SUPOS_USER_LAST_LOGIN_ENABLE,
  SUPOS_USER_TIPS_ENABLE,
} from '@/common-types/constans.ts';
import { filter, includes, isBoolean, isEmpty, map } from 'lodash-es';
import { setUserTipsEnable } from '@/stores/base';

/**
 * 通用多分组方法
 * @param data - 要分组的数组
 * @param criteria - 分组标识对象，键是组名，值是判断逻辑的回调函数
 * @returns 分组后的结果对象
 */
type Predicate<T> = (item: T) => boolean;
export type Criteria<T> = Record<string, Predicate<T>>;

export function multiGroupByCondition<T>(data: T[], criteria: Criteria<T>): Record<string, T[]> {
  if (!data?.length) {
    return {};
  }
  const groups: Record<string, T[]> = Object.keys(criteria).reduce(
    (acc, key) => {
      acc[key] = [];
      return acc;
    },
    {} as Record<string, T[]>
  );

  groups.others = [];

  data.forEach((item) => {
    let matched = false;

    for (const [key, predicate] of Object.entries(criteria)) {
      if (predicate(item)) {
        groups[key].push(item);
        matched = true;
        break;
      }
    }

    if (!matched) {
      groups.others.push(item);
    }
  });

  return groups;
}

// 拒绝优先
export function filterArrays(arr1: string[] = [], arr2: string[] = []) {
  return arr2?.filter?.((item) => !arr1?.includes?.(item));
}

// 拒绝优先
export function filterObjectArrays(arr1: any[] = [], arr2: any[] = []) {
  return arr2.filter((item2) => !arr1.some((item1) => item1.uri === item2.uri));
}

// 判断按钮权限函数
export function matchUriWithPattern(uri: string, pattern: string) {
  // 转换模式（如 button:uns.*）为正则表达式
  const regex = new RegExp('^' + pattern.replace('*', '.*').replace(':', '\\:') + '$');
  return regex.test(uri);
}

// 处理接口返回的模式数组
export function handleButtonPermissions(patterns: string[] = [], allButtonGroup: any[]) {
  const matchedButtons: string[] = [];
  const permissions = allButtonGroup?.map((i: any) => `button:${i.code}`);
  // 如果模式包含 'button:*'，返回所有按钮权限
  if (patterns?.includes?.('button:*')) {
    return [...new Set(permissions)];
  }

  // 遍历所有的模式
  patterns.forEach((pattern) => {
    permissions.forEach((uri) => {
      if (matchUriWithPattern(uri, pattern)) {
        matchedButtons.push(uri); // 匹配成功的按钮加入到结果数组
      }
    });
  });
  return [...new Set(matchedButtons)];
}

// 拆分关于我们和高阶使用
export const filterContainerList = (containerMap: { [key: string]: ContainerItemProps } = {}) => {
  const containerList = Object.values(containerMap);
  const _containerList = containerList?.filter((f) => f.envMap?.service_is_show);
  return {
    advancedUse: _containerList?.filter((f) => f.envMap?.service_redirect_url) || [],
    aboutUs: _containerList || [],
  };
};

function buildSortedTree(data: ResourceProps[]) {
  const map: Record<string, ResourceProps> = {};
  data.forEach((item: ResourceProps) => {
    map[item.id] = { ...item, children: [] };
  });
  const tree: ResourceProps[] = [];
  data.forEach((item) => {
    const node = map[item.id];
    const parentId = item.parentId;
    if (parentId && map[parentId]) {
      map[parentId].children!.push(node);
    } else {
      // 根节点直接加入树
      tree.push(node);
    }
  });
  const sortTree = (nodes: ResourceProps[]): ResourceProps[] => {
    nodes.sort((a, b) => a.sort - b.sort); // 按 sort 升序排序
    nodes.forEach((node) => {
      if (node.children && node.children.length) sortTree(node.children);
    });
    return nodes;
  };

  return sortTree(tree);
}

export function buildResourceTrees(resources: ResourceProps[]) {
  const menuGroup = resources.filter((r) => r.type === 2 || r.type === 5);
  const homeGroup = menuGroup?.filter((r) => r.homeEnable);
  const homeTabGroup = resources?.filter((r) => r.type === 4 && r.homeEnable);
  const treeResources = resources.filter((r) => r.type !== 5 && r.type !== 4);
  return {
    // 菜单分组 不带子菜单，过滤掉空目录
    menuTree: buildSortedTree(treeResources)?.filter((f) => f.type === 2 || (f.type === 1 && f?.children?.length)),
    // home页分组 不带子菜单，过滤掉空目录
    homeTree: buildSortedTree(treeResources?.filter((r) => r.homeEnable))?.filter(
      (f) => f.type === 2 || (f.type === 1 && f?.children?.length)
    ),
    // options
    menuGroup,
    homeGroup,
    homeTabGroup,
  };
}

// 重组菜单，比如是三级菜单进行重组
export function mapResource(source: ResourceProps[] = []) {
  return source.map((item) => {
    const parent = source?.find((f) => f.id === item.parentId);
    // 菜单下面还挂着菜单特殊处理，比如 /collect 下面挂在/collect/detail
    if (parent && item?.type === 5) {
      const isFrontend = frontendPathList?.includes(parent.url! + item.url!);
      const isRemote = !isFrontend && item.urlType === 1;
      // 特殊处理subMen
      return {
        ...item,
        subMenu: true,
        parentCode: parent?.code,
        url: parent.url! + item.url!,
        isFrontend,
        isRemote,
        remoteModelName: isRemote ? item.url?.slice(1) : undefined,
      };
    } else {
      const isFrontend = frontendPathList?.includes(item.url!);
      return {
        ...item,
        parentCode: parent?.code,
        isFrontend,
        isRemote: !isFrontend && item.urlType === 1,
      };
    }
  });
}

// 根据用户路由资源组 匹配出 路由 - 如果免登或者超级管理员，不进行过滤; type为1既是目录不进行过滤;type为4的也不控制
export function filterRouteByUserResource(source: any[] = [], target: any[] = [], authEnable?: boolean) {
  if (authEnable)
    return source.filter(
      (sourceItem) =>
        sourceItem.type === 1 ||
        sourceItem.type === 4 ||
        sourceItem.subMenu === true ||
        target.some((targetItem) => {
          if (sourceItem?.urlType !== 1) {
            return '/' + sourceItem?.code?.toLowerCase?.() === targetItem.uri?.toLowerCase?.();
          } else {
            return sourceItem?.url?.toLowerCase?.() === targetItem.uri?.toLowerCase?.();
          }
        })
    );
  return source;
}

// 包含新手导航的页面路由集合，新增页面导航时务必在这里添加url
const GuidePagePaths = ['/home', '/uns'];
// 新手指引设置
export function guideConfig({ systemInfo, menuGroup, info }: { systemInfo: any; menuGroup: any; info: any }) {
  // 1.新手导航：根据authenable和token区分是否为免登录
  //      a.先获取上次免登录状态和当前比较，如果发生改变，则说明用户登录发生变化（比如由需要登陆变为免登或者免登变为需要），需要清除之前的SUPOS_USER_GUIDE_ROUTES状态，并设置新的免登状态
  //      b.然后判断当前是否为免登
  //          如果是免登录，先判断SUPOS_USER_GUIDE_ROUTES是否存在，不存在，则添加，存在则不做处理
  //          如果需要登陆，再按原有逻辑（用户第一次登录）进行引导
  // 2.tips: 用户访问时进入系统则展示tips，且可勾选不再展示（每次登录或者每次免登状态都需考虑）
  //         1).判断是否免登 2).是否为刚登录 3).判断用户是否支持展示
  const lastLoginEnable = storageOpt.getOrigin(SUPOS_USER_LAST_LOGIN_ENABLE);

  // 用户是guest表明登录成功
  const token = info?.sub === 'guest' ? undefined : 'login';
  const isLoginEnable = isBoolean(systemInfo?.authEnable) && !isEmpty(token);
  // 获取上次免登录状态和当前比较，如果发生改变，则说明用户登录发生变化,或者systemInfo?.authEnable获取失败则清除缓存
  if (!isBoolean(systemInfo?.authEnable) || lastLoginEnable !== `${isLoginEnable}`) {
    storageOpt.remove(SUPOS_USER_GUIDE_ROUTES);
    storageOpt.setOrigin(
      SUPOS_USER_LAST_LOGIN_ENABLE,
      `${isBoolean(systemInfo?.authEnable) ? isLoginEnable : systemInfo?.authEnable}`
    );
    storageOpt.remove(SUPOS_USER_TIPS_ENABLE);
  }

  // 是否为免登录：authEnable===false并且不存在token时
  const notLogin = systemInfo?.authEnable === false && !token;
  // 如果为免登录，则判断是否存在新手导航数据，存在则继续触发，不存在则添加
  if (notLogin) {
    if (!storageOpt.getOrigin(SUPOS_USER_TIPS_ENABLE)) {
      setUserTipsEnable('1');
    }
    if (!storageOpt.get(SUPOS_USER_GUIDE_ROUTES)) {
      storageOpt.set(
        SUPOS_USER_GUIDE_ROUTES,
        map(
          filter(menuGroup, (r) => includes(GuidePagePaths, r?.url)),
          (route) => ({ ...route, isVisited: false })
        )
      );
    }
  }
  // 如果是登录状态
  if (isLoginEnable) {
    // 判断用户是否手动禁用tips展示
    const tipsEnable = info?.tipsEnable;
    if (tipsEnable && !storageOpt.getOrigin(SUPOS_USER_TIPS_ENABLE)) {
      setUserTipsEnable('1');
    }
    if (!tipsEnable) {
      setUserTipsEnable('0');
    }
    const isFirstLogin = info?.firstTimeLogin;
    // 首次登录且未初始化用户引导路由信息，则需初始化该信息；已经初始化则继续使用缓存的状态
    if (isFirstLogin === 1 && !storageOpt.get(SUPOS_USER_GUIDE_ROUTES)) {
      storageOpt.set(
        SUPOS_USER_GUIDE_ROUTES,
        map(
          filter(menuGroup, (r) => includes(GuidePagePaths, r?.url)),
          (route) => ({ ...route, isVisited: false })
        )
      );
    }
    // 由于存在手动启用新手导航功能，先取消清除的逻辑
    // if (isFirstLogin !== 1) {
    //   // 非首次登录直接清除用户引导路由信息
    //   storageOpt.remove(SUPOS_USER_GUIDE_ROUTES);
    // }
  }
}
