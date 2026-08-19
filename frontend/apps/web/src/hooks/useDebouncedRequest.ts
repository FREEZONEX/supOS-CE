import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { debounce } from 'lodash-es';

interface UseDebouncedRequestOptions<TData> {
  debounceMs?: number;
  initialData: TData;
}

const useDebouncedRequest = <TArgs extends unknown[], TData>(
  requestFn: (...args: TArgs) => Promise<TData>,
  options: UseDebouncedRequestOptions<TData>
) => {
  const { debounceMs = 300, initialData } = options;
  const requestIdRef = useRef(0);
  const requestFnRef = useRef(requestFn);
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<TData>(initialData);

  useEffect(() => {
    requestFnRef.current = requestFn;
  }, [requestFn]);

  const runImmediate = useCallback(async (...args: TArgs): Promise<TData | undefined> => {
    const requestId = ++requestIdRef.current;
    setLoading(true);

    try {
      const response = await requestFnRef.current(...args);
      if (requestId === requestIdRef.current) {
        setData(response);
        return response;
      }
      return undefined;
    } finally {
      if (requestId === requestIdRef.current) {
        setLoading(false);
      }
    }
  }, []);

  const runDebounced = useMemo(() => {
    return debounce((...args: TArgs) => {
      void runImmediate(...args);
    }, debounceMs);
  }, [debounceMs, runImmediate]);

  useEffect(() => {
    return () => {
      runDebounced.cancel();
    };
  }, [runDebounced]);

  return {
    loading,
    data,
    runImmediate,
    runDebounced,
  };
};

export default useDebouncedRequest;
