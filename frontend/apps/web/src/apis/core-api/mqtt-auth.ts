import { coreApi, pageData } from './core-adapter';

export type MqttAuthClient = {
  id: number;
  name?: string;
  clientId: string;
  username: string;
  password: string;
  type: 'device' | 'flow' | 'edge' | 'system';
  typeCode?: number;
  description?: string;
  authStatus?: number;
  connectStatus?: number;
  ip?: string;
  connectedTime?: number;
  clientIdRandomSuffixEnabled?: boolean;
  createdTime?: number;
  updatedTime?: number;
};

export const listMqttAuthClients = async (params?: { type?: string; keyword?: string }) => {
  const resp = await coreApi.get('/mqtt/auth-clients', { params });
  return pageData((resp?.list || resp?.data || []) as MqttAuthClient[]);
};

export const createMqttAuthClient = (payload: {
  name?: string;
  type?: 'device' | 'system';
  description?: string;
  clientIdRandomSuffixEnabled?: boolean;
}) => coreApi.post('/mqtt/auth-clients', payload);

export const resetMqttAuthPassword = (id: number | string) => coreApi.post(`/mqtt/auth-clients/${id}/reset-password`);

export const enableMqttAuthClient = (id: number | string) => coreApi.post(`/mqtt/auth-clients/${id}/enable`);

export const disableMqttAuthClient = (id: number | string) => coreApi.post(`/mqtt/auth-clients/${id}/disable`);

export const deleteMqttAuthClient = (id: number | string) => coreApi.delete(`/mqtt/auth-clients/${id}`);
