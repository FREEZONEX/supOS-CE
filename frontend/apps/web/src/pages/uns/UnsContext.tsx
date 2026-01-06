// ThemeContext.js
import { type ReactNode, createContext, useContext } from 'react';
// import useUnsGlobalWs from '@/pages/uns/useUnsGlobalWs.ts';

interface UnsContextType {
  topologyData: { [key: string]: any };
  mountStatus: { [key: string]: any };
}

const UnsContext = createContext<UnsContextType | undefined>(undefined);

export const UnsContextProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const { topologyData = {}, mountStatus = {} } = {};
  // const { topologyData = {}, mountStatus = {} } = useUnsGlobalWs;

  return <UnsContext.Provider value={{ topologyData, mountStatus }}>{children}</UnsContext.Provider>;
};

export const useUnsContext = () => {
  const context = useContext(UnsContext);
  if (!context) {
    throw new Error('useUnsContext must be used within a UnsProvider');
  }
  return context;
};
