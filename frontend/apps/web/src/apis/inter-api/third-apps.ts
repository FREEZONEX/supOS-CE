import { ApiWrapper, CustomAxiosConfigEnum } from '@/utils';

const baseUrl = '/inter-api/supos/third-apps';

const api = new ApiWrapper(baseUrl);

// interface PackageInterface {
//   packageId: string;
// }

interface AppInterface {
  appId: string;
  properties?: string;
}

// app列表
export const pageListApi = async (data?: Record<string, unknown>) =>
  api.post('/pageList', data, {
    [CustomAxiosConfigEnum.BusinessResponse]: true,
  });

// 卸载app
export const uninstallAppApi = async (data: AppInterface) => api.post('/uninstallApp', undefined, { params: data });

// 暂停app
export const stopAppApi = async (data: AppInterface) => api.post('/stopApp', undefined, { params: data });

// 开启app
export const startAppApi = async (data: AppInterface) => api.post('/startApp', undefined, { params: data });

// 安装app
export const installAppApi = async (data: AppInterface) => api.post('/install', undefined, { params: data });

// 更新配置文件
export const updateConfigApi = async (data: AppInterface) => api.put('/properties', data);

//刷新app状态
export const refreshAppStatus = async (id: string) => api.get(`/refresh?appId=${id}`);

//上传app
export const uploadApp = async (data: any) =>
  api.uploads('/upload', data, { method: 'post', [CustomAxiosConfigEnum.BusinessResponse]: true });

//删除app
export const deleteApp = async (params: any) => api.delete('/deleteApp', { params });

//批量启动
export const batchStartApp = async (params: any) => api.post('/batchStartApp', params);

//批量暂停
export const batchStopApp = async (params: any) => api.post('/batchStopApp', params);

//是否解压完成
export const existApp = async (appId: any) => api.get(`/existApp?appId=${appId}`);
