import { useEffect } from 'react';
import { useLocation } from 'react-router';
import { childrenRoutes } from '@/routers';
import { useTranslate } from '@/hooks/index.ts';
import { setCurrentMenuInfo, useBaseStore } from '@/stores/base';
import { formatShowName } from '@/utils';

// 改变menu的名称
const useChangeMenuName = () => {
  const { pathname } = useLocation();
  const formatMessage = useTranslate();
  const menuGroup = useBaseStore((state) => state.menuGroup);

  useEffect(() => {
    // 监听location修改名称
    const pathName = pathname?.slice(1);
    const info = menuGroup?.find((f) => {
      if (f.urlType === 1) {
        return f?.url === pathname;
      } else {
        return f.code === pathName;
      }
    });
    if (info) {
      setCurrentMenuInfo(info);
    } else {
      // 内置路由情况
      const interInfo = childrenRoutes?.find((f) => f.path === pathname);
      const parentInfo = menuGroup?.find((f) => {
        return f?.url === interInfo?.handle?.parentPath;
      }) || { id: '-9999', showName: '未配置', type: 2, code: '未配置', sort: 1 };
      if (interInfo && (parentInfo || interInfo?.handle?.parentPath === '/_common')) {
        setCurrentMenuInfo({
          ...parentInfo,
          showName: formatShowName({
            code: interInfo?.handle?.code,
            showName: (interInfo?.handle as any)?.showName,
            formatMessage,
            finallyShowName: parentInfo?.showName,
          }),
        });
      }
    }
  }, [pathname, menuGroup]);
};

export default useChangeMenuName;
