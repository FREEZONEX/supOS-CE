import { ApiWrapper } from '@/utils/request';

const baseUrl = '/inter-api/supos';

const api = new ApiWrapper(baseUrl);

// 系统配置
export const getSystemConfig = async () => api.get('/systemConfig');

//获取主题配置

export const getAllThemeConfig = async (params?: any) => api.get('/theme/getConfig', { params });
