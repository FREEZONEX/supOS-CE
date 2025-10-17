import { Hono } from 'hono';
import { config } from '@/config';
import fs from 'fs';
import Docker from 'dockerode';

const healthRoutes = new Hono();

// 详细健康检查
healthRoutes.get('/health/server/detailed', (c) => {
  return c.json({
    status: 'ok',
    environment: config.nodeEnv,
    version: '1.0.1',
    uptime: process.uptime(),
    memory: process.memoryUsage(),
    platform: process.platform,
    nodeVersion: process.version,
    fdasfasdfasd: 1,
  });
});

// 健康检查 - Docker容器监控
healthRoutes.get('/health', async (c) => {
  const docker = new Docker({
    host: config.dockerHost,
    port: config.dockerPort,
    ca: fs.readFileSync('/certs/ca.pem'),
    cert: fs.readFileSync('/certs/cert.pem'),
    key: fs.readFileSync('/certs/key.pem'),
    version: 'v1.47', // 根据 Docker 版本调整
  });

  const filters = {
    network: ['supos_default_network'], // 网络名称或 ID
  };

  try {
    const containers: any = await new Promise((resolve, reject) => {
      docker.listContainers({ all: true, filters: filters }, (err, containers) => {
        if (err) reject(err);
        else resolve(containers);
      });
    });

    const total = [];
    const running = [];
    const stopped = [];
    let platformStatus = 'running';

    for (const k in containers) {
      const containerName = containers[k].Names[0].substring(1);
      total.push(containerName);
      if ('running' === containers[k].State) {
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

    return c.json({ data: data });
  } catch (error) {
    return c.json(
      {
        error: 'Failed to get docker information',
        message: error,
      },
      500
    );
  }
});

export { healthRoutes };
