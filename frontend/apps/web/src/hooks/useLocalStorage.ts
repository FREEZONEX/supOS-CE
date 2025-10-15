import { useEffect, useState } from 'react';

// 监听localstorage变化
function useLocalStorage(key: string) {
  const [value, setValue] = useState(() => localStorage.getItem(key));

  useEffect(() => {
    const handleStorageChange = (event: StorageEvent) => {
      // 确保是监控的 key 且值有变化时才更新
      if (event.key === key && event.oldValue !== event.newValue) {
        setValue(event.newValue);
      }
    };

    window.addEventListener('storage', handleStorageChange);

    return () => {
      window.removeEventListener('storage', handleStorageChange);
    };
  }, [key]);

  return value;
}

export default useLocalStorage;
