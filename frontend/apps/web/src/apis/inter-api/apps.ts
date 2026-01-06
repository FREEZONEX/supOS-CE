import { ApiWrapper } from '@/utils/request';

const baseUrl = '/inter-api/supos';

const api = new ApiWrapper(baseUrl);

// 获取所有APP列表
export const getApps = async (params?: Record<string, unknown>) => api.get('/apps', { params });
// 获取单个APP列表
export const getSingleApp = async (name?: string) => api.get(`/app/${name}`);
// 获取单个html内容
export const getSingleHtml = async (name: string, htmlName?: string) => api.get(`/app/${name}/html/${htmlName}`);
// 创建App
export const createApp = async (data: any) => api.post(`/app/create`, data);
// 上传html "file":   二进制数据
export const uploadHtml = async (name: string, data: any) => api.upload(`/app/${name}`, data);
// 删除App
export const destroyApp = async (name: string) => api.delete(`/app/${name}/destroy`);
// 删除html
export const destroyHtml = async (name?: string, id?: string) => api.delete(`/app/${name}/html/${id}`);
// 设置主页
export const setHomepage = async (data: any) => api.post(`/app/homepage`, data);
