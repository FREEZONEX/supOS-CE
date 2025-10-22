import { ApiWrapper, CustomAxiosConfigEnum } from '@supos_host/utils';

const baseUrl = '/inter-api/supos/uns';

const api = new ApiWrapper(baseUrl);
// 报警列表页
export const getAlertList = async (params?: Record<string, unknown>) =>
  api
    .get('/search', {
      params,
      [CustomAxiosConfigEnum.BusinessResponse]: true,
    })
    .then((data: any) => {
      return {
        ...data,
        pageNo: data?.page?.pageNo,
        pageSize: data?.page?.pageSize,
        total: data?.page?.total,
      };
    });

export const getInstanceInfo = async (params?: Record<string, unknown>) =>
  api.get('/instance', { params, _noMessage: true }); // 查询模型字段声明

export const searchTreeData = async (params?: Record<string, unknown>) => api.get('/search', { params }); // 获取所以树数据
