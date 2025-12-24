import { Hono } from 'hono';
import { config } from '@/config';

// Hono 服务器管理模块
export class ServerManager {
  private server: any;

  constructor(private app: Hono) {}

  // 启动服务器（适配 Bun 运行时）
  start(): void {
    this.server = Bun.serve({
      port: config.port,
      fetch: this.app.fetch,
    });

    console.log(`🚀 Hono Server is running on http://localhost:${config.port}`);
    console.log(`🌍 Environment: ${config.nodeEnv}`);
    console.log('⏹️  Press Ctrl+C to stop the server');
  }

  // 设置服务器信号监听
  setupSignalHandlers(): void {
    process.on('uncaughtException', (err) => {
      console.error('Uncaught Exception:', err);
    });

    process.on('unhandledRejection', (reason, promise) => {
      console.error('Unhandled Rejection at:', promise, 'reason:', reason);
    });

    process.on('SIGINT', () => {
      this.server.stop();
      console.log('✅ Hono Server closed successfully');
      process.exit(0);
    });

    process.on('SIGTERM', () => {
      this.server.stop();
      console.log('✅ Hono Server closed successfully');
      process.exit(0);
    });
  }
}
