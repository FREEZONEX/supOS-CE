export type TransportType = 'sse' | 'stdio' | 'streamable-http';

export interface McpClientOptions {
  serverUrl: string;
  transportType?: TransportType;
  clientName: string;
  headers?: Record<string, string>;
  onMessage?: (message: Record<string, unknown>) => void;
  onError?: (error: Error) => void;
  onOpen?: () => void;
  onClose?: () => void;
  stdioConfig?: any;
}

export interface TransportConfig {
  transportType: TransportType;
  serverUrl: string;
  clientName: string;
  stdioConfig?: {
    command: string;
    args: string[];
    env: Record<string, string>;
  };
}

export interface StdioConfig {
  command: string;
  args: string[];
  env?: Record<string, string>[];
}

export interface HttpConfig {
  url?: string;
}

export type TTransportConfig = StdioConfig | HttpConfig;

export interface McpServerConfig {
  name: string;
  transportType: TransportType;
  config: TTransportConfig;
}
