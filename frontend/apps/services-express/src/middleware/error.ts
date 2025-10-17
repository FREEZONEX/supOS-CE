import { Request, Response } from 'express';
import { config } from '@/config';

/**
 * 全局错误处理中间件
 */
export const errorHandler = (req: Request, res: Response) => {
  console.error('Sup-os-edge-frontend: Server Error:', req);

  res.status(500).json({
    error: config.nodeEnv === 'production' ? 'Internal Server Error' : req,
    message: 'Something went wrong on our end',
    timestamp: new Date().toISOString(),
  });
};

/**
 * 404处理中间件
 */
export const notFoundHandler = (req: Request, res: Response) => {
  res.status(404).json({
    error: 'Sup-os-edge-frontend: Not Found',
    message: 'The requested resource was not found',
    path: req.path,
    timestamp: new Date().toISOString(),
  });
};
