import { useLocation, useNavigate } from 'react-router';
import { useState } from 'react';
import { useTabsContext } from '@/contexts/tabs-context.ts';
import { canModifyParentHref } from '@/utils/common';

const useNavigateForIframe = ({ path, replaceCurrent = false }: { path: string; replaceCurrent?: boolean }) => {
  const navigate = useNavigate();
  const location = useLocation();
  const { onRemoveTab } = useTabsContext();
  const [security] = useState<boolean | -1>(() => canModifyParentHref());

  const onClick = () => {
    if (security === -1) {
      return;
    }
    if (!security) {
      const currentPath = location.pathname;
      if (replaceCurrent) {
        navigate(path, { replace: true });
        if (currentPath !== path) {
          setTimeout(() => {
            onRemoveTab?.(currentPath);
          }, 0);
        }
      } else {
        navigate(path);
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
