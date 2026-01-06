import { createContext, useContext } from 'react';

export interface TabsContextProps {
  onRefreshTab: (routePath?: string) => void;
  onCloseTab: (routePath?: string) => void;
  onCloseOtherTab: (routePath?: string) => void;
}

const defaultValue = {
  onRefreshTab: () => {},
  onCloseTab: () => {},
  onCloseOtherTab: () => {},
};

export const TabsContext = createContext<{ current?: TabsContextProps }>({
  current: defaultValue,
});
export const useTabsContext = () => {
  const a: any = useContext(TabsContext);
  return {
    TabsContext: a,
    onRefreshTab: a?.current?.onRefreshTab,
    onCloseTab: a?.current?.onCloseTab,
    onCloseOtherTab: a?.current?.onCloseOtherTab,
  };
};
