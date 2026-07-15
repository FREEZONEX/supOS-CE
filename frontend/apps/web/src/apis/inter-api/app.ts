import { ApiWrapper } from '@/utils/request';

const baseUrl = '/inter-api/supos/app';

const api = new ApiWrapper(baseUrl);

// 安装app
export const installApp = async (data?: Record<string, unknown>) => api.post('/install', data);
