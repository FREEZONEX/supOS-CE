import mqtt, { type MqttClient } from 'mqtt';
import { useEffect, useRef, useState } from 'react';
import { getMqttCredentials } from '../api/instances';

export type TopicMessage = { payload: Record<string, unknown>; raw: string; ts: number };

export interface MqttCredentialInput {
  username: string;
  password: string;
  clientId: string;
}

// 订阅单个 topic 的最新消息（EMQX WebSocket :8083）。
// 凭据默认来自 /api/core/anchor/mqtt-credentials（需登录）；
// 免登录场景（扫码 viewer）可直接传入 credentials（来自 qr-config）；
// suspended=true 时断开连接（页签隐藏省资源），恢复后自动重连重订阅。
export function useMqttTopic(topic: string | undefined, credentials?: MqttCredentialInput, suspended?: boolean) {
  const [message, setMessage] = useState<TopicMessage | null>(null);
  const [connected, setConnected] = useState(false);
  const clientRef = useRef<MqttClient | null>(null);
  const credentialsKey = credentials ? `${credentials.username}\n${credentials.clientId}` : '';

  useEffect(() => {
    if (!topic || suspended) return;
    let disposed = false;
    let client: MqttClient | null = null;

    (credentials ? Promise.resolve(credentials) : getMqttCredentials())
      .then((cred) => {
        if (disposed) return;
        // 默认 WS 端口 8083（deploy compose OS_MQTT_WEBSOCKET_PORT）；https 环境走 8084 wss
        const secure = window.location.protocol === 'https:';
        const port = secure ? 8084 : 8083;
        const url = `${secure ? 'wss' : 'ws'}://${window.location.hostname}:${port}/mqtt`;
        client = mqtt.connect(url, {
          username: cred.username,
          password: cred.password,
          // 后端按 base clientId + 随机后缀校验；每个连接用唯一后缀，多个页面并发预览互不挤占
          clientId: `${cred.clientId}-${Math.random().toString(36).slice(2, 10)}`,
          clean: true,
          connectTimeout: 5000,
          reconnectPeriod: 3000,
        });
        clientRef.current = client;
        client.on('connect', () => {
          setConnected(true);
          client?.subscribe(topic, (err) => {
            if (err) console.error('[anchor-mqtt] subscribe failed:', err);
          });
        });
        client.on('close', () => setConnected(false));
        client.on('message', (_t, raw) => {
          const text = raw.toString();
          try {
            const payload = JSON.parse(text) as Record<string, unknown>;
            setMessage({ payload, raw: text, ts: Date.now() });
          } catch {
            setMessage({ payload: {}, raw: text, ts: Date.now() });
          }
        });
      })
      .catch((e) => console.error('[anchor-mqtt] credentials failed:', e));

    return () => {
      disposed = true;
      client?.end(true);
      clientRef.current = null;
      setConnected(false);
      setMessage(null);
    };
    // credentials 对象按内容 key 比较，避免调用方每次渲染传新对象导致重连
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [topic, credentialsKey, suspended]);

  return { message, connected };
}

// 多 topic 订阅（场景运行时）：单客户端订阅一组 topic，消息经回调直接派发（避免高频 setState）。
// 凭据默认来自登录接口；免登录场景（扫码 viewer）可直接传入 credentials。
export function useMqttTopics(
  topics: string[],
  onMessage: (topic: string, payload: Record<string, unknown>) => void,
  credentials?: MqttCredentialInput,
  suspended?: boolean
) {
  const [connected, setConnected] = useState(false);
  const handlerRef = useRef(onMessage);
  handlerRef.current = onMessage;
  const topicsKey = [...new Set(topics)].sort().join('\n');
  const credentialsKey = credentials ? `${credentials.username}\n${credentials.clientId}` : '';

  useEffect(() => {
    const uniqueTopics = topicsKey ? topicsKey.split('\n') : [];
    if (uniqueTopics.length === 0 || suspended) return;
    let disposed = false;
    let client: MqttClient | null = null;

    (credentials ? Promise.resolve(credentials) : getMqttCredentials())
      .then((cred) => {
        if (disposed) return;
        const secure = window.location.protocol === 'https:';
        const port = secure ? 8084 : 8083;
        const url = `${secure ? 'wss' : 'ws'}://${window.location.hostname}:${port}/mqtt`;
        client = mqtt.connect(url, {
          username: cred.username,
          password: cred.password,
          clientId: `${cred.clientId}-${Math.random().toString(36).slice(2, 10)}`,
          clean: true,
          connectTimeout: 5000,
          reconnectPeriod: 3000,
        });
        client.on('connect', () => {
          setConnected(true);
          uniqueTopics.forEach((topic) => client?.subscribe(topic, () => {}));
        });
        client.on('close', () => setConnected(false));
        client.on('message', (topic, raw) => {
          try {
            handlerRef.current(topic, JSON.parse(raw.toString()) as Record<string, unknown>);
          } catch {
            /* 非 JSON 消息忽略 */
          }
        });
      })
      .catch((e) => console.error('[anchor-mqtt] credentials failed:', e));

    return () => {
      disposed = true;
      client?.end(true);
      setConnected(false);
    };
    // credentials 对象按内容 key 比较，避免调用方每次渲染传新对象导致重连
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [topicsKey, credentialsKey, suspended]);

  return { connected };
}

// UNS payload 里的字段值可能嵌套（{"J1_Cur": 1.2} 或 {"fields": {...}}），统一取顶层数字
export function payloadNumber(payload: Record<string, unknown> | undefined, key: string): number | undefined {
  if (!payload || !key) return undefined;
  const value = payload[key];
  const num = typeof value === 'string' ? Number(value) : (value as number);
  return typeof num === 'number' && Number.isFinite(num) ? num : undefined;
}

export function payloadKeys(payload: Record<string, unknown> | undefined): string[] {
  if (!payload) return [];
  return Object.keys(payload).filter((key) => {
    const value = payload[key];
    return typeof value === 'number' || (typeof value === 'string' && value !== '' && Number.isFinite(Number(value)));
  });
}
