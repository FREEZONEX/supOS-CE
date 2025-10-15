import { useNavigate } from 'react-router';
import { useEffect, useState } from 'react';
import { canModifyParentHref } from '@/utils/common';

const useNavigateForIframe = ({ path }: { path: string }) => {
  const navigate = useNavigate();
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
      navigate(path);
    } else {
      window.parent.location.href = path;
    }
  };

  return {
    security: security !== -1,
    onClick,
  };
};

export default useNavigateForIframe;
