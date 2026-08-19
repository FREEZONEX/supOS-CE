import {
  CenterCircle,
  ChevronLeft,
  Download,
  Edit,
  Link as LinkIcon,
  Maximize,
  QrCode,
  Upload,
} from '@carbon/icons-react';
import { App, Badge, Button, Dropdown, InputNumber, Modal, Segmented, Tooltip, Tree } from 'antd';
import { EllipsisVertical } from 'lucide-react';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useNavigate, useParams } from 'react-router';
import {
  getInstance,
  parseBinding,
  updateInstance,
  type InstanceBindingPayload,
  type InstanceInfo,
} from '../../api/instances';
import { useSuspended } from '../../bridge/use-suspended';
import { ResizableSplit } from '../../components/resizable-split';
import { TopicSelect } from '../../components/topic-select';
import { t } from '../../i18n';
import { payloadKeys } from '../../mqtt/use-mqtt';
import { useMqttTopic } from '../../mqtt/use-mqtt';
import { getUnsNode } from '../../api/uns';
import { getModelFormat } from '../../utils/thumbnail';
import { flattenNodeTree } from '../../viewer/gltf-nodes';
import { useSplatViewer } from '../../viewer/splat-viewer';
import { useModelViewer } from '../../viewer/use-model-viewer';
import { DataTaggingPanel } from './data-tagging';
import { MappingTable } from './mapping-table';
import InstanceQrModal from './qr-modal';
import { TagOverlay } from './tag-overlay';

// 视口右上角按钮（对齐云端 ViewportIconButton 风格）
const viewportButtonStyle: React.CSSProperties = {
  width: 32,
  height: 32,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  borderRadius: 4,
  cursor: 'pointer',
  border: '1px solid var(--ui-header-splitter-color, #e0e0e0)',
  background: 'var(--ui-bg-color)',
  color: 'var(--ui-text-color)',
};

export default function InstanceDetailPage() {
  const { modelId, instanceId } = useParams();
  const navigate = useNavigate();
  const { message } = App.useApp();
  const [instance, setInstance] = useState<InstanceInfo | null>(null);
  const [binding, setBinding] = useState<InstanceBindingPayload | null>(null);
  const [, setDirty] = useState(false);
  const [activeTab, setActiveTab] = useState('mapping');
  const [topicEditOpen, setTopicEditOpen] = useState(false);
  const [heightEditOpen, setHeightEditOpen] = useState(false);
  const [pendingTopic, setPendingTopic] = useState<{ unsNodeId: string; topic: string }>();
  const [pendingHeight, setPendingHeight] = useState<number | null>(null);
  const [schemaKeys, setSchemaKeys] = useState<string[]>([]);
  const [placementMode, setPlacementMode] = useState(false);
  const [selectedTagId, setSelectedTagId] = useState<string | null>(null);
  // 实时绑定开关（对齐云端视口绿点按钮）：关闭时暂停 MQTT 驱动 3D
  const [liveEnabled, setLiveEnabled] = useState(true);
  const [qrModalOpen, setQrModalOpen] = useState(false);

  // SPLAT 模型（对齐云端）：只支持静态数据标点，无节点树/运动映射，视口走高斯泼溅查看器
  const isSplat = instance ? getModelFormat(instance.modelOriginFile || instance.modelName) === 'splat' : false;
  const { containerRef, viewer } = useModelViewer(instance && !isSplat ? instance.modelFileUrl : undefined);
  const { containerRef: splatContainerRef, viewer: splatViewer } = useSplatViewer(
    instance && isSplat ? instance.modelFileUrl : undefined
  );
  const overlayViewer = isSplat ? splatViewer : viewer;
  useEffect(() => {
    if (isSplat) setActiveTab('tagging');
  }, [isSplat]);
  // 页签隐藏/窗口后台 → 暂停渲染循环与 MQTT 订阅，恢复后自动重连（消息 effect 会立即驱动最新姿态）
  const suspended = useSuspended();
  useEffect(() => {
    viewer.setSuspended(suspended);
    splatViewer.setSuspended(suspended);
  }, [suspended, viewer, splatViewer]);
  const { message: mqttMessage, connected } = useMqttTopic(instance?.topic || undefined, undefined, suspended);
  const lastPayloadRef = useRef<Record<string, unknown>>({});
  const bindingImportRef = useRef<HTMLInputElement>(null);

  const refresh = useCallback(() => {
    if (!instanceId) return;
    getInstance(instanceId)
      .then((data) => {
        setInstance(data);
        setBinding(parseBinding(data.bindingJson));
        setDirty(false);
      })
      .catch((e) => message.error((e as Error).message));
  }, [instanceId, message]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  // UNS 节点 schema 字段（payload key 候选，和实时消息 key 合并）+ 真实 namespace 路径（头部展示用，订阅仍用 alias）
  const [topicPath, setTopicPath] = useState('');
  useEffect(() => {
    if (!instance?.unsNodeId) {
      setTopicPath('');
      return;
    }
    getUnsNode(instance.unsNodeId)
      .then((node) => {
        setSchemaKeys((node.fields || []).map((field) => field.name));
        setTopicPath(String(node.namespace || node.path || ''));
      })
      .catch(() => {
        setSchemaKeys([]);
        setTopicPath('');
      });
  }, [instance?.unsNodeId]);

  const availableKeys = useMemo(() => {
    const keys = new Set(schemaKeys);
    payloadKeys(mqttMessage?.payload).forEach((key) => keys.add(key));
    return Array.from(keys);
  }, [schemaKeys, mqttMessage?.payload]);

  // 绑定变化 → 重建 viewer 绑定；新消息 → 驱动 3D 运动
  useEffect(() => {
    if (viewer.status !== 'ready' || !binding) return;
    viewer.setBindings(binding.motionMappings);
    if (Object.keys(lastPayloadRef.current).length) viewer.applyPayload(lastPayloadRef.current);
  }, [viewer.status, binding, viewer]);

  useEffect(() => {
    if (!mqttMessage) return;
    lastPayloadRef.current = mqttMessage.payload;
    if (liveEnabled) viewer.applyPayload(mqttMessage.payload);
  }, [mqttMessage, viewer, liveEnabled]);

  // 实时绑定关闭 → 模型回到初始静态姿态；重新开启 → 立即应用最近一条数据
  useEffect(() => {
    if (viewer.status !== 'ready') return;
    if (!liveEnabled) viewer.resetPose();
    else if (Object.keys(lastPayloadRef.current).length) viewer.applyPayload(lastPayloadRef.current);
  }, [liveEnabled, viewer]);

  const running = connected && Boolean(mqttMessage);

  // 即时自动保存（对齐云端：修改即入队防抖提交，成功后 toast）
  const persistTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const bindingRef = useRef<InstanceBindingPayload | null>(null);
  bindingRef.current = binding;
  const pendingToastRef = useRef<string | undefined>(undefined);

  const schedulePersist = useCallback(
    (toast?: string) => {
      if (toast) pendingToastRef.current = toast;
      if (persistTimerRef.current) clearTimeout(persistTimerRef.current);
      persistTimerRef.current = setTimeout(async () => {
        const current = bindingRef.current;
        if (!instance || !current) return;
        try {
          await updateInstance(instance.id, { bindingJson: JSON.stringify(current) });
          setDirty(false);
          if (pendingToastRef.current) {
            message.success(pendingToastRef.current);
            pendingToastRef.current = undefined;
          }
        } catch (e) {
          message.error((e as Error).message);
        }
      }, 500);
    },
    [instance, message]
  );

  const patchBinding = useCallback(
    (updater: (prev: InstanceBindingPayload) => InstanceBindingPayload, toast?: string) => {
      setBinding((prev) => (prev ? updater(prev) : prev));
      setDirty(true);
      schedulePersist(toast);
    },
    [schedulePersist]
  );

  const saveTopic = async () => {
    if (!instance || !pendingTopic?.topic) return;
    try {
      await updateInstance(instance.id, { unsNodeId: pendingTopic.unsNodeId, topic: pendingTopic.topic });
      setTopicEditOpen(false);
      refresh();
    } catch (e) {
      message.error((e as Error).message);
    }
  };

  const saveHeight = async () => {
    if (!instance || !pendingHeight || pendingHeight <= 0 || pendingHeight > 1000) return;
    try {
      await updateInstance(instance.id, { height: pendingHeight });
      setHeightEditOpen(false);
      message.success(t('instance.heightOk'));
      refresh();
    } catch (e) {
      message.error((e as Error).message);
    }
  };

  // Edit Selected：右侧面板切换为节点树编辑（对齐云端流程）
  const [nodeEditorOpen, setNodeEditorOpen] = useState(false);
  const [checkedPaths, setCheckedPaths] = useState<string[]>([]);
  // 节点树异步加载，defaultExpandAll 不生效——树数据到位后受控全展开
  const [expandedPaths, setExpandedPaths] = useState<string[]>([]);
  useEffect(() => {
    setExpandedPaths(flattenNodeTree(viewer.nodeTree).map((node) => node.path));
  }, [viewer.nodeTree]);
  const confirmSelectedNodes = () => {
    const flat = flattenNodeTree(viewer.nodeTree);
    const byPath = new Map(flat.map((node) => [node.path, node]));
    const nextSelected = checkedPaths
      .map((path) => byPath.get(path))
      .filter(Boolean)
      .map((node) => ({ nodeID: node!.tier0Id, name: node!.title, path: node!.path }));
    const keptIDs = new Set(nextSelected.map((node) => node.nodeID));
    patchBinding(
      (prev) => ({
        ...prev,
        selectedObjects: nextSelected,
        // 移除的节点连带清理其运动映射
        motionMappings: prev.motionMappings.filter((item) => keptIDs.has(item.nodeID)),
      }),
      t('instance.mapping.saved')
    );
    setNodeEditorOpen(false);
    viewer.highlight([]);
  };

  // 3D 上的运动标签（仅 Data Tagging tab 显示，按可见性设置过滤，key 去重取首个绑定节点）
  const motionLabels = useMemo(() => {
    if (activeTab !== 'tagging' || !binding) return [];
    const visibility = binding.motionTagVisibility ?? { mode: 'showAll' as const, custom: {} };
    if (visibility.mode === 'hideAll') return [];
    const seen = new Set<string>();
    const labels: { key: string; path: string }[] = [];
    for (const mapping of binding.motionMappings) {
      if (!mapping.value || seen.has(mapping.value)) continue;
      seen.add(mapping.value);
      if (visibility.mode === 'custom' && visibility.custom[mapping.value] === false) continue;
      labels.push({ key: mapping.value, path: mapping.path });
    }
    return labels;
  }, [activeTab, binding]);

  const handleViewerClick = (event: React.MouseEvent) => {
    if (!placementMode || activeTab !== 'tagging' || !binding) return;
    const hit = isSplat
      ? splatViewer.raycastAt(event.clientX, event.clientY)
      : viewer.raycastAt(event.clientX, event.clientY);
    if (!hit) return;
    window.dispatchEvent(new CustomEvent('anchor:create-tag', { detail: hit }));
  };

  // 选中标点 → 挂平移 gizmo（对齐云端 marker TransformControls）：拖动实时预览，松手提交保存
  const [tagPreview, setTagPreview] = useState<{ id: string; pos: { x: number; y: number; z: number } } | null>(null);
  useEffect(() => {
    if (viewer.status !== 'ready') return;
    const tag =
      activeTab === 'tagging' && selectedTagId
        ? binding?.dataTags.find((item) => item.id === selectedTagId)
        : undefined;
    if (!tag) {
      viewer.setMarkerGizmo(null);
      return;
    }
    viewer.setMarkerGizmo({
      position: { x: tag.x ?? 0, y: tag.y ?? 0, z: tag.z ?? 0 },
      onPreview: (pos) => setTagPreview({ id: tag.id, pos }),
      onCommit: (pos) => {
        setTagPreview(null);
        patchBinding(
          (prev) => ({
            ...prev,
            dataTags: prev.dataTags.map((item) => (item.id === tag.id ? { ...item, ...pos } : item)),
          }),
          t('instance.tags.saved')
        );
      },
    });
    return () => viewer.setMarkerGizmo(null);
  }, [viewer, viewer.status, activeTab, selectedTagId, binding, patchBinding]);

  // 拖动中的标签用预览坐标渲染 overlay（仅对当前选中标签生效，选中变化后旧预览自动失效）
  const overlayTags = useMemo(() => {
    if (!binding) return [];
    if (!tagPreview || tagPreview.id !== selectedTagId) return binding.dataTags;
    return binding.dataTags.map((item) => (item.id === tagPreview.id ? { ...item, ...tagPreview.pos } : item));
  }, [binding, tagPreview, selectedTagId]);

  return (
    <div className="flex h-full flex-col">
      <div className="page-title-bar">
        <div className="flex min-w-0 flex-1 items-center gap-3">
          <Button
            variant="outlined"
            color="default"
            style={{ paddingLeft: '5.5px', gap: '3px' }}
            onClick={() => navigate(`/model/${modelId}`)}
          >
            <span className="inline-flex items-center gap-2">
              <ChevronLeft size={16} />
              {t('common.back')}
            </span>
          </Button>
          <div className="truncate">{instance?.name ?? t('common.loading')}</div>
        </div>
        <Badge
          status={running ? 'success' : connected ? 'processing' : 'default'}
          text={
            <span className="text-xs font-normal" style={{ color: 'var(--ui-description-card-color)' }}>
              {running
                ? t('instance.status.running')
                : connected
                  ? t('instance.mqtt.waiting')
                  : t('instance.mqtt.disconnected')}
            </span>
          }
        />
        <input
          ref={bindingImportRef}
          type="file"
          accept="application/json,.json"
          className="hidden"
          onChange={async (e) => {
            const file = e.target.files?.[0];
            if (!file) return;
            try {
              const parsed = JSON.parse(await file.text());
              if (!parsed || typeof parsed !== 'object' || !Array.isArray(parsed.motionMappings)) {
                throw new Error('invalid');
              }
              patchBinding(() => ({
                ...parsed,
                selectedObjects: parsed.selectedObjects ?? [],
                motionMappings: parsed.motionMappings ?? [],
                dataTags: parsed.dataTags ?? [],
              }));
              message.success(t('instance.binding.imported'));
            } catch {
              message.error(t('instance.binding.invalid'));
            } finally {
              if (bindingImportRef.current) bindingImportRef.current.value = '';
            }
          }}
        />
        <div className="flex items-center gap-2">
          <Button type="primary" icon={<QrCode size={16} />} onClick={() => setQrModalOpen(true)}>
            {t('instance.qr')}
          </Button>
          <Dropdown
            placement="bottomRight"
            menu={{
              items: [
                {
                  key: 'export',
                  icon: <Download size={16} />,
                  label: t('instance.binding.export'),
                  disabled: !binding,
                },
                {
                  key: 'import',
                  icon: <Upload size={16} />,
                  label: t('instance.binding.import'),
                },
              ],
              onClick: ({ key }) => {
                if (key === 'import') {
                  bindingImportRef.current?.click();
                  return;
                }
                if (key === 'export') {
                  if (!binding || !instance) return;
                  const blob = new Blob([JSON.stringify(binding, null, 2)], { type: 'application/json' });
                  const url = URL.createObjectURL(blob);
                  const link = document.createElement('a');
                  link.href = url;
                  link.download = `${instance.name}-binding.json`;
                  link.click();
                  URL.revokeObjectURL(url);
                }
              },
            }}
          >
            <Button type="text" className="page-title-more" icon={<EllipsisVertical size={16} />} />
          </Dropdown>
        </div>
      </div>

      <div className="min-h-0 flex-1">
        <ResizableSplit
          defaultLeftWidth={38}
          minLeftWidth={24}
          maxLeftWidth={48}
          left={
            /* 左：viewer（含标签 overlay），宽度可拖拽（对齐云端 38%/24-48%） */
            <div
              className="relative h-full w-full"
              style={{ cursor: placementMode ? 'crosshair' : undefined }}
              onClick={handleViewerClick}
            >
              <div ref={containerRef} className="h-full w-full" style={isSplat ? { display: 'none' } : undefined} />
              {isSplat ? <div ref={splatContainerRef} className="h-full w-full" /> : null}
              {binding ? (
                <TagOverlay
                  tags={overlayTags}
                  motionLabels={motionLabels}
                  viewer={overlayViewer}
                  suspended={suspended}
                  payload={mqttMessage?.payload}
                  messageTs={mqttMessage?.ts}
                  selectedTagId={selectedTagId}
                  onSelectTag={(id) => {
                    setSelectedTagId(id);
                    setPlacementMode(false);
                    setActiveTab('tagging');
                  }}
                />
              ) : null}
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
              {/* 视口右上按钮（对齐云端）：实时绑定开关（按连接状态着色）/ 重置视角 / 全屏 */}
              <div className="absolute top-3 right-3 z-10 flex gap-1">
                <Tooltip title={t('viewer.live')}>
                  <button
                    type="button"
                    style={viewportButtonStyle}
                    onClick={() => setLiveEnabled((prev) => !prev)}
                    aria-pressed={liveEnabled}
                  >
                    <span
                      className="h-3.5 w-3.5 rounded-full border"
                      style={{
                        borderColor: 'var(--ui-description-card-color)',
                        background: !liveEnabled
                          ? 'transparent'
                          : running
                            ? 'var(--ui-status-active-text, #53b483)'
                            : connected
                              ? 'var(--ui-t-yellow-color-40, #f1c21b)'
                              : 'transparent',
                      }}
                    />
                  </button>
                </Tooltip>
                <Tooltip title={t('viewer.resetView')}>
                  <button type="button" style={viewportButtonStyle} onClick={() => viewer.resetView()}>
                    <CenterCircle size={16} />
                  </button>
                </Tooltip>
                <Tooltip title={t('viewer.fullscreen')}>
                  <button
                    type="button"
                    style={viewportButtonStyle}
                    onClick={() => {
                      const el = containerRef.current?.parentElement;
                      if (document.fullscreenElement) void document.exitFullscreen();
                      else void el?.requestFullscreen();
                    }}
                  >
                    <Maximize size={16} />
                  </button>
                </Tooltip>
              </div>
            </div>
          }
          right={
            /* 右：topic/height + tabs */
            <div className="flex h-full min-w-0 flex-col overflow-hidden">
              {/* 顶部元信息条：Topic + 源文件 + 高度（与 Model 详情同构） */}
              <div
                className="flex items-center gap-4 border-b px-5 py-4"
                style={{ borderColor: 'var(--ui-line-color, #e0e0e0)' }}
              >
                <div className="min-w-0 flex-[1.4]">
                  <div className="text-xs" style={{ color: 'var(--ui-description-card-color)' }}>
                    {t('instance.topic')}
                  </div>
                  <div className="mt-0.5 flex items-center gap-1 text-sm font-medium">
                    <LinkIcon size={14} style={{ color: 'var(--ui-status-active-text, #53b483)', flexShrink: 0 }} />
                    <span className="truncate" title={topicPath || instance?.topic || undefined}>
                      {topicPath || instance?.topic || '--'}
                    </span>
                    <Button
                      type="text"
                      size="small"
                      icon={<Edit size={14} />}
                      onClick={() => {
                        setPendingTopic(
                          instance?.unsNodeId ? { unsNodeId: instance.unsNodeId, topic: instance.topic } : undefined
                        );
                        setTopicEditOpen(true);
                      }}
                    />
                  </div>
                </div>
                <div className="min-w-0 flex-1">
                  <div className="text-xs" style={{ color: 'var(--ui-description-card-color)' }}>
                    {t('model.detail.originFile')}
                  </div>
                  <div
                    className="mt-0.5 truncate text-sm font-medium"
                    title={instance?.modelOriginFile || undefined}
                  >
                    {instance?.modelOriginFile || '--'}
                  </div>
                </div>
                <div className="min-w-0 flex-1">
                  <div className="text-xs" style={{ color: 'var(--ui-description-card-color)' }}>
                    {t('model.height')}
                  </div>
                  <div className="mt-0.5 flex items-center gap-1 text-sm font-medium">
                    <span className="flex items-baseline gap-1">
                      <span>
                        {(instance?.height || viewer.modelHeight) > 0
                          ? (instance?.height || viewer.modelHeight).toFixed(3)
                          : '--'}
                      </span>
                      <span style={{ color: 'var(--ui-description-card-color)' }}>M</span>
                    </span>
                    <Button
                      type="text"
                      size="small"
                      icon={<Edit size={14} />}
                      onClick={() => {
                        setPendingHeight(instance?.height || Number(viewer.modelHeight.toFixed(3)) || null);
                        setHeightEditOpen(true);
                      }}
                    />
                  </div>
                </div>
              </div>

              {nodeEditorOpen ? (
                <div className="flex min-h-0 flex-1 flex-col px-5 py-3">
                  <div className="mb-2 text-base font-semibold">{t('model.nodes.editTitle')}</div>
                  <div
                    className="min-h-0 flex-1 overflow-auto rounded border p-2"
                    style={{ borderColor: 'var(--ui-header-splitter-color, #e0e0e0)' }}
                  >
                    <Tree
                      checkable
                      checkStrictly
                      expandedKeys={expandedPaths}
                      onExpand={(keys) => setExpandedPaths(keys.map(String))}
                      checkedKeys={{ checked: checkedPaths, halfChecked: [] }}
                      treeData={viewer.nodeTree}
                      fieldNames={{ key: 'path', title: 'title', children: 'children' }}
                      onCheck={(checked) => {
                        const keys = Array.isArray(checked) ? checked : checked.checked;
                        setCheckedPaths(keys.map(String));
                      }}
                      titleRender={(node) => (
                        <span
                          onMouseEnter={() => viewer.highlight([(node as { path: string }).path])}
                          onMouseLeave={() => viewer.highlight(checkedPaths)}
                        >
                          {(node as { title: string }).title}
                        </span>
                      )}
                    />
                  </div>
                  <div className="mt-3 flex shrink-0 justify-end gap-2">
                    <Button
                      onClick={() => {
                        setNodeEditorOpen(false);
                        viewer.highlight([]);
                      }}
                    >
                      {t('common.cancel')}
                    </Button>
                    <Button type="primary" onClick={confirmSelectedNodes}>
                      {t('common.confirm')}
                    </Button>
                  </div>
                </div>
              ) : (
                <div className="flex min-h-0 flex-1 flex-col">
                  {/* 双模式用 Segmented：整行浅底铺满，避免 Tabs 左贴右空 */}
                  {!isSplat ? (
                    <div
                      className="flex shrink-0 items-center border-b px-5 py-2.5"
                      style={{
                        borderColor: 'var(--ui-line-color, #e0e0e0)',
                        background: 'var(--ui-charttop-bg-color)',
                      }}
                    >
                      <Segmented
                        value={activeTab}
                        onChange={(value) => setActiveTab(String(value))}
                        options={[
                          { label: t('instance.tab.mapping'), value: 'mapping' },
                          { label: t('instance.tab.tagging'), value: 'tagging' },
                        ]}
                      />
                    </div>
                  ) : null}
                  <div className="min-h-0 flex-1 overflow-auto px-5 pt-3">
                    {activeTab === 'mapping' && !isSplat && binding ? (
                      <MappingTable
                        binding={binding}
                        availableKeys={availableKeys}
                        payload={mqttMessage?.payload}
                        onChange={patchBinding}
                        onHover={(path) => viewer.highlight(path ? [path] : [])}
                        onEditSelected={() => {
                          setCheckedPaths(binding.selectedObjects.map((node) => node.path));
                          setNodeEditorOpen(true);
                        }}
                      />
                    ) : null}
                    {activeTab === 'tagging' && binding ? (
                      <DataTaggingPanel
                        binding={binding}
                        availableKeys={availableKeys}
                        placementMode={placementMode}
                        onPlacementModeChange={setPlacementMode}
                        onChange={patchBinding}
                        selectedTagId={selectedTagId}
                        onSelectTag={setSelectedTagId}
                      />
                    ) : null}
                  </div>
                </div>
              )}
            </div>
          }
        />
      </div>

      <Modal title={t('instance.topic')} open={topicEditOpen} onCancel={() => setTopicEditOpen(false)} onOk={saveTopic}>
        <TopicSelect value={pendingTopic} onChange={setPendingTopic} />
      </Modal>

      <Modal
        title={t('instance.height')}
        open={heightEditOpen}
        onCancel={() => setHeightEditOpen(false)}
        onOk={saveHeight}
      >
        <InputNumber
          min={0.001}
          max={1000}
          precision={3}
          style={{ width: '100%' }}
          value={pendingHeight}
          onChange={setPendingHeight}
        />
        <div className="mt-1 text-xs" style={{ color: 'var(--ui-description-card-color)' }}>
          {t('instance.heightHint')}
        </div>
      </Modal>

      <InstanceQrModal
        open={qrModalOpen}
        onClose={() => setQrModalOpen(false)}
        instance={instance ? { id: instance.id, name: instance.name } : null}
      />
    </div>
  );
}
