import { Hono } from 'hono';
import { openApiRouter } from './open-api';
import { mcpApiRouter } from './copilotkit';

// 路由注册模块
export function registerRoutes(app: Hono) {
  // 注册健康检查路由
  app.route('/', openApiRouter);

  // 注册MCP管理路由
  app.route('/', mcpApiRouter);

  // 根路径
  app.get('/', (c) =>
    c.json(
      {
        message: '🚀 Services Hono API Server',
        version: '1.0.0',
        timestamp: new Date().toISOString(),
        endpoints: {
          health: '/open-api/health/*',
          mcpManagement: '/copilotkit/mcp/*',
        },
      },
      200
    )
  );

  // 404 中间件
  app.notFound((c) => {
    return c.json(
      {
        error: 'Not Found',
        message: 'The requested resource was not found',
      },
      404
    );
  });
  console.log('✅ Hono路由注册完成');
}

export { openApiRouter };
