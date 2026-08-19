import { Badge, Spin, Tag, message as antdMessage } from 'antd';
import { Box, Copy } from '@carbon/icons-react';
import { useEffect, useMemo, useRef, useState } from 'react';
import { useSearchParams } from 'react-router';
import type { MessageInstance } from 'antd/es/message/interface';
import { fetchQrConfig, parseBinding, type QrConfig } from '../../api/instances';
import { useSuspended } from '../../bridge/use-suspended';
import { parsePlacement, parseSceneConfig, type SceneQrConfig } from '../../api/scenes';
import { t } from '../../i18n';
import { useMqttTopic, useMqttTopics, type TopicMessage } from '../../mqtt/use-mqtt';
import { getModelFormat } from '../../utils/thumbnail';
import { useSplatViewer } from '../../viewer/splat-viewer';
import { useModelViewer } from '../../viewer/use-model-viewer';
import { useSceneEditor } from '../../viewer/use-scene-editor';
import { TagOverlay } from '../instance-detail/tag-overlay';

// 扫码分享 viewer（对齐云端 tier0-anchor 移动端页面）：
// 免登录，?configUrl= 指向实例或场景的 qr-config（按 payload 形状区分），
// 上半屏 3D 实时驱动，下半部信息卡 + Payload 实时表；场景模式支持 ?focusInstanceId= 聚焦。
type ViewerConfig = QrConfig | SceneQrConfig;

const isSceneConfig = (config: ViewerConfig): config is SceneQrConfig => 'scene' in config;

function statusBadge(running: boolean, connected: boolean) {
  return (
    <Badge
      status={running ? 'success' : connected ? 'processing' : 'default'}
      text={
        <span className="text-xs" style={{ color: 'var(--ui-description-card-color)' }}>
          {running
            ? t('instance.status.running')
            : connected
              ? t('instance.mqtt.waiting')
              : t('instance.mqtt.disconnected')}
        </span>
      }
    />
  );
}

function TopicChip({ topic, messageApi }: { topic: string; messageApi: MessageInstance }) {
  return (
    <button
      type="button"
      className="flex max-w-full cursor-pointer items-center gap-2 rounded px-2 py-1 text-xs"
      style={{
        background: 'var(--ui-primary-bg)',
        color: 'var(--ui-description-card-color)',
        border: '1px solid var(--ui-header-splitter-color, #e0e0e0)',
      }}
      onClick={() => {
        navigator.clipboard
          .writeText(topic)
          .then(() => messageApi.success(t('viewer.qr.topicCopied')))
          .catch(() => messageApi.error(t('viewer.qr.copyFailed')));
      }}
    >
      <span className="truncate">{topic}</span>
      <Copy size={12} className="shrink-0" />
    </button>
  );
}

function PayloadTable({ message, connected }: { message: TopicMessage | null; connected: boolean }) {
  return (
    <div className="overflow-hidden rounded" style={{ border: '1px solid var(--ui-header-splitter-color, #e0e0e0)' }}>
      <div className="px-4 py-2 text-sm font-semibold" style={{ color: 'var(--ui-text-color)' }}>
        Payload
      </div>
      <table className="w-full text-sm">
        <thead>
          <tr className="text-xs" style={{ color: 'var(--ui-description-card-color)' }}>
            <th className="px-4 py-2 text-left font-medium">{t('viewer.qr.key')}</th>
            <th className="px-4 py-2 text-left font-medium">{t('viewer.qr.value')}</th>
          </tr>
        </thead>
        <tbody>
          {message && Object.keys(message.payload).length > 0 ? (
            Object.entries(message.payload).map(([key, value]) => (
              <tr key={key} style={{ borderTop: '1px solid var(--ui-header-splitter-color, #e0e0e0)' }}>
                <td className="px-4 py-2 text-xs font-semibold" style={{ color: 'var(--ui-description-card-color)' }}>
                  {key}
                </td>
                <td className="px-4 py-2 text-xs font-semibold" style={{ color: 'var(--ui-text-color)' }}>
                  {typeof value === 'object' ? JSON.stringify(value) : String(value)}
                </td>
              </tr>
            ))
          ) : (
            <tr style={{ borderTop: '1px solid var(--ui-header-splitter-color, #e0e0e0)' }}>
              <td
                colSpan={2}
                className="px-4 py-4 text-center text-xs"
                style={{ color: 'var(--ui-description-card-color)' }}
              >
                {connected ? t('instance.mqtt.waiting') : t('instance.mqtt.disconnected')}
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}

function ViewerShell({
  badge,
  viewport,
  children,
}: {
  badge: React.ReactNode;
  viewport: React.ReactNode;
  children: React.ReactNode;
}) {
  return (
    <div className="flex min-h-screen flex-col" style={{ background: 'var(--ui-bg-color)' }}>
      <div
        className="sticky top-0 z-50 flex items-center justify-between px-4 py-3"
        style={{
          background: 'var(--ui-primary-bg)',
          borderBottom: '1px solid var(--ui-header-splitter-color, #e0e0e0)',
        }}
      >
        <div className="text-lg font-medium" style={{ color: 'var(--ui-text-color)' }}>
          Tier0
        </div>
        {badge}
      </div>
      <div className="relative h-[52vh] w-full shrink-0">{viewport}</div>
      {children}
    </div>
  );
}

function InstanceQrViewer({ config, messageApi }: { config: QrConfig; messageApi: MessageInstance }) {
  const binding = useMemo(() => parseBinding(config.bindingJson), [config.bindingJson]);
  // SPLAT 模型走高斯泼溅查看器（静态标点，无运动映射）
  const isSplat = getModelFormat(config.model.originFile || config.model.name) === 'splat';
  const { containerRef, viewer } = useModelViewer(!isSplat ? config.model.fileUrl : undefined);
  const { containerRef: splatContainerRef, viewer: splatViewer } = useSplatViewer(
    isSplat ? config.model.fileUrl : undefined
  );
  const overlayViewer = isSplat ? splatViewer : viewer;
  // 无宿主直开：仅 visibilitychange 生效（手机切后台停渲染/断订阅）
  const suspended = useSuspended();
  useEffect(() => {
    viewer.setSuspended(suspended);
    splatViewer.setSuspended(suspended);
  }, [suspended, viewer, splatViewer]);
  const credentials = useMemo(
    () => ({ username: config.mqtt.username, password: config.mqtt.password, clientId: config.mqtt.clientId }),
    [config.mqtt.username, config.mqtt.password, config.mqtt.clientId]
  );
  const { message, connected } = useMqttTopic(config.mqtt.topic || undefined, credentials, suspended);

  useEffect(() => {
    if (viewer.status !== 'ready') return;
    viewer.setBindings(binding.motionMappings);
  }, [viewer, viewer.status, binding]);

  useEffect(() => {
    if (!message) return;
    viewer.applyPayload(message.payload);
  }, [message, viewer]);

  return (
    <ViewerShell
      badge={statusBadge(connected && Boolean(message), connected)}
      viewport={
        <>
          <div ref={containerRef} className="h-full w-full" style={isSplat ? { display: 'none' } : undefined} />
          {isSplat ? <div ref={splatContainerRef} className="h-full w-full" /> : null}
          <TagOverlay
            tags={binding.dataTags}
            viewer={overlayViewer}
            payload={message?.payload}
            messageTs={message?.ts}
          />
          {overlayViewer.status !== 'ready' ? (
            <div
              className="absolute inset-0 flex items-center justify-center text-sm"
              style={{ background: 'var(--ui-bg-color)', color: 'var(--ui-description-card-color)' }}
            >
              {overlayViewer.status === 'error'
                ? `${t('model.viewer.error')}: ${overlayViewer.error}`
                : t('model.viewer.loading')}
            </div>
          ) : null}
        </>
      }
    >
      <div className="px-4 py-4">
        <div className="mb-2 flex items-center gap-2">
          <Box size={20} style={{ color: 'var(--ui-text-color)' }} />
          <div className="font-medium" style={{ color: 'var(--ui-text-color)' }}>
            {config.instance.name}
          </div>
        </div>
        <Tag bordered={false} color="lime">
          {config.model.name}
        </Tag>
      </div>
      <div className="px-4 pb-3">
        <TopicChip topic={config.mqtt.topic} messageApi={messageApi} />
      </div>
      <div className="flex-1 px-4 pb-6">
        <PayloadTable message={message} connected={connected} />
      </div>
    </ViewerShell>
  );
}

function SceneQrViewer({
  config,
  focusInstanceId,
  messageApi,
}: {
  config: SceneQrConfig;
  focusInstanceId: string;
  messageApi: MessageInstance;
}) {
  const editor = useSceneEditor({ interactive: false });
  // 无宿主直开：仅 visibilitychange 生效
  const suspended = useSuspended();
  useEffect(() => {
    editor.setSuspended(suspended);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [suspended, editor.setSuspended]);
  const [loadedCount, setLoadedCount] = useState(0);
  const focusItem = useMemo(
    () => config.scene.items.find((item) => String(item.instanceId) === focusInstanceId) ?? config.scene.items[0],
    [config.scene.items, focusInstanceId]
  );
  const [focusMessage, setFocusMessage] = useState<TopicMessage | null>(null);
  const focusTopicRef = useRef(focusItem?.instance.topic || '');
  focusTopicRef.current = focusItem?.instance.topic || '';

  // 装载场景：视口/光照配置 + 全部条目（模型 + 摆放 + 运动绑定），完成后聚焦
  useEffect(() => {
    if (!editor.ready) return;
    editor.applyConfig(parseSceneConfig(config.scene.configJson));
    let cancelled = false;
    Promise.all(
      config.scene.items.map((item) =>
        editor.addItem({
          key: String(item.instanceId),
          instanceId: item.instanceId,
          fileUrl: item.instance.modelFileUrl,
          instanceHeight: item.instance.height,
          placement: parsePlacement(item.placementJson),
          motionMappings: parseBinding(item.instance.bindingJson).motionMappings,
        })
      )
    ).then(() => {
      if (cancelled) return;
      setLoadedCount(config.scene.items.length);
      if (focusInstanceId && config.scene.items.some((item) => String(item.instanceId) === focusInstanceId)) {
        editor.frameItem(focusInstanceId);
      } else {
        editor.frameAll();
      }
    });
    return () => {
      cancelled = true;
    };
    // editor 引用稳定（hook 内部 useCallback），仅在 ready / 配置变化时装载
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [editor.ready, config, focusInstanceId]);

  const topics = useMemo(
    () => config.scene.items.map((item) => item.instance.topic).filter(Boolean),
    [config.scene.items]
  );
  const credentials = useMemo(
    () => ({ username: config.mqtt.username, password: config.mqtt.password, clientId: config.mqtt.clientId }),
    [config.mqtt.username, config.mqtt.password, config.mqtt.clientId]
  );
  const { connected } = useMqttTopics(
    topics,
    (topic, payload) => {
      for (const item of config.scene.items) {
        if (item.instance.topic === topic) editor.applyPayloadToItem(String(item.instanceId), payload);
      }
      if (topic === focusTopicRef.current) setFocusMessage({ payload, raw: '', ts: Date.now() });
    },
    credentials
  );

  return (
    <ViewerShell
      badge={statusBadge(connected && Boolean(focusMessage), connected)}
      viewport={
        <>
          <div ref={editor.containerRef} className="h-full w-full" />
          {loadedCount < config.scene.items.length ? (
            <div
              className="absolute inset-0 flex items-center justify-center text-sm"
              style={{ background: 'var(--ui-bg-color)', color: 'var(--ui-description-card-color)' }}
            >
              {t('model.viewer.loading')}
            </div>
          ) : null}
        </>
      }
    >
      <div className="px-4 py-4">
        <div className="mb-2 flex items-center gap-2">
          <Box size={20} style={{ color: 'var(--ui-text-color)' }} />
          <div className="font-medium" style={{ color: 'var(--ui-text-color)' }}>
            {config.scene.name}
          </div>
        </div>
        <div className="flex flex-wrap gap-1">
          {config.scene.items.map((item) => (
            <Tag
              key={item.instanceId}
              bordered={false}
              color={focusItem && item.instanceId === focusItem.instanceId ? 'lime' : undefined}
            >
              {item.instance.name}
            </Tag>
          ))}
        </div>
      </div>
      {focusItem?.instance.topic ? (
        <div className="px-4 pb-3">
          <TopicChip topic={focusItem.instance.topic} messageApi={messageApi} />
        </div>
      ) : null}
      <div className="flex-1 px-4 pb-6">
        <PayloadTable message={focusMessage} connected={connected} />
      </div>
    </ViewerShell>
  );
}

export default function QrViewerPage() {
  const [params] = useSearchParams();
  const configUrl = params.get('configUrl') || '';
  const focusInstanceId = params.get('focusInstanceId') || '';
  const [config, setConfig] = useState<ViewerConfig | null>(null);
  // 缺 configUrl 属静态错误，直接从初始值派生；fetch 失败在回调里覆盖
  const [error, setError] = useState(() => (configUrl ? '' : t('viewer.qr.noConfig')));
  const [messageApi, contextHolder] = antdMessage.useMessage();

  useEffect(() => {
    if (!configUrl) return;
    fetchQrConfig(configUrl)
      .then((data) => setConfig(data as ViewerConfig))
      .catch((e) => setError((e as Error).message));
  }, [configUrl]);

  if (error) {
    return (
      <div
        className="flex min-h-screen items-center justify-center p-6 text-center text-sm"
        style={{ color: 'var(--ui-description-card-color)' }}
      >
        {error}
      </div>
    );
  }
  if (!config) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Spin />
      </div>
    );
  }

  return (
    <>
      {contextHolder}
      {isSceneConfig(config) ? (
        <SceneQrViewer config={config} focusInstanceId={focusInstanceId} messageApi={messageApi} />
      ) : (
        <InstanceQrViewer config={config} messageApi={messageApi} />
      )}
    </>
  );
}
