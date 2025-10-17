import { config as envConfig } from 'dotenv';

// 加载环境
envConfig();

// 环境配置
export const config = {
  // 服务器配置
  port: parseInt(process.env.PORT || '3001', 10),
  nodeEnv: process.env.NODE_ENV || 'development',
  // Docker配置
  dockerHost: process.env.DOCKER_HOST || 'localhost',
  dockerPort: parseInt(process.env.DOCKER_PORT || '2375', 10),
} as const;

export type Config = typeof config;
