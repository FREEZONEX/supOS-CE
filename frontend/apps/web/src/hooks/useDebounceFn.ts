import { useEffect, useMemo, useRef } from 'react';
import { debounce } from 'lodash-es';

type AnyFunction = (...args: any[]) => void;

const useDebounceFn = <T extends AnyFunction>(callback: T, delay = 300) => {
  const callbackRef = useRef(callback);
  callbackRef.current = callback;

  const debouncedCallback = useMemo(() => {
    return debounce((...args: Parameters<T>) => {
      callbackRef.current(...args);
    }, delay);
  }, [delay]);

  useEffect(() => {
    return () => {
      debouncedCallback.cancel();
    };
  }, [debouncedCallback]);

  return debouncedCallback;
};

export default useDebounceFn;
