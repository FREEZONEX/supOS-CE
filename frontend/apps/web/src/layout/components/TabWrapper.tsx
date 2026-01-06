import { type ReactNode, useRef, useState } from 'react';
import { TabsLifecycleContext } from '@/contexts/tabs-lifecycle-context.ts';
import { useMount, useUpdateEffect } from 'ahooks';

const TabWrapper = ({ children, isActive }: { children: ReactNode; isActive: boolean }) => {
  const [activated, setActivated] = useState(isActive);
  const isFirstRender = useRef(true);
  const isShowRef = useRef(true);

  useMount(() => {
    isFirstRender.current = false;
  });

  useUpdateEffect(() => {
    if (isActive && !activated) {
      setActivated(true);
      isShowRef.current = true;
    } else if (!isActive && activated) {
      setActivated(false);
      isShowRef.current = false;
    }
  }, [isActive, activated]); // 监听 isActive 和 activated 的变化

  const activate = (callback: () => void) => {
    if (!activated && isActive && !isFirstRender.current) {
      callback();
    }
  };

  const unActivate = (callback: () => void) => {
    if (activated && !isActive) {
      callback();
    }
  };

  return (
    <TabsLifecycleContext.Provider value={{ activate, unActivate, isShowRef }}>
      {children}
    </TabsLifecycleContext.Provider>
  );
};

export { TabWrapper };
