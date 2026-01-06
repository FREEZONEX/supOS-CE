import { createContext, type RefObject, useContext, useEffect } from 'react';

export const TabsLifecycleContext = createContext<{
  isShowRef?: RefObject<boolean>;
  activate: (_: () => void) => void;
  unActivate: (_: () => void) => void;
}>({
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  activate: (_) => {},
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  unActivate: (_) => {},
});
export const useTabLifecycle = () => useContext(TabsLifecycleContext);

// 自定义 Hook：激活时触发
export const useActivate = (cb: () => void) => {
  const { activate } = useTabLifecycle();
  useEffect(() => {
    activate(cb);
  }, [activate, cb]);
};

// 自定义 Hook：未激活时触发
export const useUnActivate = (cb: () => void) => {
  const { unActivate } = useTabLifecycle();
  useEffect(() => {
    unActivate(cb);
  }, [unActivate, cb]);
};

// 判断页面是否显示
export const usePageIsShow = () => {
  const { isShowRef } = useTabLifecycle();
  return isShowRef;
};
