const localHostnames = new Set(['localhost', '127.0.0.1', '0.0.0.0', '::1']);

interface NormalizeLocalAppUrlOptions {
  allowPrivateHostRewrite?: boolean;
}

const isPrivateIPv4 = (hostname: string) => {
  const parts = hostname.split('.').map((part) => Number(part));
  if (parts.length !== 4 || parts.some((part) => !Number.isInteger(part) || part < 0 || part > 255)) {
    return false;
  }
  const [first, second] = parts;
  return (
    first === 10 ||
    first === 127 ||
    (first === 192 && second === 168) ||
    (first === 172 && second >= 16 && second <= 31)
  );
};

const shouldUseCurrentLocalHost = (targetHostname: string) => {
  if (typeof window === 'undefined') {
    return false;
  }
  const currentHostname = window.location.hostname;
  if (!localHostnames.has(currentHostname) || currentHostname === targetHostname) {
    return false;
  }
  return localHostnames.has(targetHostname);
};

export const normalizeLocalAppUrl = (rawUrl: string, options: NormalizeLocalAppUrlOptions = {}) => {
  if (!rawUrl || !/^https?:\/\//i.test(rawUrl)) {
    return rawUrl;
  }
  try {
    const url = new URL(rawUrl);
    const shouldRewriteHost =
      shouldUseCurrentLocalHost(url.hostname) ||
      (options.allowPrivateHostRewrite &&
        typeof window !== 'undefined' &&
        localHostnames.has(window.location.hostname) &&
        isPrivateIPv4(url.hostname));
    if (shouldRewriteHost) {
      url.hostname = window.location.hostname;
      return url.toString();
    }
  } catch {
    return rawUrl;
  }
  return rawUrl;
};
