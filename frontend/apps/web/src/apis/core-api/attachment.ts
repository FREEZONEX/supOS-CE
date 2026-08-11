// 模板实例的文件上传

import { coreApi } from './core-adapter';

// 上传
export const commonUploadAttachment = async (data: any) =>
  coreApi.uploads(`/assets`, data, {
    method: 'post',
  });
