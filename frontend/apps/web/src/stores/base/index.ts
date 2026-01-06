import { type UseBoundStoreWithEqualityFn, createWithEqualityFn } from 'zustand/traditional';
import type { StoreApi } from 'zustand';
import { shallow } from 'zustand/vanilla/shallow';
import type { DataItem, ResourceProps, UserInfoProps } from '@/stores/types.ts';
import { storageOpt } from '@/utils/storage';
import { APP_TITLE, SUPOS_LANG, SUPOS_UNS_TREE, SUPOS_USER_TIPS_ENABLE } from '@/common-types/constans.ts';
import { getPersonConfigApi } from '@/apis/inter-api/uns.ts';
import { getSystemConfig } from '@/apis/inter-api/system-config.ts';
import { getUserInfo } from '@/apis/inter-api/auth';

import type { TBaseStore } from '@/stores/base/type.ts';
import { initI18n, defaultLanguage, useI18nStore } from '../i18n-store.ts';
import { getRoutesResourceApi } from '@/apis/inter-api/resource.ts';
import {
  type Criteria,
  filterArrays,
  filterContainerList,
  filterObjectArrays,
  guideConfig,
  handleButtonPermissions,
  multiGroupByCondition,
  buildResourceTrees,
  filterRouteByUserResource,
  mapResource,
} from '../utils.ts';
import { getLangListApi } from '@/apis/inter-api/i18n.ts';

/**
 * 获取语言包
 */
export const getLangList = async () => {
  try {
    const langList = await getLangListApi();
    useI18nStore.setState({
      langList: langList,
    });
    return langList;
  } catch (e) {
    console.log(e);
    const langList = [
      {
        hasUsed: true,
        id: 1,
        languageCode: 'zh-CN',
        languageName: '中文（简体）',
        languageType: 1,
        label: '中文（简体）',
        value: 'zh-CN',
      },
      {
        hasUsed: true,
        id: 2,
        languageCode: 'en-US',
        languageName: 'English',
        languageType: 1,
        label: 'English',
        value: 'en-US',
      },
    ];
    useI18nStore.setState({
      langList: langList,
    });
    return langList;
  }
};
/**
 * @description: 系统基础store 路由、用户信息、系统信息、当前菜单信息等
 *
 * currentUserInfo: 用户相关信息，包含：用户角色，用户存在的操作权限buttonList,拒绝优先操作资源组denyButtonGroup，操作资源组buttonGroup 等；
 * 导航路由信息-根据权限显示：menuTree, menuGroup
 * home页路由信息-根据权限显示：homeTree, homeGroup
 * home页Tab: homeTabGroup（暂未控制权限）
 * 原始菜单组（导航的不控制权限 含父级目录）: originMenu
 * 所有按钮组（不控制权限）: allButtonGroup
 * **/
export const initBaseContent = {
  originMenu: [],
  menuTree: [],
  homeTree: [],
  menuGroup: [],
  homeGroup: [],
  homeTabGroup: [],
  allButtonGroup: [],
  currentUserInfo: {},
  systemInfo: { appTitle: '' },
  dataBaseType: [],
  dashboardType: [],
  userTipsEnable: storageOpt.getOrigin(SUPOS_USER_TIPS_ENABLE) || '',
  pluginList: [],
  buttonList: [],
  loading: true,
};

export const useBaseStore: UseBoundStoreWithEqualityFn<StoreApi<TBaseStore>> = createWithEqualityFn(
  () => initBaseContent,
  shallow
);

// 设置用户tipsEnable
export const setUserTipsEnable = (value: string) => {
  storageOpt.setOrigin(SUPOS_USER_TIPS_ENABLE, value);
  useBaseStore.setState({
    userTipsEnable: value,
  });
};

const criteria: Criteria<DataItem> = {
  buttonGroup: (item: any) => item?.uri?.includes('button:'),
};

// edge版本 用户默认支持所有的权限和菜单
// 更新路由基础方法 (私有)
const updateBaseStore = async (isFirst: boolean = false) => {
  if (isFirst) {
    try {
      // 首次需要同时拿到用户信息的url和路由
      const [{ value: resource, reason }, { value: info }, { value: systemInfo }]: any = await Promise.allSettled([
        getRoutesResourceApi(),
        getUserInfo(),
        getSystemConfig(),
      ]);

      // 国际化语言包list
      await getLangList();

      // 通过用户的资源池  拿到 - 菜单资源 和 操作资源
      const { buttonGroup, others } = multiGroupByCondition(info?.resourceList, criteria);
      // 拿到 拒绝优先的 菜单资源、 操作资源
      const { buttonGroup: denyButtonGroup, others: denyOthers } = multiGroupByCondition(
        info?.denyResourceList,
        criteria
      );
      // 整合出用户路由资源组
      const userRoutesResourceList = filterObjectArrays(denyOthers, others);
      // 过滤后的路由组,含home\home_tab\menu及目录 去除操作权限
      const allRoutes = filterRouteByUserResource(
        mapResource(resource?.filter((r: ResourceProps) => r.type !== 3)),
        userRoutesResourceList,
        systemInfo?.authEnable && !info?.superAdmin
      );
      // 剔除未启用的路由
      const enableRoutes = allRoutes?.filter((f) => f.enable);
      // 获取终极菜单
      const { homeTree, homeTabGroup, homeGroup, menuGroup, menuTree } = buildResourceTrees(enableRoutes);
      const allButtonGroup = resource?.filter((r: ResourceProps) => r.type === 3);
      const _buttonList =
        systemInfo?.authEnable === false || info?.superAdmin === true
          ? handleButtonPermissions(['button:*'], allButtonGroup) || []
          : filterArrays(
              handleButtonPermissions(denyButtonGroup?.map((i: any) => i.uri) || [], allButtonGroup) || [],
              handleButtonPermissions(buttonGroup?.map((i: any) => i.uri) || [], allButtonGroup) || []
            ) || [];
      // 储存用户信息
      storageOpt.set('personInfo', {
        username: info?.preferredUsername,
      });
      const containerList = filterContainerList(systemInfo?.containerMap || {});
      // 个人用户设置
      const currentUserInfo = {
        ...info,
        roleList: info?.roleList || [],
        roleString: info?.roleList?.map((i: any) => i.roleName)?.join('/') || '',
        buttonList: buttonGroup?.map((i: any) => i.uri) || [],
        pageList: userRoutesResourceList || [],
        superAdmin: info?.superAdmin,
        denyButtonGroup,
        buttonGroup,
      };
      useBaseStore.setState({
        ...initBaseContent,
        homeTree,
        homeTabGroup,
        homeGroup,
        menuGroup,
        menuTree,
        originMenu: resource,
        allButtonGroup,
        // pluginList,
        routesStatus: reason?.status,
        currentUserInfo,
        systemInfo: {
          ...(systemInfo ?? {}),
          appTitle: systemInfo?.appTitle || APP_TITLE,
        },
        containerList,
        buttonList: _buttonList,
        dataBaseType: systemInfo?.containerMap?.tdengine?.envMap?.service_is_show ? ['tdengine'] : ['timescale'],
        mqttBrokeType: systemInfo?.containerMap?.emqx?.name,
        dashboardType:
          containerList.aboutUs
            ?.filter((i) => ['fuxa', 'grafana'].includes(i.name) && i.envMap?.service_is_show)
            ?.map((m) => m.name) ?? [],
      });
      // 设置新手引导
      guideConfig({ systemInfo, menuGroup, info });

      // 设置unsTree信息
      const unsTreeInfo = storageOpt.get(SUPOS_UNS_TREE);
      if (unsTreeInfo) {
        storageOpt.set(SUPOS_UNS_TREE, { ...unsTreeInfo, state: { lazyTree: systemInfo?.lazyTree } });
      } else {
        storageOpt.set(SUPOS_UNS_TREE, { state: { lazyTree: systemInfo?.lazyTree }, version: 0 });
      }
      // 请求国际化语言
      const _lang = await fetchUserLanguage({
        userId: currentUserInfo?.sub,
        lang: systemInfo?.lang,
      });
      // 首次需要初始化语言包
      return await initI18n(_lang);
    } catch (_) {
      console.log(_);
      // 首次需要初始化语言包
      return await initI18n(storageOpt.getOrigin(SUPOS_LANG) || defaultLanguage);
    }
  } else {
    const baseState = useBaseStore.getState();
    // 重新获取菜单
    return getRoutesResourceApi().then((resource: ResourceProps[]) => {
      const allRoutes = filterRouteByUserResource(
        mapResource(resource.filter((r: ResourceProps) => r.type !== 3)),
        baseState?.currentUserInfo?.pageList,
        baseState.systemInfo?.authEnable && !baseState.currentUserInfo?.superAdmin
      );
      const enableRoutes = allRoutes?.filter((f) => f.enable);
      const { homeTree, homeTabGroup, homeGroup, menuGroup, menuTree } = buildResourceTrees(enableRoutes);
      const allButtonGroup = resource?.filter((r: ResourceProps) => r.type === 3);
      const _buttonList =
        baseState?.systemInfo?.authEnable === false || baseState?.currentUserInfo?.superAdmin === true
          ? handleButtonPermissions(['button:*'], allButtonGroup) || []
          : filterArrays(
              handleButtonPermissions(
                baseState?.currentUserInfo?.denyButtonGroup?.map((i: any) => i.uri) || [],
                allButtonGroup
              ) || [],
              handleButtonPermissions(
                baseState?.currentUserInfo?.buttonGroup?.map((i: any) => i.uri) || [],
                allButtonGroup
              ) || []
            ) || [];
      useBaseStore.setState({
        homeTree,
        homeTabGroup,
        homeGroup,
        menuGroup,
        menuTree,
        originMenu: resource,
        allButtonGroup,
        buttonList: _buttonList,
      });
      return allRoutes;
    });
  }
};

// 初始化获取baseStore
export const fetchBaseStore = async (isFirst: boolean = false): Promise<any> => {
  return updateBaseStore(isFirst).finally(() => {
    useBaseStore.setState({
      loading: false,
    });
  });
};

// 设置当前菜单信息
export const setCurrentMenuInfo = (data: ResourceProps) => {
  useBaseStore.setState({
    currentMenuInfo: data,
  });
};

// 手动更新用户信息
export const updateForUserInfo = (info: UserInfoProps) => {
  useBaseStore.setState({
    currentUserInfo: {
      ...useBaseStore.getState().currentUserInfo,
      ...info,
    },
  });
};

export const setPluginList = (pluginList: any[]) => {
  useBaseStore.setState({
    pluginList,
  });
};

const fetchUserLanguage = async (info: { userId?: string; lang?: string }) => {
  const { lang, userId } = info;
  try {
    if (!userId) {
      return import.meta.env.REACT_APP_LOCAL_LANG || lang || storageOpt.getOrigin(SUPOS_LANG) || defaultLanguage;
    } else {
      const response = await getPersonConfigApi(userId);
      return import.meta.env.REACT_APP_LOCAL_LANG || response.mainLanguage;
    }
  } catch (error) {
    console.error('配置请求失败', error);
    return import.meta.env.REACT_APP_LOCAL_LANG || lang || storageOpt.getOrigin(SUPOS_LANG) || defaultLanguage;
  }
};

export const fetchSystemInfo = async (fetchRoute?: boolean): Promise<any> => {
  await getSystemConfig().then((systemInfo) => {
    const containerList = filterContainerList(systemInfo?.containerMap || {});
    useBaseStore.setState({
      systemInfo: {
        ...(systemInfo ?? {}),
        appTitle: systemInfo?.appTitle || APP_TITLE,
      },
      containerList,
      dataBaseType: systemInfo?.containerMap?.tdengine?.envMap?.service_is_show ? ['tdengine'] : ['timescale'],
      mqttBrokeType: systemInfo?.containerMap?.emqx?.name,
      dashboardType:
        containerList.aboutUs
          ?.filter((i) => ['fuxa', 'grafana'].includes(i.name) && i.envMap?.service_is_show)
          ?.map((m) => m.name) ?? [],
    });
    if (fetchRoute) {
      fetchBaseStore?.();
    }
  });
};
