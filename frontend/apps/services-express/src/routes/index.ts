import { Express } from 'express';
import { openApiRouter } from './open-api';
import { mcpApiRouter } from './copilotkit';
// import { copilotKitRoutes } from './copilotkit';

// 路由注册模块
export function registerRoutes(app: Express) {
  // 注册健康检查路由
  app.use('/open-api', openApiRouter);

  // 注册MCP管理路由
  app.use('/copilotkit/mcp', mcpApiRouter);

  // 根路径
  app.get('/', (_, res) => {
    res.json({
      message: '🚀 Services Express API Server',
      version: '1.0.0',
      timestamp: new Date().toISOString(),
      endpoints: {
        health: '/open-api/health',
        mcpManagement: '/mcp-management',
      },
    });
  });

  console.log('✅ Express路由注册完成');
}

export { openApiRouter, mcpApiRouter };
