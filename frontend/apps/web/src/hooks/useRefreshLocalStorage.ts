// 监听localStorage值变化
import { useEffect, useState } from 'react';

const useRefreshLocalStorage = (localStorageKey: string) => {
  const [storageValue, setStorageValue] = useState(localStorage.getItem(localStorageKey!));

  useEffect(() => {
    const originalSetItem = localStorage.setItem;
    localStorage.setItem = function (key, newValue) {
      const setItemEvent = new CustomEvent('setItemEvent', {
        detail: { key, newValue },
      });
      window.dispatchEvent(setItemEvent);
      originalSetItem.apply(this, [key, newValue]);
    };

    const handleSetItemEvent = (event: any) => {
      const customEvent = event;
      if (event.detail.key === localStorageKey) {
        const updatedValue = customEvent.detail.newValue;
        setStorageValue(updatedValue);
      }
    };

    window.addEventListener('setItemEvent', handleSetItemEvent);

    return () => {
      window.removeEventListener('setItemEvent', handleSetItemEvent);
      localStorage.setItem = originalSetItem;
    };
  }, [localStorageKey]);

  // 返回当前的 storageValue
  return [storageValue];
};

export default useRefreshLocalStorage;
