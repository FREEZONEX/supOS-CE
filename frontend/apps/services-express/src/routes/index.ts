import { Express } from 'express';
import { openApiRouter } from './open-api';
// import { copilotKitRoutes } from './copilotkit';

// 路由注册模块
export function registerRoutes(app: Express) {
  // 注册健康检查路由
  app.use('/open-api', openApiRouter);

  // 根路径
  app.get('/', (_, res) => {
    res.json({
      message: '🚀 Services Express API Server',
      version: '1.0.0',
      timestamp: new Date().toISOString(),
      endpoints: {
        health: '/open-api/health',
      },
    });
  });

  console.log('✅ Express路由注册完成');
}

export { openApiRouter };
