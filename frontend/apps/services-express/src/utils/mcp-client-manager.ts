import { MCPClient } from './mcp-client';
import { parseTransportUrl } from './path';
import { McpClientOptions } from '@/types';

// MCP客户端缓存条目接口
interface MCPClientEntry {
  client: MCPClient;
  endpoint: string;
  lastUsed: number;
  isConnected: boolean;
}

// 缓存TTL (30分钟)
const CLIENT_CACHE_TTL = 30 * 60 * 1000;

/**
 * MCP客户端管理器类
 */
export class MCPClientManager {
  // MCP客户端缓存映射
  private mcpClientCache: Map<string, MCPClientEntry>;

  constructor() {
    this.mcpClientCache = new Map<string, MCPClientEntry>();
  }

  /**
   * 清理过期的MCP客户端缓存 - 主动触发
   */
  async cleanupExpiredClients(): Promise<void> {
    const now = Date.now();
    const expiredKeys: string[] = [];

    // 查找过期的客户端
    for (const [key, entry] of this.mcpClientCache.entries()) {
      if (now - entry.lastUsed > CLIENT_CACHE_TTL) {
        expiredKeys.push(key);
      }
    }

    // 清理过期的客户端
    for (const key of expiredKeys) {
      const entry = this.mcpClientCache.get(key);
      if (entry) {
        try {
          // 断开连接
          if (entry.isConnected) {
            await entry.client.close();
          }
        } catch (error) {
          console.error(`Error disconnecting MCP client for key ${key}:`, error);
        } finally {
          // 从缓存中删除
          this.mcpClientCache.delete(key);
          console.log(`Cleaned up expired MCP client for endpoint: ${entry.endpoint}`);
        }
      }
    }

    if (expiredKeys.length > 0) {
      console.log(`Cleaned up ${expiredKeys.length} expired MCP clients`);
    }
  }

  /**
   * 健康检查和连接状态管理
   */
  private async checkAndRepairClient(entry: MCPClientEntry): Promise<boolean> {
    try {
      // 如果客户端未连接，尝试重新连接
      if (!entry.isConnected) {
        console.log(`Reconnecting MCP client for endpoint: ${entry.endpoint}`);
        await entry.client.connect();
        entry.isConnected = true;
        entry.lastUsed = Date.now();
        console.log(`Successfully reconnected MCP client for endpoint: ${entry.endpoint}`);
      }

      // 更新最后使用时间
      entry.lastUsed = Date.now();
      return true;
    } catch (error) {
      console.error(`Error checking/reparing MCP client for endpoint ${entry.endpoint}:`, error);
      // 标记为未连接，将在下次调用时尝试重新连接
      entry.isConnected = false;
      return false;
    }
  }

  /**
   * 创建新的MCP客户端
   */
  private async createNewMCPClient(config: any, props: any): Promise<MCPClient> {
    console.log(`Creating new MCP client for endpoint: ${config.endpoint}`);

    // 创建客户端配置
    const clientOptions: McpClientOptions = {
      serverUrl: props.serverUrl,
      transportType: props.transportType,
      clientName: props.clientName || 'copilotkit-mcp-client',
      headers: props.headers,
      stdioConfig: props.stdioConfig,
    };

    // 创建新的MCP客户端
    const mcpClient = new MCPClient(clientOptions);

    // 连接到服务器
    await mcpClient.connect();

    console.log(`Successfully created and connected MCP client for endpoint: ${config.endpoint}`);
    return mcpClient;
  }

  /**
   * 获取或创建MCP客户端的主要方法
   */
  async getOrCreateMCPClient(config: any): Promise<MCPClient> {
    const endpoint = config.endpoint;
    const cacheKey = endpoint;

    // 检查缓存中是否存在客户端
    const cachedEntry = this.mcpClientCache.get(cacheKey);

    if (cachedEntry) {
      console.log(`Found cached MCP client for endpoint: ${endpoint}`);

      // 检查和修复客户端连接
      const isHealthy = await this.checkAndRepairClient(cachedEntry);

      if (isHealthy) {
        // 返回缓存的客户端
        return cachedEntry.client;
      } else {
        // 健康检查失败，从缓存中移除
        this.mcpClientCache.delete(cacheKey);
        console.log(`Removed unhealthy MCP client from cache for endpoint: ${endpoint}`);
      }
    }

    // 没有缓存或缓存无效，创建新的客户端
    const props = parseTransportUrl(endpoint);
    const mcpClient = await this.createNewMCPClient(config, props);

    // 缓存新的客户端
    const newEntry: MCPClientEntry = {
      client: mcpClient,
      endpoint: endpoint,
      lastUsed: Date.now(),
      isConnected: true,
    };

    this.mcpClientCache.set(cacheKey, newEntry);
    console.log(`Cached new MCP client for endpoint: ${endpoint}`);

    return mcpClient;
  }

  /**
   * 手动清理特定的MCP客户端缓存
   */
  async removeMCPClient(endpoint: string): Promise<void> {
    const cacheKey = endpoint;
    const entry = this.mcpClientCache.get(cacheKey);

    if (entry) {
      try {
        // 断开连接
        if (entry.isConnected) {
          await entry.client.close();
        }
      } catch (error) {
        console.error(`Error disconnecting MCP client for endpoint ${endpoint}:`, error);
      } finally {
        // 从缓存中删除
        this.mcpClientCache.delete(cacheKey);
        console.log(`Removed MCP client from cache for endpoint: ${endpoint}`);
      }
    }
  }

  /**
   * 获取当前缓存的客户端数量
   */
  getMCPClientCount(): number {
    return this.mcpClientCache.size;
  }

  /**
   * 清理所有MCP客户端缓存
   */
  async removeAllMCPClient(): Promise<void> {
    const endpoints = Array.from(this.mcpClientCache.keys());

    for (const endpoint of endpoints) {
      await this.removeMCPClient(endpoint);
    }

    console.log('Removed all MCP clients from cache');
  }
}
