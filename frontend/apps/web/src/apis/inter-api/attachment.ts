// 模板实例的文件上传

import { ApiWrapper } from '@/utils/request';

const baseUrl = '/inter-api/supos/attachment';

const api = new ApiWrapper(baseUrl);

// 上传
export const commonUploadAttachment = async (data: any) =>
  api.uploads(`/upload`, data, {
    method: 'post',
  });
