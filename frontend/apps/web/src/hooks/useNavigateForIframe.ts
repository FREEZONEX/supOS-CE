import { useLocation, useNavigate } from 'react-router';
import { useEffect, useState } from 'react';
import { useTabsContext } from '@/contexts/tabs-context.ts';
import { canModifyParentHref } from '@/utils/common';

const useNavigateForIframe = ({ path, replaceCurrent = false }: { path: string; replaceCurrent?: boolean }) => {
  const navigate = useNavigate();
  const location = useLocation();
  const { onCloseTab } = useTabsContext();
  const [security, setSecurity] = useState<boolean | -1>(true);
  useEffect(() => {
    const result = canModifyParentHref();
    setSecurity(result);
  }, []);

  const onClick = () => {
    if (security === -1) {
      return;
    }
    if (!security) {
      const currentPath = location.pathname;
      if (replaceCurrent) {
        navigate(path, { replace: true });
      } else {
        navigate(path);
      }
      if (replaceCurrent && currentPath !== path) {
        setTimeout(() => {
          onCloseTab?.(currentPath);
        }, 0);
      }
    } else {
      if (replaceCurrent) {
        window.parent.location.replace(path);
      } else {
        window.parent.location.href = path;
      }
    }
  };

  return {
    security: security !== -1,
    onClick,
  };
};

export default useNavigateForIframe;
