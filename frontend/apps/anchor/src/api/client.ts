// 同源调用 backend 内部 API（iframe 内共享主应用 session cookie）
export class ApiError extends Error {
  status: number;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

type Envelope<T> = { code: number; msg: string; data: T };

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, { credentials: 'include', ...init });
  const contentType = res.headers.get('content-type') || '';
  if (!contentType.includes('application/json')) {
    if (!res.ok) throw new ApiError(res.status, await res.text());
    return undefined as T;
  }
  const body = (await res.json()) as Envelope<T>;
  if (!res.ok || (body.code && body.code !== 200)) {
    throw new ApiError(body.code || res.status, body.msg || res.statusText);
  }
  return body.data;
}

export const api = {
  get: <T>(path: string, params?: Record<string, string | number | undefined>) => {
    const search = new URLSearchParams();
    Object.entries(params || {}).forEach(([k, v]) => {
      if (v !== undefined && v !== '') search.set(k, String(v));
    });
    const query = search.toString();
    return request<T>(query ? `${path}?${query}` : path);
  },
  post: <T>(path: string, body?: unknown) =>
    request<T>(path, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body ?? {}),
    }),
  put: <T>(path: string, body?: unknown) =>
    request<T>(path, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body ?? {}),
    }),
  delete: <T>(path: string) => request<T>(path, { method: 'DELETE' }),
  upload: <T>(path: string, form: FormData) => request<T>(path, { method: 'POST', body: form }),
};
