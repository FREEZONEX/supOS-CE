import { MCPClient } from './mcp-client';
import { parseTransportUrl } from './path';
import { McpClientOptions, TransportConfig } from '@/types';

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
          // 断开连接并清理MCP内部服务
          if (entry.isConnected) {
            await entry.client.close();
          }
        } catch (error) {
          console.error(`清理MCP客户端连接时出错 (端点: ${entry.endpoint}):`, error);
        } finally {
          // 从缓存中删除
          this.mcpClientCache.delete(key);
          console.log(`已清理过期MCP客户端 (端点: ${entry.endpoint})`);
        }
      }
    }

    if (expiredKeys.length > 0) {
      console.log(`已清理 ${expiredKeys.length} 个过期MCP客户端`);
    }
  }

  /**
   * 健康检查和连接状态管理
   */
  private async checkAndRepairClient(entry: MCPClientEntry): Promise<boolean> {
    try {
      // 如果客户端未连接，尝试重新连接
      if (!entry.isConnected) {
        console.log(`正在重新连接MCP客户端 (端点: ${entry.endpoint})`);
        await entry.client.connect();
        entry.isConnected = true;
        entry.lastUsed = Date.now();
        console.log(`成功重新连接MCP客户端 (端点: ${entry.endpoint})`);
      }

      // 更新最后使用时间
      entry.lastUsed = Date.now();
      return true;
    } catch (error) {
      console.error(`检查/修复MCP客户端时出错 (端点: ${entry.endpoint}):`, error);
      // 标记为未连接，将在下次调用时尝试重新连接
      entry.isConnected = false;
      return false;
    }
  }

  /**
   * 创建新的MCP客户端
   */
  private async createNewMCPClient(config: any, props: any): Promise<MCPClient> {
    console.log(`正在创建新的MCP客户端 (端点: ${config.endpoint})`);

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

    console.log(`成功创建并连接MCP客户端 (端点: ${config.endpoint})`);
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
      console.log(`找到缓存的MCP客户端 (端点: ${endpoint})`);

      // 检查和修复客户端连接
      const isHealthy = await this.checkAndRepairClient(cachedEntry);

      if (isHealthy) {
        // 返回缓存的客户端
        return cachedEntry.client;
      } else {
        // 健康检查失败，从缓存中移除
        this.mcpClientCache.delete(cacheKey);
        console.log(`从缓存中移除不健康的MCP客户端 (端点: ${endpoint})`);
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
    console.log(`已缓存新的MCP客户端 (端点: ${endpoint})`);

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
        // 断开连接并清理MCP内部服务
        if (entry.isConnected) {
          await entry.client.close();
        }
      } catch (error) {
        console.error(`断开MCP客户端连接时出错 (端点: ${endpoint}):`, error);
      } finally {
        // 从缓存中删除
        this.mcpClientCache.delete(cacheKey);
        console.log(`已从缓存中移除MCP客户端 (端点: ${endpoint})`);
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

    console.log('已清理所有MCP客户端缓存');
  }

  /**
   * 获取只读的MCP客户端缓存视图
   */
  getMCPClientCache(): ReadonlyArray<Readonly<Omit<MCPClientEntry & TransportConfig, 'client'>>> {
    return Array.from(this.mcpClientCache.entries()).map(([key, entry]) => {
      const config = parseTransportUrl(entry.endpoint);
      return {
        endpoint: entry.endpoint || key,
        lastUsed: entry.lastUsed,
        isConnected: entry.isConnected,
        ...config,
      };
    });
  }
  /**
   * 根据endpoint获取特定MCP客户端的工具列表
   */
  async getToolsListByEndpoint(endpoint?: string): Promise<{
    success: boolean;
    data: Array<{
      endpoint: string;
      tools: Array<{
        name: string;
        description: string;
      }>;
    }>;
    error?: string;
  }> {
    try {
      const result: Array<{
        endpoint: string;
        tools: Array<{
          name: string;
          description: string;
        }>;
      }> = [];

      // 如果指定了endpoint，只获取该endpoint的工具列表
      if (endpoint) {
        const entry = this.mcpClientCache.get(endpoint);
        if (entry && entry.isConnected) {
          try {
            const toolsMap = await entry.client.tools();
            const tools = Object.entries(toolsMap).map(([name, tool]) => ({
              name,
              description: tool.description || '',
            }));
            result.push({ endpoint, tools });
          } catch (error) {
            console.error(`获取工具列表失败 (端点: ${endpoint}):`, error);
          }
        }
      } else {
        // 获取所有已连接客户端的工具列表
        for (const [endpointKey, entry] of this.mcpClientCache.entries()) {
          if (entry.isConnected) {
            try {
              const toolsMap = await entry.client.tools();
              const tools = Object.entries(toolsMap).map(([name, tool]) => ({
                name,
                description: tool.description || '',
              }));
              result.push({ endpoint: endpointKey, tools });
            } catch (error) {
              console.error(`获取工具列表失败 (端点: ${endpointKey}):`, error);
            }
          }
        }
      }

      return { success: true, data: result };
    } catch (error) {
      console.error('获取工具列表失败:', error);
      return {
        success: false,
        data: [],
        error: `获取工具列表失败: ${error instanceof Error ? error.message : String(error)}`,
      };
    }
  }

  /**
   * 刷新指定的MCP客户端 - 通过重新连接来刷新状态
   * @param endpoint MCP客户端端点
   * @returns 刷新结果
   */
  async refreshMCPClient(endpoint: string): Promise<{ success: boolean; message: string }> {
    const entry = this.mcpClientCache.get(endpoint);

    if (!entry) {
      return { success: false, message: `MCP客户端 ${endpoint} 不存在` };
    }

    try {
      console.log(`开始刷新MCP客户端: ${endpoint}`);

      // 如果客户端已连接，先断开连接
      if (entry.isConnected) {
        await entry.client.close();
        entry.isConnected = false;
        console.log(`已断开MCP客户端连接: ${endpoint}`);
      }

      // 重新连接客户端
      await entry.client.connect();
      entry.isConnected = true;
      entry.lastUsed = Date.now();

      // 清除工具缓存，强制重新获取工具列表
      entry.client.clearToolsCache();

      console.log(`成功刷新MCP客户端: ${endpoint}`);
      return { success: true, message: `MCP客户端 ${endpoint} 刷新成功` };
    } catch (error) {
      console.error(`刷新MCP客户端 ${endpoint} 失败:`, error);
      entry.isConnected = false;
      return {
        success: false,
        message: `MCP客户端 ${endpoint} 刷新失败: ${error instanceof Error ? error.message : String(error)}`,
      };
    }
  }

  /**
   * 重启指定的MCP客户端 - 先停止再重新创建
   * @param endpoint MCP客户端端点
   * @returns 重启结果
   */
  async restartMCPClient(endpoint: string): Promise<{ success: boolean; message: string }> {
    const cacheKey = endpoint;
    const entry = this.mcpClientCache.get(cacheKey);

    if (!entry) {
      return { success: false, message: `MCP客户端 ${endpoint} 不存在` };
    }

    try {
      console.log(`开始重启MCP客户端: ${endpoint}`);

      // 保存客户端配置信息
      const config = { endpoint };
      const props = parseTransportUrl(endpoint);

      // 先停止客户端
      await this.removeMCPClient(endpoint);

      // 重新创建客户端
      const newClient = await this.createNewMCPClient(config, props);

      // 更新缓存
      const newEntry: MCPClientEntry = {
        client: newClient,
        endpoint: endpoint,
        lastUsed: Date.now(),
        isConnected: true,
      };

      this.mcpClientCache.set(cacheKey, newEntry);

      console.log(`成功重启MCP客户端: ${endpoint}`);
      return { success: true, message: `MCP客户端 ${endpoint} 重启成功` };
    } catch (error) {
      console.error(`重启MCP客户端 ${endpoint} 失败:`, error);
      return {
        success: false,
        message: `MCP客户端 ${endpoint} 重启失败: ${error instanceof Error ? error.message : String(error)}`,
      };
    }
  }

  /**
   * 停止指定的MCP客户端 - 断开连接但不从缓存中移除
   * @param endpoint MCP客户端端点
   * @returns 停止结果
   */
  async stopMCPClient(endpoint: string): Promise<{ success: boolean; message: string }> {
    const entry = this.mcpClientCache.get(endpoint);

    if (!entry) {
      return { success: false, message: `MCP客户端 ${endpoint} 不存在` };
    }

    try {
      console.log(`开始停止MCP客户端: ${endpoint}`);

      // 断开连接
      if (entry.isConnected) {
        await entry.client.close();
        entry.isConnected = false;
        console.log(`已停止MCP客户端: ${endpoint}`);
      } else {
        console.log(`MCP客户端 ${endpoint} 已经处于停止状态`);
      }

      // 更新最后使用时间
      entry.lastUsed = Date.now();

      console.log(`成功停止MCP客户端: ${endpoint}`);
      return { success: true, message: `MCP客户端 ${endpoint} 停止成功` };
    } catch (error) {
      console.error(`停止MCP客户端 ${endpoint} 失败:`, error);
      return {
        success: false,
        message: `MCP客户端 ${endpoint} 停止失败: ${error instanceof Error ? error.message : String(error)}`,
      };
    }
  }
}

const mcpManager = new MCPClientManager();

export { mcpManager };
