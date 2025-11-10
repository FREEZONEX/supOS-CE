import { MCPTool, MCPClient as MCPClientInterface } from '@copilotkit/runtime';
import { Client } from '@modelcontextprotocol/sdk/client/index.js';
import { SSEClientTransport } from '@modelcontextprotocol/sdk/client/sse.js';
import { StdioClientTransport } from '@modelcontextprotocol/sdk/client/stdio.js';
import { StreamableHTTPClientTransport } from '@modelcontextprotocol/sdk/client/streamableHttp.js';
import { JSONRPCMessage } from '@modelcontextprotocol/sdk/types';
import { McpClientOptions, TransportType } from './types';

export class MCPClient implements MCPClientInterface {
  private client: Client;
  private transport: SSEClientTransport | StdioClientTransport | StreamableHTTPClientTransport;
  private transportType: TransportType;
  private serverUrl: URL;
  private onMessage: (message: Record<string, unknown>) => void;
  private onError: (error: Error) => void;
  private onOpen: () => void;
  private onClose: () => void;
  private isConnected = false;
  private headers?: Record<string, string>;
  private stdioConfig?: any;

  private toolsCache: Record<string, MCPTool> | null = null;

  constructor(options: McpClientOptions) {
    this.serverUrl =
      options.transportType === 'stdio' ? new URL('http://localhost:3000') : new URL(options.serverUrl || '');
    this.transportType = options.transportType || 'stdio';
    this.headers = options.headers || {};
    this.stdioConfig = options.stdioConfig || {};
    this.onMessage = options.onMessage || ((message) => console.log('Message received:', message));
    this.onError = options.onError || ((error) => console.error('Error:', error));
    this.onOpen = options.onOpen || (() => console.log('Connection opened'));
    this.onClose = options.onClose || (() => console.log('Connection closed'));
    this.transport = this.createTransport(this.transportType);

    this.client = new Client({
      name: 'cpk-mcp-client',
      version: '0.0.1',
    });

    this.transport.onmessage = this.handleMessage.bind(this);
    this.transport.onerror = this.handleError.bind(this);
    this.transport.onclose = this.handleClose.bind(this);
  }

  private handleMessage(message: JSONRPCMessage): void {
    try {
      this.onMessage(message as Record<string, unknown>);
    } catch (error) {
      this.onError(error instanceof Error ? error : new Error(`Failed to handle message: ${error}`));
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

  public async connect(): Promise<void> {
    try {
      console.log('Connecting to MCP server:', this.serverUrl.href);
      await this.client.connect(this.transport as any);
      this.isConnected = true;
      console.log('Connected to MCP server');
      this.onOpen();
    } catch (error) {
      console.error('Failed to connect to MCP server:', error);
      this.onError(error instanceof Error ? error : new Error(String(error)));
      throw error;
    }
  }

  public async tools(): Promise<Record<string, MCPTool>> {
    try {
      if (this.toolsCache) {
        console.log('有tools');
        return this.toolsCache;
      }
      console.log('重加tools');
      const rawToolsResult = await this.client.listTools();

      const toolsMap: Record<string, MCPTool> = {};

      if (rawToolsResult) {
        if (typeof rawToolsResult === 'object' && 'tools' in rawToolsResult && Array.isArray(rawToolsResult.tools)) {
          rawToolsResult.tools.forEach((tool: any) => {
            if (tool && typeof tool === 'object' && 'name' in tool) {
              let requiredParams: string[] = [];

              if (
                tool.inputSchema &&
                typeof tool.inputSchema === 'object' &&
                'required' in tool.inputSchema &&
                Array.isArray(tool.inputSchema.required)
              ) {
                requiredParams = tool.inputSchema.required;
              }

              let enhancedDescription = tool.description || '';

              if (requiredParams.length > 0) {
                enhancedDescription += `\nRequired parameters: ${requiredParams.join(', ')}`;
              }

              const exampleInput = this.deriveExampleInput(tool.inputSchema, tool.name);
              if (exampleInput) {
                enhancedDescription += `\nExample usage: ${exampleInput}`;
              }

              toolsMap[tool.name] = {
                description: enhancedDescription,
                schema: tool.inputSchema || {},
                execute: async (args: Record<string, unknown>) => {
                  return this.callTool(tool.name, args);
                },
              };
            }
          });
        } else if (Array.isArray(rawToolsResult)) {
          rawToolsResult.forEach((tool: any) => {
            if (tool && typeof tool === 'object' && 'name' in tool) {
              let requiredParams: string[] = [];

              if (
                tool.inputSchema &&
                typeof tool.inputSchema === 'object' &&
                'required' in tool.inputSchema &&
                Array.isArray(tool.inputSchema.required)
              ) {
                requiredParams = tool.inputSchema.required;
              }

              let enhancedDescription = tool.description || '';

              if (requiredParams.length > 0) {
                enhancedDescription += `\nRequired parameters: ${requiredParams.join(', ')}`;
              }

              const exampleInput = this.deriveExampleInput(tool.inputSchema, tool.name);
              if (exampleInput) {
                enhancedDescription += `\nExample usage: ${exampleInput}`;
              }

              toolsMap[tool.name] = {
                description: enhancedDescription,
                schema: tool.inputSchema || {},
                execute: async (args: Record<string, unknown>) => {
                  return this.callTool(tool.name, args);
                },
              };
            }
          });
        }
      }

      this.toolsCache = toolsMap;
      console.log('toolsMap:', toolsMap);
      return toolsMap;
    } catch (error) {
      console.error('Error fetching tools:', error);
      return {};
    }
  }

  public async close(): Promise<void> {
    return this.disconnect();
  }

  public async disconnect(): Promise<void> {
    try {
      this.toolsCache = null;
      await this.transport.close();
      this.isConnected = false;
      console.log('Disconnected from MCP server');
    } catch (error) {
      console.error('Error disconnecting from MCP server:', error);
      this.onError(error instanceof Error ? error : new Error(String(error)));
    }
  }

  public async callTool(name: string, args: Record<string, unknown>): Promise<any> {
    try {
      console.log(`Calling tool: ${name} with args:`, JSON.stringify(args, null, 2));

      const fixedArgs = this.normalizeToolArgs(args);
      const processedArgs = this.processStringifiedJsonArgs(fixedArgs);

      console.log(`Processed args for ${name}:`, JSON.stringify(processedArgs.params, null, 2));

      if (processedArgs.params && Object.keys(processedArgs.params).length > 0) {
        console.log(`Calling tool ${name} with processed params:`, JSON.stringify(processedArgs.params, null, 2));
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
      console.error(`Error calling tool ${name}:`, error);
      throw error;
    }
  }

  private normalizeToolArgs(args: Record<string, unknown>): Record<string, unknown> {
    if ('params' in args && args.params !== null && typeof args.params === 'object') {
      const paramsObj = args.params as Record<string, unknown>;
      if ('params' in paramsObj) {
        console.log('Detected double-nested params, fixing structure');
        return paramsObj;
      }
    }
    return args;
  }

  private processStringifiedJsonArgs(args: Record<string, unknown>): Record<string, unknown> {
    const result: Record<string, unknown> = {};

    for (const [key, value] of Object.entries(args)) {
      if (typeof value === 'string') {
        try {
          const trimmedValue = value.trim();
          if (trimmedValue.startsWith('{') || trimmedValue.startsWith('[') || trimmedValue.startsWith('\'"')) {
            const parsedValue = JSON.parse(value);
            result[key] = parsedValue;
          } else {
            result[key] = value;
          }
        } catch (e) {
          console.log('JSON parsing error for key:', key, 'value:', value, 'error:', e);
          result[key] = value;
        }
      } else if (Array.isArray(value)) {
        result[key] = value.map((item) =>
          typeof item === 'object' && item !== null
            ? this.processStringifiedJsonArgs(item as Record<string, unknown>)
            : item
        );
      } else if (value !== null && typeof value === 'object') {
        result[key] = this.processStringifiedJsonArgs(value as Record<string, unknown>);
      } else {
        result[key] = value;
      }
    }
    return result;
  }

  private deriveExampleInput(inputSchema: any, toolName: string): string | null {
    if (!inputSchema) return null;

    try {
      if (toolName.toLowerCase().includes('asana_create')) {
        return '{ "params": { "data": { "name": "Task name", "notes": "Task description" } } }';
      }

      if (inputSchema.type === 'object' && inputSchema.properties) {
        const example: Record<string, any> = {};
        const props = inputSchema.properties;

        if (Array.isArray(inputSchema.required)) {
          inputSchema.required.forEach((key: string) => {
            if (key in props) {
              if (props[key].type === 'object' && props[key].properties) {
                example[key] = this.createExampleObject(props[key]);
              } else if (props[key].type === 'string') {
                example[key] = `"Example ${key}"`;
              } else if (props[key].type === 'number') {
                example[key] = 123;
              } else if (props[key].type === 'boolean') {
                example[key] = true;
              } else {
                example[key] = null;
              }
            }
          });
        }
        return JSON.stringify(example, null, 2);
      }
      return null;
    } catch (error) {
      console.error('Error creating example input:', error);
      return null;
    }
  }

  private createExampleObject(schema: any): Record<string, any> {
    const result: Record<string, any> = {};

    if (schema.type !== 'object' || !schema.properties) {
      return result;
    }

    const props = schema.properties;

    if (Array.isArray(schema.required)) {
      schema.required.forEach((key: string) => {
        if (key in props) {
          if (props[key].type === 'object' && props[key].properties) {
            result[key] = this.createExampleObject(props[key]);
          } else if (props[key].type === 'string') {
            result[key] = `Example ${key}`;
          } else if (props[key].type === 'number') {
            result[key] = 123;
          } else if (props[key].type === 'boolean') {
            result[key] = true;
          } else {
            result[key] = null;
          }
        }
      });
    }
    return result;
  }

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
