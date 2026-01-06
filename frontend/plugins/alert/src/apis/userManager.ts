import { ApiWrapper } from '@supos_host/utils';

const baseUrl = '/inter-api/supos/userManage';

const api = new ApiWrapper(baseUrl);

// 获取用户信息 - select使用
export const searchUserManageList = async (data?: Record<string, unknown>) =>
  api.post('/pageList', data).then((data: any) => {
    return data?.map?.((item: any) => ({
      label: item.preferredUsername,
      value: item.id,
    }));
  });
