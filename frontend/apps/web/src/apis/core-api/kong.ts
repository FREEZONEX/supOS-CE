/* eslint-disable @typescript-eslint/no-unused-vars */
import { coreApi } from './core-adapter';

const toKongRoute = (item: any) => {
  const id = item.routeKey || item.id;
  const serviceId = item.targetUrl || 'static';
  return {
    id,
    name: item.name || id,
    protocols: ['http', 'https'],
    methods: item.methods || ['GET'],
    paths: [item.pathPattern || '/'],
    strip_path: Boolean(item.stripPrefix),
    preserve_host: false,
    service: { id: serviceId },
    _raw: item,
  };
};

const toGatewayRoute = (data: Record<string, any>, routeKey?: string) => {
  const firstPath = Array.isArray(data.paths) ? data.paths[0] : data.pathPattern || '/';
  const targetUrl = data.targetUrl || data.service?.id || 'http://example.invalid';
  return {
    routeKey: routeKey || data.id || data.name || `route-${Date.now()}`,
    name: data.name || routeKey || `route-${Date.now()}`,
    description: data.description || '',
    targetType: 'reverseProxy',
    matchType: 'prefix',
    pathPattern: firstPath.endsWith('/**') ? firstPath : `${String(firstPath).replace(/\/$/, '')}/**`,
    methods: data.methods || ['GET'],
    targetUrl,
    targetPath: '/',
    stripPrefix: Boolean(data.strip_path ?? data.stripPrefix),
    rewritePath: '',
    authPolicy: data.authPolicy || 'login',
    resourceKey: data.resourceKey || '',
    timeoutMs: Number(data.timeoutMs || 30000),
    priority: Number(data.priority || 100),
    enabled: data.enabled !== false,
  };
};

const routePage = async () => {
  const resp = await coreApi.get('/gateway/routes');
  return { data: (resp?.list || []).map(toKongRoute) };
};

const servicePage = async () => {
  const routes = await routePage();
  const seen = new Map<string, any>();
  routes.data.forEach((route: any) => {
    const id = route.service?.id || 'static';
    if (!seen.has(id)) {
      seen.set(id, {
        id,
        name: id.replace(/^https?:\/\//, ''),
        protocol: id.startsWith('https://') ? 'https' : 'http',
        host: id,
        path: '/',
      });
    }
  });
  return { data: Array.from(seen.values()) };
};

export const getServices = (_params?: Record<string, unknown>) => servicePage();
export const getService = async (idOrName: string) => ({ id: idOrName, name: idOrName });
export const createService = async (_data: Record<string, unknown>) =>
  Promise.reject(new Error('service create is represented by gateway route targetUrl'));
export const updateService = async (_id: string, _data: Record<string, unknown>) => ({});
export const deleteService = async (_id: string) => ({});

export const getRoutes = (_params?: Record<string, unknown>) => routePage();
export const getRoute = async (id: string) => {
  const routes = await routePage();
  return routes.data.find((item: any) => item.id === id);
};
export const createRoute = async (data: Record<string, unknown>) =>
  coreApi.post('/gateway/routes', toGatewayRoute(data));
export const updateRoute = async (id: string, data: Record<string, unknown>) =>
  coreApi.put(`/gateway/routes/${id}`, toGatewayRoute(data, id));
export const deleteRoute = async (id: string) => coreApi.delete(`/gateway/routes/${id}`);

const emptyPage = { data: [] };
export const getConsumers = async (_params?: Record<string, unknown>) => emptyPage;
export const getConsumer = async (_idOrUsername: string) => ({});
export const createConsumer = async (_data: Record<string, unknown>) => ({});
export const updateConsumer = async (_id: string, _data: Record<string, unknown>) => ({});
export const deleteConsumer = async (_id: string) => ({});
export const getPlugins = async (_params?: Record<string, unknown>) => emptyPage;
export const getPlugin = async (_id: string) => ({});
export const getEnabledPlugins = async () => ({ enabled_plugins: ['jwt', 'key-auth', 'cors'] });
export const getPluginSchema = async (_name: string) => ({ fields: [] });
export const createPlugin = async (_data: Record<string, unknown>) => ({});
export const updatePlugin = async (_id: string, _data: Record<string, unknown>) => ({});
export const deletePlugin = async (_id: string) => ({});
export const getCertificates = async (_params?: Record<string, unknown>) => emptyPage;
export const getCertificate = async (_id: string) => ({});
export const createCertificate = async (_data: Record<string, unknown>) => ({});
export const updateCertificate = async (_id: string, _data: Record<string, unknown>) => ({});
export const deleteCertificate = async (_id: string) => ({});
export const getServiceRoutes = async (serviceId: string) => {
  const routes = await routePage();
  return { data: routes.data.filter((item: any) => item.service?.id === serviceId) };
};

export const getConsumerBasicAuth = async (_consumerId: string) => emptyPage;
export const createConsumerBasicAuth = async (_consumerId: string, _data: Record<string, unknown>) => ({});
export const deleteConsumerBasicAuth = async (_consumerId: string, _credId: string) => ({});
export const getConsumerKeyAuth = async (_consumerId: string) => emptyPage;
export const createConsumerKeyAuth = async (_consumerId: string, _data: Record<string, unknown>) => ({});
export const deleteConsumerKeyAuth = async (_consumerId: string, _credId: string) => ({});
export const getConsumerHmacAuth = async (_consumerId: string) => emptyPage;
export const createConsumerHmacAuth = async (_consumerId: string, _data: Record<string, unknown>) => ({});
export const deleteConsumerHmacAuth = async (_consumerId: string, _credId: string) => ({});
export const getConsumerOAuth2 = async (_consumerId: string) => emptyPage;
export const createConsumerOAuth2 = async (_consumerId: string, _data: Record<string, unknown>) => ({});
export const deleteConsumerOAuth2 = async (_consumerId: string, _credId: string) => ({});
export const getConsumerJwt = async (_consumerId: string) => emptyPage;
export const createConsumerJwt = async (_consumerId: string, _data: Record<string, unknown>) => ({});
export const deleteConsumerJwt = async (_consumerId: string, _credId: string) => ({});

export const getKongRoutesApi = async (key?: string) => {
  const res: any = await getRoutes({ size: 1000 });
  const routes = res?.data ?? [];
  if (!key) return routes;
  const keyword = key.toLowerCase();
  return routes.filter((item: any) =>
    String(item?.name ?? '')
      .toLowerCase()
      .includes(keyword)
  );
};

export const getKongInfo = async () => ({ version: 'embedded-gateway' });
export const getKongStatus = async () => ({ database: 'reachable' });
