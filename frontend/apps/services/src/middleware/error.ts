import { config } from '@/config';

/**
 * 全局错误处理中间件
 */
export const errorHandler = (err: Error, c: any) => {
  console.error('Sup-os-edge-frontend: Server Error:', err);

  return c.json(
    {
      error: config.nodeEnv === 'production' ? 'Internal Server Error' : err.message,
      message: 'Something went wrong on our end',
      timestamp: new Date().toISOString(),
    },
    500
  );
};

/**
 * 404处理中间件
 */
export const notFoundHandler = (c: any) => {
  return c.json(
    {
      error: 'Sup-os-edge-frontend: Not Found',
      message: 'The requested resource was not found',
      path: c.req.path,
      timestamp: new Date().toISOString(),
    },
    404
  );
};
