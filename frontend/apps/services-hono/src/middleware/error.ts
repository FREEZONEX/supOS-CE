import { Context, Next } from 'hono';
import { HTTPException } from 'hono/http-exception';
/**
 * 全局错误处理中间件
 */
export const errorHandler = async (c: Context, next: Next) => {
  try {
    await next();
  } catch (err) {
    console.error('Error:', err);

    if (err instanceof HTTPException) {
      return err.getResponse();
    }

    return c.json(
      {
        error: 'Internal Server Error',
        message: err instanceof Error ? err.message : 'Unknown error',
      },
      500
    );
  }
};
