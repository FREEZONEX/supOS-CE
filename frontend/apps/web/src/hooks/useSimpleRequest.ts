import { useCallback, useEffect, useRef, useState } from 'react';
interface UsePaginationParams {
  fetchApi: any; // API 函数，返回数据及总数
  // 默认参数
  defaultParams?: any;
  // 首次是否请求
  firstNotGetData?: boolean;
  onSuccessCallback?: (data: any) => void;
  autoRefresh?: boolean; // 是否自动刷新
  refreshInterval?: number; // 自动刷新间隔，单位毫秒
  onOnceSuccessCallback?: (data: any) => void;
}

const useSimpleRequest = <T>({
  fetchApi,
  defaultParams = {},
  firstNotGetData,
  onSuccessCallback,
  autoRefresh = false,
  refreshInterval = 5000, // 默认5秒
  onOnceSuccessCallback,
}: UsePaginationParams) => {
  const firstUpdate = useRef(firstNotGetData === true);
  const isOnceRef = useRef(true);
  const [data, setData] = useState<T[]>([]);
  const [loading, setLoading] = useState(false);
  const refreshTimerRef = useRef<NodeJS.Timeout | null>(null);
  const abortControllerRef = useRef<AbortController | null>(null);
  const [paramsData, setParamsData] = useState({
    searchFormData: {},
  });

  const cancelRequest = useCallback(() => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
      abortControllerRef.current = null;
    }
  }, []);

  const clearTime = useCallback(() => {
    if (refreshTimerRef.current) {
      clearTimeout(refreshTimerRef.current);
      refreshTimerRef.current = null;
    }
  }, []);

  const getData = (needLoading = true) => {
    clearTime();
    cancelRequest(); // 取消之前的请求
    abortControllerRef.current = new AbortController(); // 创建新的 AbortController
    if (needLoading) {
      setLoading(true);
    }
    fetchApi(
      {
        ...defaultParams,
        ...paramsData.searchFormData,
      },
      {
        signal: abortControllerRef.current.signal, // 传递 signal
      }
    )
      .then((data: any) => {
        if (abortControllerRef.current?.signal.aborted) {
          return; // 如果请求被中止，则不处理响应
        }
        setData(data);
        if (onSuccessCallback) {
          onSuccessCallback?.(data);
        }
        if (isOnceRef.current && onOnceSuccessCallback) {
          isOnceRef.current = false;
          onOnceSuccessCallback?.(data);
        }
      })
      .catch((error: any) => {
        if (error.name === 'AbortError') {
          console.log('Fetch aborted');
        } else {
          // 处理其他错误
          console.error('Fetch error:', error);
        }
      })
      .finally(() => {
        setLoading(false);
        if (autoRefresh) {
          refreshTimerRef.current = setTimeout(() => getData(false), refreshInterval);
        }
      });
  };

  useEffect(() => {
    if (firstUpdate.current) {
      firstUpdate.current = false;
      if (autoRefresh && !firstNotGetData) {
        // 如果首次不请求数据，但开启了自动刷新，则在首次渲染后启动定时器
        refreshTimerRef.current = setTimeout(() => getData(false), refreshInterval);
      }
      return;
    }
    getData();

    return () => {
      clearTime();
      cancelRequest(); // 组件卸载时取消请求
    };
  }, [paramsData, autoRefresh, refreshInterval]);

  const reload = useCallback(() => {
    clearTime();
    cancelRequest();
    setParamsData((o) => ({
      ...o,
    }));
  }, []);

  // 参数请求
  const setSearchParams = useCallback((value: any, reset: boolean = true) => {
    clearTime();
    cancelRequest();
    setParamsData((o) => {
      const searchParams = reset ? value : Object.assign(o.searchFormData || {}, value || {});
      return {
        ...o,
        searchFormData: searchParams || {},
      };
    });
  }, []);

  const clearData = () => {
    clearTime();
    cancelRequest();
    setData([]);
  };

  return {
    setLoading,
    loading,
    data,
    reload,
    refreshRequest: getData,
    setSearchParams,
    clearData,
  };
};

export default useSimpleRequest;
