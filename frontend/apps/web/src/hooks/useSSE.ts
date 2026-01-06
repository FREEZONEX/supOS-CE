import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

/**
 * SSE连接状态枚举
 */
export enum SSEConnectionStatus {
  CONNECTING = 'connecting',
  OPEN = 'open',
  CLOSED = 'closed',
  ERROR = 'error',
}

/**
 * SSE事件类型
 */
export interface SSEEvent {
  type: string;
  data: any;
  lastEventId?: string;
  origin?: string;
}

/**
 * SSE Hook配置选项
 */
export interface UseServerSentEventsOptions {
  /** 是否自动连接，默认为true */
  autoConnect?: boolean;
  /** 请求凭据，默认为'same-origin' */
  credentials?: RequestCredentials;
  /** 连接关闭回调 */
  onClose?: (event: Event) => void;
  /** 错误处理回调 */
  onError?: (error: Event) => void;
  /** 消息接收回调 */
  onMessage?: (event: SSEEvent) => void;
  /** 连接打开回调 */
  onOpen?: (event: Event) => void;
}

/**
 * SSE Hook返回值类型
 */
export interface UseServerSentEventsReturn {
  /** 当前连接状态 */
  status: SSEConnectionStatus;
  /** 最后接收到的消息 */
  lastMessage: SSEEvent | null;
  /** 最后发生的错误 */
  lastError: Event | null;
  /** 手动连接方法 */
  connect: () => void;
  /** 手动断开连接方法 */
  disconnect: () => void;
}

/**
 * 使用Server-Sent Events的React Hook
 * 提供与WebSocket类似的API设计和使用体验
 *
 * @param url SSE服务器URL
 * @param options 配置选项
 * @returns SSE Hook返回值
 */
const useSSE = (url: string, options: UseServerSentEventsOptions = {}): UseServerSentEventsReturn => {
  const { autoConnect = true, credentials = 'same-origin', onClose, onError, onMessage, onOpen } = options;

  const eventSourceRef = useRef<EventSource | null>(null);
  const callbacksRef = useRef({ onClose, onError, onMessage, onOpen });

  const [status, setStatus] = useState<SSEConnectionStatus>(SSEConnectionStatus.CLOSED);
  const [lastMessage, setLastMessage] = useState<SSEEvent | null>(null);
  const [lastError, setLastError] = useState<Event | null>(null);

  // 更新回调函数引用
  useEffect(() => {
    callbacksRef.current = { onClose, onError, onMessage, onOpen };
  }, [onClose, onError, onMessage, onOpen]);

  /**
   * 断开连接
   */
  const disconnect = useCallback(() => {
    if (eventSourceRef.current) {
      eventSourceRef.current.close();
      eventSourceRef.current = null;
    }

    setStatus(SSEConnectionStatus.CLOSED);

    if (callbacksRef.current.onClose) {
      callbacksRef.current.onClose(new Event('close'));
    }
  }, []);

  /**
   * 处理消息接收
   */
  const handleMessage = useCallback((event: MessageEvent) => {
    const sseEvent: SSEEvent = {
      type: event.type,
      data: event.data,
      lastEventId: event.lastEventId,
      origin: event.origin,
    };

    setLastMessage(sseEvent);

    // 调用自定义消息处理器
    if (callbacksRef.current.onMessage) {
      callbacksRef.current.onMessage(sseEvent);
    }
  }, []);

  /**
   * 建立连接
   */
  const connect = useCallback(() => {
    // 清理现有连接
    disconnect();

    if (!url) {
      console.warn('SSE URL is required');
      return;
    }

    try {
      setStatus(SSEConnectionStatus.CONNECTING);

      // 创建EventSource实例
      const eventSource = new EventSource(url, {
        withCredentials: true,
      });
      console.log('SSE连接建立:', url);

      eventSourceRef.current = eventSource;

      // 连接打开事件
      eventSource.onopen = (event: Event) => {
        setStatus(SSEConnectionStatus.OPEN);
        if (callbacksRef.current.onOpen) {
          callbacksRef.current.onOpen(event);
        }
      };

      // 错误事件
      eventSource.onerror = (event: Event) => {
        setLastError(event);
        setStatus(SSEConnectionStatus.ERROR);

        if (callbacksRef.current.onError) {
          callbacksRef.current.onError(event);
        }
      };

      // 消息事件
      eventSource.onmessage = handleMessage;
    } catch (error) {
      console.error('Failed to create SSE connection:', error);
      setLastError(new Event('connection-error'));
      setStatus(SSEConnectionStatus.ERROR);
    }
  }, [url, credentials, disconnect, handleMessage]);

  // 自动连接效果
  useEffect(() => {
    if (autoConnect) {
      connect();
    }

    return () => {
      disconnect();
    };
  }, [autoConnect]);

  // URL变化时重新连接
  useEffect(() => {
    if (autoConnect && url) {
      connect();
    }
  }, [url, autoConnect, connect]);

  return useMemo(
    () => ({
      status,
      lastMessage,
      lastError,
      connect,
      disconnect,
    }),
    [status, lastMessage, lastError, connect, disconnect]
  );
};

export default useSSE;
