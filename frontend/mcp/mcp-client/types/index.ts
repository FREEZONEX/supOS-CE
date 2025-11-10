export type TransportType = 'stdio' | 'sse' | 'streamable-http';

export interface McpClientOptions {
  transportType: TransportType;
  serverUrl?: string;
  headers?: Record<string, string>;
  stdioConfig?: any;
  onMessage?: (message: Record<string, unknown>) => void;
  onError?: (error: Error) => void;
  onOpen?: () => void;
  onClose?: () => void;
}
