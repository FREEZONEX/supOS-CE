import {
  authApi,
  buttonUrisForResourceKeys,
  coreResourcesForResourceKeys,
  routeUrisForResourceKeys,
  type CoreResource,
} from './core-adapter';
import { removeToken, setToken } from '@/utils/auth';

const isSuperAdmin = (user?: any) => {
  const roleCode = String(user?.roleCode || '')
    .trim()
    .toLowerCase();
  if (roleCode === 'admin' || roleCode === 'builder') {
    return true;
  }
  return ['tier0', 'admin'].includes(String(user?.userName || '').trim());
};
const defaultHomePage = '/uns';
const defaultProfileEmail = 'tier0@example.com';
const validEmailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

const normalizeProfileEmail = (email?: string, user?: any) => {
  const value = String(email || '').trim();
  if (validEmailPattern.test(value)) return value;
  return isSuperAdmin(user) && String(user?.userName || '').trim() === 'tier0' ? defaultProfileEmail : value;
};

const normalizeHomePage = (homePage?: string | null) => {
  const value = String(homePage || '').trim();
  return value.startsWith('/') && !value.startsWith('//') ? value : '';
};

const routeUri = (item: any) => String(item?.uri || '').trim();
const isLaunchpadAppHomePage = (value: string) => /^\/launchpad\/[^/?#]+\/[^/?#]+\/?$/i.test(value);

export const resolveAllowedHomePage = (homePage?: string | null, resourceList: any[] = [], superAdmin = false) => {
  const normalized = normalizeHomePage(homePage);
  if (superAdmin) {
    return normalized || defaultHomePage;
  }
  const allowedRoutes = resourceList.map(routeUri).filter((uri) => uri.startsWith('/'));
  if (!allowedRoutes.length) {
    return defaultHomePage;
  }
  const permissionRoute = isLaunchpadAppHomePage(normalized) ? '/launchpad' : normalized;
  if (permissionRoute && allowedRoutes.some((uri) => uri.toLowerCase() === permissionRoute.toLowerCase())) {
    return normalized;
  }
  return allowedRoutes.find((uri) => uri.toLowerCase() === defaultHomePage) || allowedRoutes[0] || defaultHomePage;
};

const getCurrentUserConfig = async () => {
  try {
    return await authApi.get('/config');
  } catch {
    return { homePage: defaultHomePage };
  }
};
const getPermissionResources = async (user: any = {}) => {
  return coreResourcesForResourceKeys(user?.resourceKeys || []);
}

const buildUserPermissions = (resources: CoreResource[] = [], user?: any, config?: any) => {
  const resourceList = routeUrisForResourceKeys(resources, user?.resourceKeys || []);
  const uiButtons = (user?.uiActions || []).filter(
    (item: string) => item === 'button:*' || String(item || '').startsWith('button:')
  );
  const buttonList = uiButtons.length > 0 ? uiButtons : buttonUrisForResourceKeys(resources, user?.resourceKeys || []);
  const superAdmin = isSuperAdmin(user);

  return {
    resourceList,
    buttonList: superAdmin ? ['button:*'] : buttonList,
    homePage: resolveAllowedHomePage(config?.homePage, resourceList, superAdmin),
    superAdmin,
  };
};

export const loginApi = async (data?: Record<string, unknown>) => {
  const result = await authApi.post('/login', data);
  if (result?.token) {
    setToken(result.token, { expires: 0.5, path: '/' });
  }
  const [resources, config] = await Promise.all([getPermissionResources(result), getCurrentUserConfig()]);
  const permissions = buildUserPermissions(resources as CoreResource[], result, config);
  return {
    ...result,
    sub: String(result?.userId || ''),
    preferredUsername: result?.userName,
    firstName: result?.nickName || result?.userName,
    forceChangePassword: Boolean(result?.forceChangePassword),
    phone: result?.phone,
    email: normalizeProfileEmail(result?.email, result),
    ...permissions,
  };
};
// 获取用户信息
export const getUserInfo = async () => {
  const user = await authApi.get('/me');
  const [resources, config] = await Promise.all([getPermissionResources(user), getCurrentUserConfig()]);
  const coreResources = resources as CoreResource[];
  const permissions = buildUserPermissions(coreResources, user, config);
  return {
    firstTimeLogin: user?.isRandomPwd ? 1 : 0,
    forceChangePassword: Boolean(user?.forceChangePassword),
    tipsEnable: 0,
    sub: String(user?.userId || ''),
    preferredUsername: user?.userName,
    firstName: user?.nickName || user?.userName,
    email: normalizeProfileEmail(user?.email, user),
    phone: user?.phone,
    enabled: true,
    roleList: [],
    resourceKeys: user?.resourceKeys || [],
    ...permissions,
    denyResourceList: [],
  };
};

export const logoutApi = async () => {
  try {
    await authApi.post('/logout');
  } finally {
    removeToken();
  }
};
