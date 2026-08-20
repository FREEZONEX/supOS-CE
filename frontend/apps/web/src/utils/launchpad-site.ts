const isBrowser = () => typeof window !== 'undefined';

const normalizePort = (value?: string) => String(value || '').trim();

const getRuntimeLaunchpadPort = () => {
  if (!isBrowser()) {
    return '';
  }
  return normalizePort(window.__TIER0_RUNTIME_CONFIG__?.launchpadPort);
};

const isRuntimeLaunchpadStandalone = () => Boolean(isBrowser() && window.__TIER0_RUNTIME_CONFIG__?.launchpadStandalone);

const getProcessLaunchpadPort = () => {
  if (typeof process === 'undefined') {
    return '';
  }
  return normalizePort(process.env?.LAUNCHPAD_PORT);
};

export const getLaunchpadStandalonePort = () =>
  getRuntimeLaunchpadPort() || normalizePort(import.meta.env.VITE_LAUNCHPAD_PORT) || getProcessLaunchpadPort();

export const isLaunchpadRoutePath = (pathname: string) =>
  pathname === '/launchpad' || pathname.startsWith('/launchpad/');

export const isLaunchpadUserSettingsPath = (pathname: string) =>
  pathname === '/settings' ||
  pathname === '/settings/profile' ||
  pathname === '/settings/preferences' ||
  pathname === '/settings/security';

export const isLaunchpadStandalonePort = () => {
  if (isRuntimeLaunchpadStandalone()) {
    return true;
  }
  const standalonePort = getLaunchpadStandalonePort();
  return Boolean(standalonePort && isBrowser() && window.location.port === standalonePort);
};

export const isLaunchpadStandaloneAllowedPath = (pathname: string) =>
  isLaunchpadRoutePath(pathname) || (isLaunchpadStandalonePort() && isLaunchpadUserSettingsPath(pathname));

export const isLaunchpadSiteRequest = (pathname: string) =>
  isLaunchpadStandalonePort() && (isLaunchpadRoutePath(pathname) || isLaunchpadUserSettingsPath(pathname));
