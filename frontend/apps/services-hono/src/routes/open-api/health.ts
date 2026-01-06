import { Hono } from 'hono';
import { config } from '@/config';
import Docker from 'dockerode';
import * as fs from 'node:fs';

const healthRouter = new Hono();

// 详细健康检查
healthRouter.get('/health/server/detailed', (c) =>
  c.json(
    {
      status: 'ok',
      environment: config.nodeEnv,
      version: '1.0.1',
      uptime: process.uptime(),
      memory: process.memoryUsage(),
      platform: process.platform,
      nodeVersion: process.version,
    },
    200
  )
);

// 健康检查 - Docker容器监控
healthRouter.get('/health', async (c) => {
  try {
    const docker = new Docker({
      host: config.dockerHost,
      port: config.dockerPort,
      ca: fs.readFileSync('/certs/ca.pem'),
      cert: fs.readFileSync('/certs/cert.pem'),
      key: fs.readFileSync('/certs/key.pem'),
      version: 'v1.47', // 根据 Docker 版本调整
    });

    const filters = {
      network: ['tier0_edge_network'], // 网络名称或 ID
    };

    // 使用Promise方式调用Docker API
    const containers: any[] = await new Promise((resolve, reject) => {
      docker.listContainers({ all: true, filters: filters }, (err, containers) => {
        if (err) reject(err);
        else resolve(containers || []);
      });
    });

    const total: string[] = [];
    const running: string[] = [];
    const stopped: string[] = [];
    let platformStatus = 'running';

    for (const container of containers) {
      const containerName = container.Names[0].substring(1);
      total.push(containerName);
      if ('running' === container.State) {
        running.push(containerName);
      } else {
        stopped.push(containerName);
        if (containerName === 'backend') {
          platformStatus = 'stop';
        }
      }
    }

    const data = {
      status: platformStatus,
      overview: {
        total: total.length,
        running: running.length,
        stop: stopped.length,
      },
      container: {
        running: running,
        stop: stopped,
      },
    };

    return c.json({ data: data }, 200);
  } catch (error) {
    console.error('supos-docker-health-error:', error);
    return c.json(
      {
        error: 'Docker健康检查失败',
        message: error instanceof Error ? error.message : '未知错误',
      },
      500
    );
  }
});

export { healthRouter };
