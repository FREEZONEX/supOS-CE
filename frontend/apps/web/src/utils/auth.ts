import Cookies from 'js-cookie';
import { useBaseStore } from '@/stores/base';
import { APP_COMMUNITY_TOKEN, LOGIN_URL, SYS_COMMUNITY_TOKEN } from '@/common-types/constans.ts';

const isBrowser = () => typeof window !== 'undefined';

export const isDevScopedTokenEnabled = () => {
  if (!import.meta.env.DEV || !isBrowser()) {
    return false;
  }
  const hostname = window.location.hostname;
  return (
    Boolean(window.location.port) || hostname === 'localhost' || hostname === '127.0.0.1' || hostname === '0.0.0.0'
  );
};

const devScopedTokenKey = () => `${APP_COMMUNITY_TOKEN}:${window.location.host}`;

const getDevScopedToken = () => {
  if (!isDevScopedTokenEnabled()) {
    return undefined;
  }
  try {
    return window.sessionStorage.getItem(devScopedTokenKey()) || undefined;
  } catch {
    return undefined;
  }
};

const setDevScopedToken = (token: string) => {
  if (!isDevScopedTokenEnabled()) {
    return;
  }
  try {
    window.sessionStorage.setItem(devScopedTokenKey(), token);
  } catch {
    // sessionStorage can be unavailable in strict browser privacy modes.
  }
};

const removeDevScopedToken = () => {
  if (!isDevScopedTokenEnabled()) {
    return;
  }
  try {
    window.sessionStorage.removeItem(devScopedTokenKey());
  } catch {
    // ignore storage cleanup failures
  }
};

export function getToken() {
  if (isDevScopedTokenEnabled()) {
    return getDevScopedToken();
  }
  return Cookies.get(APP_COMMUNITY_TOKEN) || Cookies.get(SYS_COMMUNITY_TOKEN);
}

export function setToken(token: string, options?: any) {
  setDevScopedToken(token);
  Cookies.remove(SYS_COMMUNITY_TOKEN, { path: '/' });
  return Cookies.set(APP_COMMUNITY_TOKEN, token, { path: '/', ...(options || {}) });
}

const tokenCookiePaths = () => {
  const paths = new Set(['/', LOGIN_URL, '/api', '/api/core', '/api/core/auth']);
  if (typeof window !== 'undefined') {
    const parts = window.location.pathname.split('/').filter(Boolean);
    let current = '';
    parts.forEach((part) => {
      current += `/${part}`;
      paths.add(current);
    });
  }
  return Array.from(paths);
};

const tokenCookieDomains = () => {
  if (typeof window === 'undefined') {
    return [undefined];
  }
  const hostname = window.location.hostname;
  const domains: Array<string | undefined> = [undefined, hostname];
  if (hostname && !hostname.includes(':') && !/^\d+\.\d+\.\d+\.\d+$/.test(hostname)) {
    domains.push(`.${hostname}`);
  }
  return domains;
};

const removeTokenCookie = (name: string) => {
  Cookies.remove(name);
  tokenCookiePaths().forEach((path) => {
    tokenCookieDomains().forEach((domain) => {
      const options = domain ? { path, domain } : { path };
      Cookies.remove(name, options);
    });
  });
};

export function removeToken() {
  removeDevScopedToken();
  removeTokenCookie(SYS_COMMUNITY_TOKEN);
  removeTokenCookie(APP_COMMUNITY_TOKEN);
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
  const state = useBaseStore.getState();
  if (state.systemInfo?.authEnable === false || state.currentUserInfo?.superAdmin === true) {
    return true;
  }
  // button: 必须要有
  const perms = state.buttonList;
  // 重构期存在后端资源树已下发、按钮资源还未完全补齐的情况；没有按钮权限列表时不误隐藏页面操作。
  if (!perms?.length) {
    return true;
  }
  if (perms?.includes('button:*')) {
    return true;
  }
  if (auth instanceof Array) {
    return auth?.some((item) => perms?.includes(item.startsWith('button:') ? item : 'button:' + item));
  } else {
    return perms?.includes(auth.startsWith('button:') ? auth : 'button:' + auth);
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
