import { Express } from 'express';
import { config } from '@/config';

// 服务器管理模块
export class ServerManager {
  private server: any;

  constructor(private app: Express) {}

  // 启动服务器
  start(): void {
    this.server = this.app.listen(config.port, () => {
      console.log(`🚀 Express Server is running on http://localhost:${config.port}`);
      console.log(`🌍 Environment: ${config.nodeEnv}`);
      console.log('⏹️  Press Ctrl+C to stop the server');
    });
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
      this.server.close(() => {
        console.log('✅ Express Server closed successfully');
        process.exit(0);
      });
    });

    process.on('SIGTERM', () => {
      this.server.close((err: any) => {
        if (err) {
          console.error('Error closing server:', err);
          process.exit(1);
        }
        console.log('✅ Express Server closed successfully');
        process.exit(0);
      });
    });
  }
}
