import React, { createContext, useContext, useState } from 'react';

interface UnsTreeMapContextType {
  isTreeMapVisible: boolean;
  setTreeMapVisible: (visible: boolean) => void;
}

const UnsTreeMapContext = createContext<UnsTreeMapContextType | undefined>(undefined);

export const UnsTreeMapProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [isTreeMapVisible, setTreeMapVisible] = useState(false);

  return (
    <UnsTreeMapContext.Provider value={{ isTreeMapVisible, setTreeMapVisible }}>{children}</UnsTreeMapContext.Provider>
  );
};

export const useUnsTreeMapContext = () => {
  const context = useContext(UnsTreeMapContext);
  if (!context) {
    throw new Error('useTreeMapContext must be used within a TreeMapProvider');
  }
  return context;
};
