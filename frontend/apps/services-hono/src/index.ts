import { Hono } from 'hono';
import { logger } from 'hono/logger';
import { ServerManager } from './server';
import { registerRoutes } from '@/routes';
import { errorHandler, copilotkitHandler } from '@/middleware';

const app = new Hono();

app.use(logger());

// 创建并启动服务器
const serverManager = new ServerManager(app);
// 注册路由
registerRoutes(app);

// 注册 copilotkit 路由
app.use('/copilotkit', copilotkitHandler);
// 中间件: 自定义错误处理
app.use('*', errorHandler);

serverManager.setupSignalHandlers();
serverManager.start();

export default app;
