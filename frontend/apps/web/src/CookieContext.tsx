import { useCookies } from 'react-cookie';
import { LOGIN_URL, OMC_MODEL, SUPOS_COMMUNITY_TOKEN, SUPOS_USER_TIPS_ENABLE } from '@/common-types/constans';
import { useUpdateEffect } from 'ahooks';
import { message } from 'antd';
import { SUPOS_USER_GUIDE_ROUTES } from '@/common-types/constans';
import { storageOpt } from '@/utils/storage';
import { useBaseStore } from '@/stores/base';
import Cookies from 'js-cookie';

// 登录失效控制
const CookieContext = () => {
  const systemInfo = useBaseStore((state) => state.systemInfo);

  const [cookies] = useCookies([SUPOS_COMMUNITY_TOKEN]);

  useUpdateEffect(() => {
    // cookie发生改变删除guide routes信息
    storageOpt.remove(SUPOS_USER_GUIDE_ROUTES);
    // cookie发生改变重置tips展示状态
    storageOpt.remove(SUPOS_USER_TIPS_ENABLE);
    // 清空
    storageOpt.remove('personInfo');

    if (!cookies?.[SUPOS_COMMUNITY_TOKEN]) {
      if (import.meta.env.MODE === 'development') {
        message.error('开发环境cookie已失效，重新登录生产环境环境，然后复制生产环境的cookie使用');
      } else {
        if (Cookies.get(OMC_MODEL)) {
          console.warn('omc——cookie失效');
          window.location.href = '/403';
        } else {
          console.log('登录cookie不存在，要跳转到登录页');
          window.location.href = systemInfo?.loginPath || LOGIN_URL;
        }
      }
    }
  }, [cookies?.[SUPOS_COMMUNITY_TOKEN]]);

  return null;
};

export default CookieContext;
