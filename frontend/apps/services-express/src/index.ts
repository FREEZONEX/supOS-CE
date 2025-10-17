import express from 'express';
import { copilotkitHandler, errorHandler, notFoundHandler } from '@/middleware';
import { registerRoutes } from '@/routes';
import { ServerManager } from '@/server';

// 创建Express应用
const app = express();

// 应用基础中间件

// 应用自定义中间件
// copilotkit => ai
app.use('/copilotkit', copilotkitHandler);

// 注册所有路由
registerRoutes(app);

// 应用错误处理中间件
app.use(notFoundHandler);
app.use(errorHandler);

// 创建服务器管理器并启动服务器
const serverManager = new ServerManager(app);
serverManager.setupSignalHandlers();
serverManager.start();

export default app;
