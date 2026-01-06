import { useLocation, useNavigate } from 'react-router';
import { useEffect } from 'react';
import { useBaseStore } from '@/stores/base';

const NotPage = () => {
  const { pathname } = useLocation();
  const navigate = useNavigate();
  useEffect(() => {
    const originMenu = useBaseStore.getState().originMenu;
    const isAuthRoute = originMenu?.find((f) => '/' + f.code === pathname || f?.url === pathname);
    if (isAuthRoute?.code || pathname === '/403') {
      // 如果没权限
      navigate('/403');
    } else {
      // 如果没菜单
      navigate('/404');
    }
  }, [pathname]);

  return <div></div>;
};

export default NotPage;
