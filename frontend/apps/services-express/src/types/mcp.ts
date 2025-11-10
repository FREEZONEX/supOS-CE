export type TransportType = 'sse' | 'stdio' | 'streamable-http';

export interface McpClientOptions {
  serverUrl: string;
  transportType?: TransportType;
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
  stdioConfig?: {
    command: string;
    args: string[];
    env: Record<string, string>;
  };
}
