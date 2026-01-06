import express, { Request, Response } from 'express';
import { config } from '@/config';
import fs from 'fs';
import Docker from 'dockerode';

const healthRouter = express.Router();

// 详细健康检查
healthRouter.get('/health/server/detailed', (_: Request, res: Response) => {
  res.status(200).json({
    status: 'ok',
    environment: config.nodeEnv,
    version: '1.0.1',
    uptime: process.uptime(),
    memory: process.memoryUsage(),
    platform: process.platform,
    nodeVersion: process.version,
  });
});

// 健康检查 - Docker容器监控
healthRouter.get('/health', (_: Request, res: Response) => {
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
  docker.listContainers({ all: true, filters: filters }, (err, containers: any) => {
    if (err) throw err;
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
    res.status(200).json({ data: data });
  });
});

export { healthRouter };
