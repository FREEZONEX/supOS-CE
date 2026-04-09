import { useCallback, useEffect, useRef, useState } from 'react';

interface UseKongTableParams {
  fetchApi: (params?: Record<string, unknown>) => Promise<any>;
  pageSize?: number;
}

/**
 * Kong Admin API uses cursor-based pagination (offset, not page numbers).
 * This hook wraps that into a simple load/refresh pattern.
 */
const useKongTable = <T = any>({ fetchApi, pageSize = 100 }: UseKongTableParams) => {
  const [data, setData] = useState<T[]>([]);
  const [loading, setLoading] = useState(false);
  const [offset, setOffset] = useState<string | undefined>(undefined);
  const [hasNext, setHasNext] = useState(false);
  const mountedRef = useRef(true);

  const fetchData = useCallback(
    async (nextOffset?: string) => {
      setLoading(true);
      try {
        const params: Record<string, unknown> = { size: pageSize };
        if (nextOffset) params.offset = nextOffset;
        const res: any = await fetchApi(params);
        if (!mountedRef.current) return;
        setData(res?.data ?? []);
        setOffset(res?.offset);
        setHasNext(!!res?.offset);
      } catch {
        if (!mountedRef.current) return;
        setData([]);
      } finally {
        if (mountedRef.current) setLoading(false);
      }
    },
    [fetchApi, pageSize]
  );

  const refresh = useCallback(() => fetchData(), [fetchData]);

  const loadNext = useCallback(() => {
    if (hasNext && offset) fetchData(offset);
  }, [fetchData, hasNext, offset]);

  useEffect(() => {
    mountedRef.current = true;
    fetchData();
    return () => {
      mountedRef.current = false;
    };
  }, [fetchData]);

  return { data, loading, setLoading, refresh, hasNext, loadNext };
};

export default useKongTable;
