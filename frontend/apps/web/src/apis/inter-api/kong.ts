import service from '@/utils/request';

const KONG_BASE = '/kong-admin';

const kong = {
  get: (url: string, params?: Record<string, unknown>) =>
    service.request({ url: KONG_BASE + url, method: 'get', params, _noCode: true } as any),
  post: (url: string, data?: Record<string, unknown>) =>
    service.request({ url: KONG_BASE + url, method: 'post', data, _noCode: true } as any),
  patch: (url: string, data?: Record<string, unknown>) =>
    service.request({ url: KONG_BASE + url, method: 'patch', data, _noCode: true } as any),
  del: (url: string) => service.request({ url: KONG_BASE + url, method: 'delete', _noCode: true } as any),
};

// ─── Services ────────────────────────────────────────────────────────
export const getServices = (params?: Record<string, unknown>) => kong.get('/services', params);

export const getService = (idOrName: string) => kong.get(`/services/${idOrName}`);

export const createService = (data: Record<string, unknown>) => kong.post('/services', data);

export const updateService = (id: string, data: Record<string, unknown>) => kong.patch(`/services/${id}`, data);

export const deleteService = (id: string) => kong.del(`/services/${id}`);

// ─── Routes ──────────────────────────────────────────────────────────
export const getRoutes = (params?: Record<string, unknown>) => kong.get('/routes', params);

export const getRoute = (id: string) => kong.get(`/routes/${id}`);

export const createRoute = (data: Record<string, unknown>) => kong.post('/routes', data);

export const updateRoute = (id: string, data: Record<string, unknown>) => kong.patch(`/routes/${id}`, data);

export const deleteRoute = (id: string) => kong.del(`/routes/${id}`);

// ─── Consumers ───────────────────────────────────────────────────────
export const getConsumers = (params?: Record<string, unknown>) => kong.get('/consumers', params);

export const getConsumer = (idOrUsername: string) => kong.get(`/consumers/${idOrUsername}`);

export const createConsumer = (data: Record<string, unknown>) => kong.post('/consumers', data);

export const updateConsumer = (id: string, data: Record<string, unknown>) => kong.patch(`/consumers/${id}`, data);

export const deleteConsumer = (id: string) => kong.del(`/consumers/${id}`);

// ─── Plugins ─────────────────────────────────────────────────────────
export const getPlugins = (params?: Record<string, unknown>) => kong.get('/plugins', params);

export const getPlugin = (id: string) => kong.get(`/plugins/${id}`);

export const getEnabledPlugins = () => kong.get('/plugins/enabled');

export const getPluginSchema = (name: string) => kong.get(`/plugins/schema/${name}`);

export const createPlugin = (data: Record<string, unknown>) => kong.post('/plugins', data);

export const updatePlugin = (id: string, data: Record<string, unknown>) => kong.patch(`/plugins/${id}`, data);

export const deletePlugin = (id: string) => kong.del(`/plugins/${id}`);

// ─── Certificates ────────────────────────────────────────────────────
export const getCertificates = (params?: Record<string, unknown>) => kong.get('/certificates', params);

export const getCertificate = (id: string) => kong.get(`/certificates/${id}`);

export const createCertificate = (data: Record<string, unknown>) => kong.post('/certificates', data);

export const updateCertificate = (id: string, data: Record<string, unknown>) => kong.patch(`/certificates/${id}`, data);

export const deleteCertificate = (id: string) => kong.del(`/certificates/${id}`);

// ─── Service Routes ───────────────────────────────────────────────────
export const getServiceRoutes = (serviceId: string, params?: Record<string, unknown>) =>
  kong.get(`/services/${serviceId}/routes`, params);

// ─── Consumer Credentials - Basic Auth ───────────────────────────────
export const getConsumerBasicAuth = (consumerId: string) => kong.get(`/consumers/${consumerId}/basic-auth`);
export const createConsumerBasicAuth = (consumerId: string, data: Record<string, unknown>) =>
  kong.post(`/consumers/${consumerId}/basic-auth`, data);
export const deleteConsumerBasicAuth = (consumerId: string, credId: string) =>
  kong.del(`/consumers/${consumerId}/basic-auth/${credId}`);

// ─── Consumer Credentials - API Keys ─────────────────────────────────
export const getConsumerKeyAuth = (consumerId: string) => kong.get(`/consumers/${consumerId}/key-auth`);
export const createConsumerKeyAuth = (consumerId: string, data: Record<string, unknown>) =>
  kong.post(`/consumers/${consumerId}/key-auth`, data);
export const deleteConsumerKeyAuth = (consumerId: string, credId: string) =>
  kong.del(`/consumers/${consumerId}/key-auth/${credId}`);

// ─── Consumer Credentials - HMAC ─────────────────────────────────────
export const getConsumerHmacAuth = (consumerId: string) => kong.get(`/consumers/${consumerId}/hmac-auth`);
export const createConsumerHmacAuth = (consumerId: string, data: Record<string, unknown>) =>
  kong.post(`/consumers/${consumerId}/hmac-auth`, data);
export const deleteConsumerHmacAuth = (consumerId: string, credId: string) =>
  kong.del(`/consumers/${consumerId}/hmac-auth/${credId}`);

// ─── Consumer Credentials - OAuth2 ───────────────────────────────────
export const getConsumerOAuth2 = (consumerId: string) => kong.get(`/consumers/${consumerId}/oauth2`);
export const createConsumerOAuth2 = (consumerId: string, data: Record<string, unknown>) =>
  kong.post(`/consumers/${consumerId}/oauth2`, data);
export const deleteConsumerOAuth2 = (consumerId: string, credId: string) =>
  kong.del(`/consumers/${consumerId}/oauth2/${credId}`);

// ─── Consumer Credentials - JWT ──────────────────────────────────────
export const getConsumerJwt = (consumerId: string) => kong.get(`/consumers/${consumerId}/jwt`);
export const createConsumerJwt = (consumerId: string, data: Record<string, unknown>) =>
  kong.post(`/consumers/${consumerId}/jwt`, data);
export const deleteConsumerJwt = (consumerId: string, credId: string) =>
  kong.del(`/consumers/${consumerId}/jwt/${credId}`);

// ─── Compatibility alias (used by menu-configuration) ────────────────
export const getKongRoutesApi = async (key?: string) => {
  const res: any = await getRoutes({ size: 1000 });
  const routes = res?.data ?? [];
  if (!key) {
    return routes;
  }

  const keyword = key.toLowerCase();
  return routes.filter((item: any) =>
    String(item?.name ?? '')
      .toLowerCase()
      .includes(keyword)
  );
};

// ─── Node info ───────────────────────────────────────────────────────
export const getKongInfo = () => kong.get('/');

export const getKongStatus = () => kong.get('/status');
