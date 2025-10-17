import { Hono } from 'hono';
import { openApi } from './open-api';
// import { copilotKitRoutes } from './copilotkit';

// 路由注册模块
export function registerRoutes(app: Hono) {
  // 注册健康检查路由
  app.route('/', openApi);
  // app.route('/', copilotKitRoutes);

  // 根路径
  app.get('/', (c) => {
    return c.json({
      message: '🚀 Services API Server',
      version: '1.0.0',
      timestamp: new Date().toISOString(),
      endpoints: {
        health: '/open-api/health',
      },
    });
  });

  console.log('✅ 路由注册完成');
}

export { openApi };
