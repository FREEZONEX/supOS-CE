import { ApiWrapper } from '@/utils/request';

const baseUrl = '/inter-api/supos/plugin';

const api = new ApiWrapper(baseUrl);

// 列表
export const getPluginListApi = async () => api.get('');
// 安装
export const installPluginApi = async (data: any) => api.post('/install', data);
// 未安装
export const unInstallPluginApi = async (params: any) => api.delete('/uninstall', { params });
// 升级
export const upgradePluginApi = async (data: any) =>
  api.uploads(`/upgrade`, data, {
    method: 'post',
  });
