import { ApiWrapper, CustomAxiosConfigEnum } from '@/utils/request';

const baseUrl = '/inter-api/supos/notify';

const api = new ApiWrapper(baseUrl);

export const queryNoticeList = async (data: any) =>
  api.post('/pageList', data, { [CustomAxiosConfigEnum.BusinessResponse]: true }); // 获取所以树数据
export const noticeToRead = async (data: any) => api.put('/read', data); // 获取所以树数据
export const deleteNotice = async (data: any) => api.delete('/delete', {}, data); // 获取数据类型
