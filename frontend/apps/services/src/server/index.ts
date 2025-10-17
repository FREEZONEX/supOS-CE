import { serve } from '@hono/node-server';
import { config } from '@/config';

// 服务器管理模块
export class ServerManager {
  private server: any;

  constructor(private app: any) {}

  // 启动服务器
  start(): void {
    this.server = serve({
      fetch: this.app.fetch,
      port: config.port,
    });

    console.log(`🚀 Server is running on http://localhost:${config.port}`);
    console.log(`🌍 Environment: ${config.nodeEnv}`);
    console.log('⏹️  Press Ctrl+C to stop the server');
  }

  // 设置服务器信号监听
  setupSignalHandlers(): void {
    process.on('SIGINT', () => {
      this.server.close();
      console.log('✅ Server closed successfully');
      process.exit(0);
    });

    process.on('SIGINT', () => {
      this.server.close((err: any) => {
        if (err) {
          console.error(err);
          process.exit(1);
        }
        process.exit(0);
      });
    });
  }

  // 获取服务器信息
  getServerInfo() {
    return {
      port: config.port,
      environment: config.nodeEnv,
      platform: process.platform,
      nodeVersion: process.version,
      uptime: process.uptime(),
    };
  }
}
