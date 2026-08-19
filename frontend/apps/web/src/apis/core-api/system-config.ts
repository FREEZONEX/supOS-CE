/* eslint-disable @typescript-eslint/no-unused-vars */
import { coreApi } from './core-adapter';

export type RuntimeMetricItem = {
  usagePercent?: number;
  usedBytes?: number;
  totalBytes?: number;
  logicalCores?: number;
  scope?: string;
  source?: string;
  path?: string;
};

export type RuntimeMetrics = {
  cpu?: RuntimeMetricItem;
  memory?: RuntimeMetricItem;
  hostCpu?: RuntimeMetricItem;
  hostMemory?: RuntimeMetricItem;
  disk?: RuntimeMetricItem;
  collectedAt?: number;
};

// 系统配置
export const getSystemConfig = async () => coreApi.get('/system/config');

export const getRuntimeMetrics = async (): Promise<RuntimeMetrics> => coreApi.get('/system/runtime-metrics');

//获取主题配置

export const getAllThemeConfig = async (_params?: any) => ({});
