import { ApiWrapper } from '@/utils/request';

const baseUrl = '/inter-api/supos/example';

const api = new ApiWrapper(baseUrl);

// 获取example列表
export const queryExamples = async () =>
  api.post('/pageList', {
    pageNo: 1,
    pageSize: 110,
  });

// 安装example
export const installExample = async (id: string) => api.post(`/install?id=${id}`);
// 卸载example
export const unInstallExample = async (id: string) => api.delete(`/uninstall?id=${id}`);
