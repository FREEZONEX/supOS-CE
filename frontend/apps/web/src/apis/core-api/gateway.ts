import { coreApi } from './core-adapter';

export interface GatewayRoute {
  routeKey: string;
  name: string;
  description?: string;
  targetType: string;
  matchType: string;
  pathPattern: string;
  methods?: string[];
  targetUrl?: string;
  targetPath?: string;
  stripPrefix?: boolean;
  rewritePath?: string;
  authPolicy: string;
  resourceKey?: string;
  timeoutMs?: number;
  priority?: number;
  enabled: boolean;
  systemBuiltin?: boolean;
}

export const listGatewayRoutes = async (): Promise<{ list: GatewayRoute[]; total: number }> => {
  const resp = await coreApi.get('/gateway/routes');
  return {
    list: resp?.list || [],
    total: Number(resp?.total || resp?.list?.length || 0),
  };
};

export const saveGatewayRoute = async (data: GatewayRoute) => {
  const payload = {
    ...data,
    targetType: data.targetType || 'reverseProxy',
    matchType: data.matchType || 'prefix',
    methods: data.methods?.length ? data.methods : ['GET'],
    authPolicy: data.authPolicy || 'login',
    timeoutMs: Number(data.timeoutMs || 10000),
    priority: Number(data.priority || 100),
    enabled: data.enabled !== false,
  };
  return coreApi.post('/gateway/routes', payload);
};

export const updateGatewayRoute = async (routeKey: string, data: GatewayRoute) => {
  const payload = {
    ...data,
    routeKey,
    targetType: data.targetType || 'reverseProxy',
    matchType: data.matchType || 'prefix',
    methods: data.methods?.length ? data.methods : ['GET'],
    authPolicy: data.authPolicy || 'login',
    timeoutMs: Number(data.timeoutMs || 10000),
    priority: Number(data.priority || 100),
    enabled: data.enabled !== false,
  };
  return coreApi.put(`/gateway/routes/${encodeURIComponent(routeKey)}`, payload);
};

export const deleteGatewayRoute = async (routeKey: string) =>
  coreApi.delete(`/gateway/routes/${encodeURIComponent(routeKey)}`);

export const testGatewayRoute = async (routeKey: string) =>
  coreApi.post(`/gateway/routes/${encodeURIComponent(routeKey)}/test`, {});
