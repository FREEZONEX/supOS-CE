import { ApiWrapper } from '@/utils/request';

const baseUrl = '/copilotkit/mcp';

const api = new ApiWrapper(baseUrl);

// 获取mcp-client列表
export const getMcpList = async () => api.get('/list');
// 删除
export const deleteMcp = async (data: any) => api.post('/delete', data);
// 新增
export const addMcp = async (data: any) => api.post('/add', data);
