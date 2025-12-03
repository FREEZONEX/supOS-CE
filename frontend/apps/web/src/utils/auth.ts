import Cookies from 'js-cookie';
import { useBaseStore } from '@/stores/base';
import { SUPOS_COMMUNITY_TOKEN } from '@/common-types/constans.ts';

export function getToken() {
  return Cookies.get(SUPOS_COMMUNITY_TOKEN);
}

export function setToken(token: string, options?: any) {
  return Cookies.set(SUPOS_COMMUNITY_TOKEN, token, options);
}

export function removeToken() {
  return Cookies.remove(SUPOS_COMMUNITY_TOKEN, { path: '/' });
}

export function getCookie(key: any) {
  return Cookies.get(key);
}

export function setCookie(key: any, value: any, expires?: any) {
  return Cookies.set(key, value, { expires });
}

export function removeCookie(key: any) {
  return Cookies.remove(key);
}

export const hasPermission = (auth: string | string[]) => {
  // 开发环境不控制权限
  if (import.meta.env.DEV) {
    return true;
  }
  // button: 必须要有
  const perms = useBaseStore.getState().buttonList;
  if (auth instanceof Array) {
    return auth?.some((item) => perms?.includes('button:' + item));
  } else {
    return perms?.includes('button:' + auth);
  }
};

export const filterPermissionToList = <T, K extends string = 'auth'>(
  list: (T & { [P in K]?: string | string[] })[],
  key: K = 'auth' as K
): T[] => {
  return list
    ?.filter((item) => {
      const itemAny = item as any;
      if (!itemAny[key]) {
        return true;
      } else {
        return hasPermission(itemAny[key]);
      }
    })
    ?.map((item) => {
      const newItem = { ...item } as any;
      delete newItem[key];
      return newItem as T;
    });
};
