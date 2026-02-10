import { ApiWrapper } from '@/utils/request';

const baseUrl = '/inter-api/supos/group';

const api = new ApiWrapper(baseUrl);

// dashboard eventflow sourflow 的分组接口
// type 1-sourceflow 2-eventflow 3-datasource

// 新增分组  Name Description type
export const addGroup = async (data: any) => api.post('', data);

// 编辑分组
export const editGroup = async (data: any) => api.put('', data);

// 删除分组  组id
export const deleteGroup = async (uid: string) => api.delete(`/${uid}`);

// 置顶 id: 组id status： 是否置顶
export const markGroup = async (id: string, status: boolean) => api.post('/operationGroupTop', { id, status });

// 移入/移出组  bizId: 文件id   id: 组id
export const optGroup = async ({ bizId, id, status }: { bizId: string; id: string; status: boolean }) =>
  api.post('/operationGroup', { id, status, bizId });

// 获取分组列表 type 1-sourceflow 2-eventflow 3-datasource ,  page , size
export const getGroupList = async (params?: Record<string, unknown>) =>
  api.get('/by-type', {
    params,
  });
