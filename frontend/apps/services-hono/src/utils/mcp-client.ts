import { MCPTool, MCPClient as MCPClientInterface } from '@copilotkit/runtime';
import { Client } from '@modelcontextprotocol/sdk/client/index.js';
import { SSEClientTransport } from '@modelcontextprotocol/sdk/client/sse.js';
import { StdioClientTransport } from '@modelcontextprotocol/sdk/client/stdio.js';
import { StreamableHTTPClientTransport } from '@modelcontextprotocol/sdk/client/streamableHttp.js';
import { JSONRPCMessage } from '@modelcontextprotocol/sdk/types';
import { McpClientOptions, TransportType } from '@/types';

export class MCPClient implements MCPClientInterface {
  // 客户端
  private client: Client;
  // 传输
  private transport: SSEClientTransport | StdioClientTransport | StreamableHTTPClientTransport;
  // 参数类型
  private transportType: TransportType;
  private serverUrl: URL;
  private clientName: string;
  private onMessage: (message: Record<string, unknown>) => void;
  private onError: (error: Error) => void;
  private onOpen: () => void;
  private onClose: () => void;
  private isConnected = false;
  private headers?: Record<string, string>;
  private stdioConfig?: any;
  // 工具缓存，避免重复获取
  private toolsCache: Record<string, MCPTool> | null = null;

  constructor(options: McpClientOptions) {
    this.serverUrl = options.transportType === 'stdio' ? new URL('http://localhost:3000') : new URL(options.serverUrl);
    this.transportType = options.transportType || 'stdio';
    this.headers = options.headers || {};
    this.stdioConfig = options.stdioConfig || {};
    this.clientName = options.clientName || 'cpk-mcp-client';
    this.onMessage = options.onMessage || ((message) => console.log('收到消息:', message));
    this.onError = options.onError || ((error) => console.error('错误:', error));
    this.onOpen = options.onOpen || (() => console.log('连接已打开'));
    this.onClose = options.onClose || (() => console.log('连接已关闭'));
    this.transport = this.createTransport(this.transportType);

    // 初始化客户端
    this.client = new Client({
      name: this.clientName,
      version: '0.0.1',
    });

    // 设置事件处理器
    this.transport.onmessage = this.handleMessage.bind(this);
    this.transport.onerror = this.handleError.bind(this);
    this.transport.onclose = this.handleClose.bind(this);
  }

  private handleMessage(message: JSONRPCMessage): void {
    try {
      this.onMessage(message as Record<string, unknown>);
    } catch (error) {
      this.onError(error instanceof Error ? error : new Error(`处理消息失败: ${error}`));
    }
  }

  private handleError(error: Error): void {
    this.onError(error);
    if (this.isConnected) {
      this.isConnected = false;
    }
  }

  private handleClose(): void {
    this.isConnected = false;
    this.onClose();
  }

  /**
   * 连接到MCP服务器
   */
  public async connect(): Promise<void> {
    try {
      console.log(`正在连接到MCP服务器 (客户端: ${this.clientName})`);

      // 连接客户端（这会连接传输层）
      await this.client.connect(this.transport as any);

      this.isConnected = true;
      console.log(`已连接到MCP服务器 (客户端: ${this.clientName})`);
      this.onOpen();
    } catch (error) {
      console.error(`连接到MCP服务器失败 (客户端: ${this.clientName}):`, error);
      this.onError(error instanceof Error ? error : new Error(String(error)));
      throw error;
    }
  }

  /**
   * 返回工具名称到MCPTool对象的映射
   * 此方法符合CopilotKit接口的预期
   */
  public async tools(): Promise<Record<string, MCPTool>> {
    try {
      // 如果缓存可用，从缓存返回
      if (this.toolsCache) {
        return this.toolsCache;
      }

      // 获取原始工具数据
      const rawToolsResult = await this.client.listTools();

      // 转换为预期格式
      const toolsMap: Record<string, MCPTool> = {};

      if (rawToolsResult) {
        const toolsArray = this.extractToolsArray(rawToolsResult);

        if (toolsArray) {
          toolsArray.forEach((tool: any) => {
            if (tool && typeof tool === 'object' && 'name' in tool) {
              toolsMap[tool.name] = this.createMCPTool(tool);
            }
          });
        }
      }

      // 缓存结果
      this.toolsCache = toolsMap;
      return toolsMap;
    } catch (error) {
      console.error(`获取工具时出错 (客户端: ${this.clientName}):`, error);
      return {};
    }
  }

  /**
   * 从原始结果中提取工具数组
   */
  private extractToolsArray(rawToolsResult: any): any[] | null {
    if (typeof rawToolsResult === 'object' && 'tools' in rawToolsResult && Array.isArray(rawToolsResult.tools)) {
      return rawToolsResult.tools;
    } else if (Array.isArray(rawToolsResult)) {
      return rawToolsResult;
    }
    return null;
  }

  /**
   * 从原始工具数据创建MCPTool
   */
  private createMCPTool(tool: any): MCPTool {
    const requiredParams = this.extractRequiredParams(tool);
    const enhancedDescription = this.buildEnhancedDescription(
      tool.description,
      requiredParams,
      tool.inputSchema,
      tool.name
    );

    return {
      description: enhancedDescription,
      schema: tool.inputSchema || {},
      execute: async (args: Record<string, unknown>) => {
        return this.callTool(tool.name, args);
      },
    };
  }

  /**
   * 从工具模式中提取必需参数
   */
  private extractRequiredParams(tool: any): string[] {
    if (
      tool.inputSchema &&
      typeof tool.inputSchema === 'object' &&
      'required' in tool.inputSchema &&
      Array.isArray(tool.inputSchema.required)
    ) {
      return tool.inputSchema.required;
    }
    return [];
  }

  /**
   * 构建包含参数信息的增强描述
   */
  private buildEnhancedDescription(
    baseDescription: string,
    requiredParams: string[],
    inputSchema: any,
    toolName: string
  ): string {
    let description = baseDescription || '';

    if (requiredParams.length > 0) {
      description += `\n必需参数: ${requiredParams.join(', ')}`;
    }

    const exampleInput = this.deriveExampleInput(inputSchema, toolName);
    if (exampleInput) {
      description += `\n使用示例: ${exampleInput}`;
    }

    return description;
  }
  /**
   * 关闭与MCP服务器的连接
   * 此方法符合CopilotKit接口的预期
   */
  public async close(): Promise<void> {
    return this.disconnect();
  }

  /**
   * 断开与MCP服务器的连接
   * (旧方法，建议使用close()以兼容CopilotKit)
   */
  public async disconnect(): Promise<void> {
    try {
      // 清理工具缓存
      this.toolsCache = null;
      // 关闭传输连接
      await this.transport.close();
      this.isConnected = false;
      console.log(`已断开与MCP服务器的连接 (客户端: ${this.clientName})`);
    } catch (error) {
      console.error(`断开MCP服务器连接时出错 (客户端: ${this.clientName}):`, error);
      this.onError(error instanceof Error ? error : new Error(String(error)));
    }
  }

  /**
   * 使用给定名称和参数调用工具
   * @param name 工具名称
   * @param args 工具参数
   * @returns 工具执行结果
   */
  public async callTool(name: string, args: Record<string, unknown>): Promise<any> {
    try {
      const fixedArgs = this.normalizeToolArgs(args);
      const processedArgs = this.processStringifiedJsonArgs(fixedArgs);

      // 检查processedArgs中的params对象是否包含对象
      if (processedArgs.params && Object.keys(processedArgs.params).length > 0) {
        return this.client.callTool({
          name: name,
          arguments: processedArgs.params as Record<string, unknown>,
        });
      } else {
        return this.client.callTool({
          name: name,
          arguments: processedArgs,
        });
      }
    } catch (error) {
      console.error(`调用工具 ${name} 时出错 (客户端: ${this.clientName}):`, error);
      throw error;
    }
  }

  /**
   * 标准化工具参数 - 检测并修复LLM工具调用中的常见模式
   * 如双重嵌套的params对象
   */
  private normalizeToolArgs(args: Record<string, unknown>): Record<string, unknown> {
    // 处理双重嵌套参数: { params: { params: { 实际数据 } } }
    if ('params' in args && args.params !== null && typeof args.params === 'object') {
      const paramsObj = args.params as Record<string, unknown>;
      if ('params' in paramsObj) {
        return paramsObj;
      }
    }

    return args;
  }

  /**
   * 处理参数，处理JSON字符串可能被传递而不是对象的情况
   */
  private processStringifiedJsonArgs(args: Record<string, unknown>): Record<string, unknown> {
    const result: Record<string, unknown> = {};

    for (const [key, value] of Object.entries(args)) {
      result[key] = this.processSingleValue(value);
    }

    return result;
  }

  /**
   * 处理单个值的JSON字符串解析
   */
  private processSingleValue(value: unknown): unknown {
    if (typeof value === 'string') {
      return this.parseJsonString(value);
    } else if (Array.isArray(value)) {
      return value.map((item) => this.processSingleValue(item));
    } else if (value !== null && typeof value === 'object') {
      return this.processStringifiedJsonArgs(value as Record<string, unknown>);
    } else {
      return value;
    }
  }

  /**
   * 如果字符串看起来像JSON则解析JSON字符串
   */
  private parseJsonString(value: string): unknown {
    const trimmedValue = value.trim();
    if (trimmedValue.startsWith('{') || trimmedValue.startsWith('[') || trimmedValue.startsWith('"')) {
      try {
        return JSON.parse(value);
      } catch (e) {
        console.error(`JSON解析错误 (客户端: ${this.clientName}):`, e);
        return value;
      }
    }
    return value;
  }

  /**
   * 从工具的inputSchema推导示例输入结构
   * 这有助于LLM理解如何正确格式化请求
   */
  private deriveExampleInput(inputSchema: any, toolName: string): string | null {
    if (!inputSchema) return null;

    try {
      if (toolName.toLowerCase().includes('asana_create')) {
        return '{ "params": { "data": { "name": "任务名称", "notes": "任务描述" } } }';
      }

      if (inputSchema.type === 'object' && inputSchema.properties) {
        const example = this.createExampleObject(inputSchema);
        return JSON.stringify(example, null, 2);
      }

      return null;
    } catch (error) {
      console.error(`创建示例输入时出错 (客户端: ${this.clientName}):`, error);
      return null;
    }
  }

  /**
   * 从对象模式创建示例对象
   */
  private createExampleObject(schema: any): Record<string, any> {
    const result: Record<string, any> = {};

    if (schema.type !== 'object' || !schema.properties) {
      return result;
    }

    const props = schema.properties;

    if (Array.isArray(schema.required)) {
      schema.required.forEach((key: string) => {
        if (key in props) {
          result[key] = this.createExampleValue(props[key]);
        }
      });
    }

    return result;
  }

  /**
   * 根据模式类型创建示例值
   */
  private createExampleValue(propertySchema: any): any {
    if (propertySchema.type === 'object' && propertySchema.properties) {
      return this.createExampleObject(propertySchema);
    } else if (propertySchema.type === 'string') {
      return `示例 ${propertySchema.title || '值'}`;
    } else if (propertySchema.type === 'number') {
      return 123;
    } else if (propertySchema.type === 'boolean') {
      return true;
    } else {
      return null;
    }
  }

  /**
   * 清理工具缓存，强制重新获取工具列表
   * 用于刷新客户端状态时使用
   */
  public clearToolsCache(): void {
    this.toolsCache = null;
    console.log(`已清理工具缓存 (客户端: ${this.clientName})`);
  }

  /**
   * 根据类型创建适当的传输
   */
  private createTransport(
    type: TransportType
  ): SSEClientTransport | StdioClientTransport | StreamableHTTPClientTransport {
    switch (type) {
      case 'stdio':
        return new StdioClientTransport(this.stdioConfig);
      case 'streamable-http':
        return new StreamableHTTPClientTransport(this.serverUrl, this.headers);
      case 'sse':
      default:
        return new SSEClientTransport(this.serverUrl, this.headers);
    }
  }
}
