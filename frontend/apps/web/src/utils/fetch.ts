type RequestMethod = 'GET' | 'POST' | 'PUT' | 'DELETE';

interface RequestConfig extends RequestInit {
  baseURL?: string;
  params?: Record<string, string>;
  timeout?: number;
  retries?: number;
}

interface StreamResponse<T> {
  onData: (callback: (chunk: T) => void) => void;
  abort: () => void;
  complete: Promise<void>;
}

class FetchHttp {
  private defaultConfig: RequestConfig = {
    baseURL: import.meta.env.VITE_API_BASE_URL,
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${import.meta.env.VITE_API_KEY}`,
    },
    timeout: 10000,
    retries: 2,
  };

  // 请求拦截器
  private requestInterceptors: Array<(config: RequestConfig) => RequestConfig> = [];

  // 响应拦截器
  private responseInterceptors: Array<(response: Response) => Response> = [];

  constructor(config?: RequestConfig) {
    this.defaultConfig = { ...this.defaultConfig, ...config };
  }

  // 核心请求方法
  async request<T = any>(url: string, method: RequestMethod = 'GET', data?: any, config?: RequestConfig): Promise<T> {
    const mergedConfig: RequestConfig = {
      ...this.defaultConfig,
      ...config,
      method,
      headers: {
        ...this.defaultConfig.headers,
        ...config?.headers,
      },
    };

    // 执行请求拦截器
    let finalConfig = mergedConfig;
    for (const interceptor of this.requestInterceptors) {
      finalConfig = interceptor(finalConfig);
    }

    // 处理 query 参数
    const fullUrl = this.buildFullUrl(url, finalConfig);
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    const body = this.getRequestBody(data, finalConfig.headers?.['Content-Type']);

    // 支持超时取消
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), finalConfig.timeout);

    try {
      let response = await fetch(fullUrl, {
        ...finalConfig,
        body,
        signal: controller.signal,
      });

      clearTimeout(timeoutId);

      // 执行响应拦截器
      for (const interceptor of this.responseInterceptors) {
        response = interceptor(response);
      }

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      return this.handleResponse<T>(response);
    } catch (error) {
      if (finalConfig.retries && finalConfig.retries > 0) {
        return this.request<T>(url, method, data, {
          ...finalConfig,
          retries: finalConfig.retries - 1,
        });
      }
      throw error;
    }
  }

  // 流式请求专用方法
  stream<T = any>(url: string, method: RequestMethod = 'GET', data?: any, config?: RequestConfig): StreamResponse<T> {
    const controller = new AbortController();
    let reader: ReadableStreamDefaultReader<Uint8Array> | null = null;
    let completed = false;

    const process = async () => {
      try {
        const response = await this.request(url, method, data, {
          ...config,
          signal: controller.signal,
        });

        if (!response.body) {
          throw new Error('No response body');
        }

        reader = response.body.getReader();
        const decoder = new TextDecoder();

        while (true) {
          // eslint-disable-next-line @typescript-eslint/ban-ts-comment
          // @ts-expect-error
          const { done, value } = await reader.read();
          if (done) {
            completed = true;
            return;
          }

          const chunk = decoder.decode(value);
          try {
            const data = JSON.parse(chunk) as T;
            dataCallbacks.forEach((cb) => cb(data));
          } catch (e) {
            console.error('Stream parse error:', e);
          }
        }
      } catch (error) {
        if (!completed) {
          errorCallbacks.forEach((cb) => cb(error as Error));
        }
      }
    };

    const dataCallbacks: Array<(chunk: T) => void> = [];
    const errorCallbacks: Array<(error: Error) => void> = [];

    process();

    return {
      onData: (callback) => {
        dataCallbacks.push(callback);
      },
      abort: () => {
        reader?.cancel();
        controller.abort();
        completed = true;
      },
      complete: new Promise((resolve) => {
        const checkComplete = () => {
          if (completed) resolve();
          else setTimeout(checkComplete, 100);
        };
        checkComplete();
      }),
    };
  }

  // 快捷方法
  get<T = any>(url: string, config?: RequestConfig) {
    return this.request<T>(url, 'GET', undefined, config);
  }

  post<T = any>(url: string, data?: any, config?: RequestConfig) {
    return this.request<T>(url, 'POST', data, config);
  }

  // 添加拦截器
  useRequestInterceptor(interceptor: (config: RequestConfig) => RequestConfig) {
    this.requestInterceptors.push(interceptor);
  }

  useResponseInterceptor(interceptor: (response: Response) => Response) {
    this.responseInterceptors.push(interceptor);
  }

  private buildFullUrl(url: string, config: RequestConfig): string {
    const baseURL = config.baseURL || '';
    const query = new URLSearchParams(config.params).toString();
    return `${baseURL}${url}${query ? `?${query}` : ''}`;
  }

  private getRequestBody(data: any, contentType?: string): BodyInit | null {
    if (!data) return null;

    if (contentType?.includes('application/json')) {
      return JSON.stringify(data);
    }

    if (data instanceof FormData) {
      return data;
    }

    return String(data);
  }

  private async handleResponse<T>(response: Response): Promise<T> {
    const contentType = response.headers.get('Content-Type');

    if (contentType?.includes('application/json')) {
      return response.json();
    }

    if (contentType?.includes('text/')) {
      return response.text() as Promise<T>;
    }

    return response.blob() as Promise<T>;
  }
}

// 创建实例
export const fetchHttp = new FetchHttp();
