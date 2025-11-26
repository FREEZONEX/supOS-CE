// 模板实例的文件上传

import { ApiWrapper, CustomAxiosConfigEnum } from '@/utils/request';

const baseUrl = '/inter-api/supos/uns';

const api = new ApiWrapper(baseUrl);

// 获取列表
export const getAttachmentsList = async (params?: Record<string, unknown>) => api.get('/attachments', { params });
// 删除
export const deleteAttachments = async (params?: Record<string, unknown>) => api.delete('/attachment', { params });
// 下载
export const getAttachment = async (params?: Record<string, unknown>) =>
  api.get('/attachment', { params, responseType: 'blob', [CustomAxiosConfigEnum.NoCode]: true });
// 上传
export const uploadAttachment = async (data: any, params: any) =>
  api.uploads(`/attachment`, data, {
    method: 'post',
    params,
  });
