import { useCallback, useEffect, useRef, useState, useMemo } from 'react';

/**
 * SSE连接状态枚举
 */
export enum SSEConnectionStatus {
  CONNECTING = 'connecting',
  OPEN = 'open',
  CLOSED = 'closed',
  ERROR = 'error',
  RECONNECTING = 'reconnecting',
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
  /** 重连间隔时间（毫秒），默认为3000 */
  reconnectInterval?: number;
  /** 最大重连次数，默认为5 */
  maxReconnectAttempts?: number;
  /** 重连延迟增长因子，默认为1.5 */
  reconnectBackoffFactor?: number;
  /** 请求凭据，默认为'same-origin' */
  credentials?: RequestCredentials;
  /** 连接超时时间（毫秒），默认为10000 */
  timeout?: number;
  /** 自定义事件处理器 */
  eventHandlers?: Record<string, (event: SSEEvent) => void>;
  /** 连接成功回调 */
  onOpen?: (event: Event) => void;
  /** 连接关闭回调 */
  onClose?: (event: Event) => void;
  /** 错误处理回调 */
  onError?: (error: Event) => void;
  /** 消息接收回调 */
  onMessage?: (event: SSEEvent) => void;
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
  /** 重连尝试次数 */
  reconnectAttempts: number;
  /** 手动连接方法 */
  connect: () => void;
  /** 手动断开连接方法 */
  disconnect: () => void;
  /** 发送消息方法（如果服务器支持双向通信） */
  send: (data: any, eventType?: string) => void;
  /** 重新连接方法 */
  reconnect: () => void;
  /** 添加事件监听器 */
  addEventListener: (type: string, listener: (event: SSEEvent) => void) => void;
  /** 移除事件监听器 */
  removeEventListener: (type: string, listener: (event: SSEEvent) => void) => void;
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
  const {
    autoConnect = true,
    reconnectInterval = 3000,
    maxReconnectAttempts = 5,
    reconnectBackoffFactor = 1.5,
    credentials = 'same-origin',
    timeout = 10000,
    eventHandlers = {},
    onOpen,
    onClose,
    onError,
    onMessage,
  } = options;

  const eventSourceRef = useRef<EventSource | null>(null);
  const reconnectTimerRef = useRef<NodeJS.Timeout | null>(null);
  const timeoutTimerRef = useRef<NodeJS.Timeout | null>(null);

  const [status, setStatus] = useState<SSEConnectionStatus>(SSEConnectionStatus.CLOSED);
  const [lastMessage, setLastMessage] = useState<SSEEvent | null>(null);
  const [lastError, setLastError] = useState<Event | null>(null);
  const [reconnectAttempts, setReconnectAttempts] = useState(0);

  const eventListenersRef = useRef<Map<string, Set<(event: SSEEvent) => void>>>(new Map());

  /**
   * 清理定时器
   */
  const clearTimers = useCallback(() => {
    if (reconnectTimerRef.current) {
      clearTimeout(reconnectTimerRef.current);
      reconnectTimerRef.current = null;
    }
    if (timeoutTimerRef.current) {
      clearTimeout(timeoutTimerRef.current);
      timeoutTimerRef.current = null;
    }
  }, []);

  /**
   * 断开连接
   */
  const disconnect = useCallback(() => {
    clearTimers();

    if (eventSourceRef.current) {
      eventSourceRef.current.close();
      eventSourceRef.current = null;
    }

    setStatus(SSEConnectionStatus.CLOSED);
    setReconnectAttempts(0);

    if (onClose) {
      onClose(new Event('close'));
    }
  }, [clearTimers, onClose]);

  /**
   * 处理连接错误
   */
  const handleError = useCallback(
    (error: Event) => {
      setLastError(error);
      setStatus(SSEConnectionStatus.ERROR);

      if (onError) {
        onError(error);
      }

      // 自动重连逻辑
      if (reconnectAttempts < maxReconnectAttempts) {
        const delay = reconnectInterval * Math.pow(reconnectBackoffFactor, reconnectAttempts);

        setStatus(SSEConnectionStatus.RECONNECTING);
        reconnectTimerRef.current = setTimeout(() => {
          setReconnectAttempts((prev) => prev + 1);
          connect();
        }, delay);
      }
    },
    [reconnectAttempts, maxReconnectAttempts, reconnectInterval, reconnectBackoffFactor, onError]
  );

  /**
   * 处理消息接收
   */
  const handleMessage = useCallback(
    (event: MessageEvent) => {
      const sseEvent: SSEEvent = {
        type: event.type,
        data: event.data,
        lastEventId: event.lastEventId,
        origin: event.origin,
      };

      setLastMessage(sseEvent);

      // 调用自定义消息处理器
      if (onMessage) {
        onMessage(sseEvent);
      }

      // 调用特定事件类型的处理器
      if (eventHandlers[event.type]) {
        eventHandlers[event.type](sseEvent);
      }

      // 调用注册的事件监听器
      const listeners = eventListenersRef.current.get(event.type);
      if (listeners) {
        listeners.forEach((listener) => listener(sseEvent));
      }
    },
    [onMessage, eventHandlers]
  );

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
        withCredentials: credentials === 'include',
      });

      eventSourceRef.current = eventSource;

      // 设置连接超时
      timeoutTimerRef.current = setTimeout(() => {
        if (eventSource.readyState !== EventSource.OPEN) {
          handleError(new Event('timeout'));
          eventSource.close();
        }
      }, timeout);

      // 连接打开事件
      eventSource.onopen = (event: Event) => {
        clearTimers();
        setStatus(SSEConnectionStatus.OPEN);
        setReconnectAttempts(0);

        if (onOpen) {
          onOpen(event);
        }
      };

      // 错误事件
      eventSource.onerror = (event: Event) => {
        clearTimers();
        handleError(event);
      };

      // 消息事件
      eventSource.onmessage = handleMessage;

      // 注册自定义事件处理器
      Object.keys(eventHandlers).forEach((eventType) => {
        eventSource.addEventListener(eventType, handleMessage as EventListener);
      });
    } catch (error) {
      console.error('Failed to create SSE connection:', error);
      handleError(new Event('connection-error'));
    }
  }, [url, credentials, timeout, eventHandlers, disconnect, clearTimers, handleError, handleMessage, onOpen]);

  /**
   * 重新连接
   */
  const reconnect = useCallback(() => {
    setReconnectAttempts(0);
    connect();
  }, [connect]);

  /**
   * 发送消息（如果服务器支持双向通信）
   */
  const send = useCallback(
    (data: any, eventType: string = 'message') => {
      if (status !== SSEConnectionStatus.OPEN) {
        console.warn('SSE connection is not open');
        return;
      }

      // 注意：标准的SSE是单向的，这里提供的是扩展功能
      // 如果需要双向通信，通常需要配合其他HTTP请求
      console.log('SSE send (not standard):', { eventType, data });

      // 这里可以扩展为发送HTTP请求到服务器
      // fetch(url, { method: 'POST', body: JSON.stringify(data) })
    },
    [status]
  );

  /**
   * 添加事件监听器
   */
  const addEventListener = useCallback(
    (type: string, listener: (event: SSEEvent) => void) => {
      if (!eventListenersRef.current.has(type)) {
        eventListenersRef.current.set(type, new Set());
      }
      eventListenersRef.current.get(type)?.add(listener);

      // 如果连接已建立，立即注册到EventSource
      if (eventSourceRef.current && status === SSEConnectionStatus.OPEN) {
        eventSourceRef.current.addEventListener(type, handleMessage as EventListener);
      }
    },
    [status, handleMessage]
  );

  /**
   * 移除事件监听器
   */
  const removeEventListener = useCallback(
    (type: string, listener: (event: SSEEvent) => void) => {
      const listeners = eventListenersRef.current.get(type);
      if (listeners) {
        listeners.delete(listener);
        if (listeners.size === 0) {
          eventListenersRef.current.delete(type);
        }
      }

      // 如果连接已建立，从EventSource移除监听器
      if (eventSourceRef.current && status === SSEConnectionStatus.OPEN) {
        eventSourceRef.current.removeEventListener(type, handleMessage as EventListener);
      }
    },
    [status, handleMessage]
  );

  // 自动连接效果
  useEffect(() => {
    if (autoConnect) {
      connect();
    }

    return () => {
      disconnect();
    };
  }, [autoConnect, connect, disconnect]);

  // URL变化时重新连接
  useEffect(() => {
    if (autoConnect && url) {
      reconnect();
    }
  }, [url, autoConnect, reconnect]);

  const returnValue = useMemo(
    () => ({
      status,
      lastMessage,
      lastError,
      reconnectAttempts,
      connect,
      disconnect,
      send,
      reconnect,
      addEventListener,
      removeEventListener,
    }),
    [
      status,
      lastMessage,
      lastError,
      reconnectAttempts,
      connect,
      disconnect,
      send,
      reconnect,
      addEventListener,
      removeEventListener,
    ]
  );

  return returnValue;
};

export default useSSE;
