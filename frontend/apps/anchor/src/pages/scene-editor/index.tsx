// 场景编辑器 — 布局对齐云端原版 scene-editor-shell：
//   画布全屏铺底；实例库左上悬浮卡片（可折叠+搜索）；变换工具条顶部居中悬浮（模式 + Local/World）；
//   右上视口按钮（重置视角/全屏）；右侧 295px 固定栏 = Scene Collection + 可拖拽分割 + Settings 三 tab。
import {
  CenterCircle,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  Link as LinkIcon,
  Maximize,
  Move,
  RotateClockwise,
  Save,
  Scale,
  TrashCan,
  Unlink,
} from '@carbon/icons-react';
import { App, Badge, Button, Empty, Input, InputNumber, Select, Slider, Spin, Switch, Tooltip } from 'antd';
import { useEffect, useMemo, useRef, useState } from 'react';
import { useNavigate, useParams } from 'react-router';
import { listInstances, parseBinding, type InstanceInfo } from '../../api/instances';
import { uploadAsset } from '../../api/models';
import {
  defaultPlacement,
  defaultSceneConfig,
  getScene,
  parsePlacement,
  parseSceneConfig,
  updateScene,
  type SceneConfigV4,
  type SceneInfo,
} from '../../api/scenes';
import { useSuspended } from '../../bridge/use-suspended';
import { t } from '../../i18n';
import { useMqttTopics } from '../../mqtt/use-mqtt';
import { LIGHT_PRESETS, useSceneEditor, type GizmoMode } from '../../viewer/use-scene-editor';

interface EditorItemMeta {
  key: string;
  instanceId: number;
  instanceName: string;
  topic: string;
}

let itemSeq = 0;
const nextKey = (instanceId: number) => `item_${instanceId}_${++itemSeq}`;

const panelStyle: React.CSSProperties = {
  background: 'var(--ui-bg-color)',
  border: '1px solid var(--ui-header-splitter-color, #e0e0e0)',
  borderRadius: 4,
};

const mutedColor = { color: 'var(--ui-description-card-color)' } as const;

export default function SceneEditorPage() {
  const { sceneId } = useParams();
  const navigate = useNavigate();
  const { message } = App.useApp();
  const editor = useSceneEditor();
  // 页签隐藏/窗口后台 → 暂停渲染循环与 MQTT 订阅
  const suspended = useSuspended();
  useEffect(() => {
    editor.setSuspended(suspended);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [suspended, editor.setSuspended]);
  const [scene, setScene] = useState<SceneInfo | null>(null);
  const [config, setConfig] = useState<SceneConfigV4 | null>(null);
  const [items, setItems] = useState<EditorItemMeta[]>([]);
  const [instances, setInstances] = useState<InstanceInfo[]>([]);
  const [mode, setModeState] = useState<GizmoMode>('translate');
  const [space, setSpaceState] = useState<'local' | 'world'>('world');
  const [saving, setSaving] = useState(false);
  const [dirty, setDirty] = useState(false);
  const [libraryCollapsed, setLibraryCollapsed] = useState(false);
  const [librarySearch, setLibrarySearch] = useState('');
  const [activeTab, setActiveTab] = useState<'viewport' | 'light' | 'object'>('viewport');
  const [collectionRatio, setCollectionRatio] = useState(0.34);
  const loadedRef = useRef(false);
  const editorAreaRef = useRef<HTMLDivElement>(null);
  const rightPanelRef = useRef<HTMLDivElement>(null);

  // 加载场景 + 实例库
  useEffect(() => {
    if (!sceneId) return;
    getScene(sceneId)
      .then((data) => {
        setScene(data);
        setConfig(parseSceneConfig(data.configJson));
      })
      .catch((e) => message.error((e as Error).message));
    listInstances({ page: 1, size: 200 })
      .then((data) => setInstances(data.list || []))
      .catch(() => setInstances([]));
  }, [sceneId, message]);

  // 编辑器就绪后恢复已保存的 items 与配置
  useEffect(() => {
    if (!editor.ready || !scene || !config || loadedRef.current) return;
    loadedRef.current = true;
    editor.applyConfig(config);
    (async () => {
      for (const item of scene.items) {
        if (!item.instance?.modelFileUrl) continue;
        const key = nextKey(item.instanceId);
        await editor.addItem({
          key,
          instanceId: item.instanceId,
          fileUrl: item.instance.modelFileUrl,
          instanceHeight: item.instance.height,
          placement: parsePlacement(item.placementJson),
          motionMappings: parseBinding(item.instance.bindingJson).motionMappings,
        });
        setItems((prev) => [
          ...prev,
          { key, instanceId: item.instanceId, instanceName: item.instance!.name, topic: item.instance!.topic || '' },
        ]);
      }
    })();
  }, [editor.ready, scene, config, editor]);

  useEffect(() => {
    if (editor.ready && config && loadedRef.current) editor.applyConfig(config);
  }, [config, editor]);

  // 场景实时运行时：订阅所有实例 topic，消息驱动绑定节点
  const itemsRef = useRef(items);
  itemsRef.current = items;
  const liveTopics = useMemo(() => items.map((item) => item.topic).filter(Boolean), [items]);
  const { connected: liveConnected } = useMqttTopics(
    liveTopics,
    (topic, payload) => {
      for (const item of itemsRef.current) {
        if (item.topic === topic) editor.applyPayloadToItem(item.key, payload);
      }
    },
    undefined,
    suspended
  );

  // 选中对象时自动切到 Object tab（对齐源端行为）
  useEffect(() => {
    if (editor.selectedKey) setActiveTab('object');
    else setActiveTab((tab) => (tab === 'object' ? 'viewport' : tab));
  }, [editor.selectedKey]);

  // 快捷键 W/E/R + Delete + Esc
  useEffect(() => {
    const onKey = (event: KeyboardEvent) => {
      if ((event.target as HTMLElement)?.tagName === 'INPUT') return;
      if (event.key === 'w' || event.key === 'W') applyMode('translate');
      if (event.key === 'e' || event.key === 'E') applyMode('rotate');
      if (event.key === 'r' || event.key === 'R') applyMode('scale');
      if ((event.key === 'Delete' || event.key === 'Backspace') && editor.selectedKey) removeItem(editor.selectedKey);
      if (event.key === 'Escape') editor.select(null);
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [editor.selectedKey, editor]);

  const applyMode = (next: GizmoMode) => {
    setModeState(next);
    editor.setMode(next);
  };

  const applySpace = (next: 'local' | 'world') => {
    setSpaceState(next);
    editor.setSpace(next);
  };

  const addInstance = async (instance: InstanceInfo) => {
    if (!instance.modelFileUrl) return;
    // 对齐云端：同一实例在一个场景中只能添加一次，重复添加弹警告
    if (items.some((item) => item.instanceId === instance.id)) {
      message.error(t('scene.editor.dupInstance'));
      return;
    }
    const key = nextKey(instance.id);
    const placement = defaultPlacement();
    placement.position = [items.length * 1.5, 0, 0];
    await editor.addItem({
      key,
      instanceId: instance.id,
      fileUrl: instance.modelFileUrl,
      instanceHeight: instance.height,
      placement,
      motionMappings: parseBinding(instance.bindingJson).motionMappings,
    });
    setItems((prev) => [
      ...prev,
      { key, instanceId: instance.id, instanceName: instance.name, topic: instance.topic || '' },
    ]);
    editor.select(key);
    setDirty(true);
  };

  const removeItem = (key: string) => {
    editor.removeItem(key);
    setItems((prev) => prev.filter((item) => item.key !== key));
    setDirty(true);
  };

  useEffect(() => {
    if (editor.transformTick > 0 && loadedRef.current) setDirty(true);
  }, [editor.transformTick]);

  const selectedPlacement = useMemo(
    () => (editor.selectedKey ? editor.getPlacement(editor.selectedKey) : null),
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [editor.selectedKey, editor.transformTick, editor]
  );

  const selectedPivot = useMemo(
    () => (editor.selectedKey ? editor.getPivot(editor.selectedKey) : 'bottomCenter'),
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [editor.selectedKey, editor.transformTick, editor]
  );

  const patchPlacement = (axisIndex: number, field: 'position' | 'rotation' | 'scale', value: number | null) => {
    if (!editor.selectedKey || value === null) return;
    const placement = editor.getPlacement(editor.selectedKey);
    if (!placement) return;
    const next = { ...placement, [field]: [...placement[field]] as [number, number, number] };
    next[field][axisIndex] = field === 'rotation' ? (value * Math.PI) / 180 : value;
    editor.setPlacement(editor.selectedKey, next);
    setDirty(true);
  };

  const handleSave = async () => {
    if (!scene || !config) return;
    setSaving(true);
    try {
      const payloadItems = items
        .map((item, index) => {
          const placement = editor.getPlacement(item.key);
          return placement
            ? { instanceId: item.instanceId, placementJson: JSON.stringify(placement), sort: index }
            : null;
        })
        .filter(Boolean) as { instanceId: number; placementJson: string; sort: number }[];

      let thumbnailAssetId: number | undefined;
      editor.select(null);
      await new Promise((resolve) => requestAnimationFrame(() => requestAnimationFrame(resolve)));
      const shot = await editor.snapshot();
      if (shot) {
        try {
          thumbnailAssetId = (await uploadAsset(shot)).fileId;
        } catch {
          /* 缩略图失败不阻塞保存 */
        }
      }
      await updateScene(scene.id, {
        configJson: JSON.stringify(config),
        items: payloadItems,
        itemsSet: true,
        thumbnailAssetId,
      });
      setDirty(false);
      message.success(t('scene.editor.saved'));
    } catch (e) {
      message.error((e as Error).message);
    } finally {
      setSaving(false);
    }
  };

  // 右栏 Collection/Settings 分割拖拽
  const handleResizeStart = (event: React.PointerEvent) => {
    event.preventDefault();
    const panel = rightPanelRef.current;
    if (!panel) return;
    const onMove = (moveEvent: PointerEvent) => {
      const rect = panel.getBoundingClientRect();
      const ratio = (moveEvent.clientY - rect.top) / rect.height;
      setCollectionRatio(Math.min(0.7, Math.max(0.15, ratio)));
    };
    const onUp = () => {
      window.removeEventListener('pointermove', onMove);
      window.removeEventListener('pointerup', onUp);
    };
    window.addEventListener('pointermove', onMove);
    window.addEventListener('pointerup', onUp);
  };

  const filteredInstances = useMemo(
    () =>
      instances.filter((instance) =>
        librarySearch ? instance.name.toLowerCase().includes(librarySearch.toLowerCase()) : true
      ),
    [instances, librarySearch]
  );

  const toolButton = (active: boolean): React.CSSProperties => ({
    width: 32,
    height: 32,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    borderRadius: 4,
    cursor: 'pointer',
    border: active ? '1px solid var(--ui-theme-color)' : '1px solid transparent',
    background: active ? 'var(--ui-primary-bg)' : 'transparent',
    color: 'var(--ui-text-color)',
  });

  const spaceButton = (active: boolean): React.CSSProperties => ({
    padding: '4px 12px',
    fontSize: 11,
    fontWeight: 600,
    letterSpacing: '0.05em',
    textTransform: 'uppercase',
    borderRadius: 4,
    cursor: 'pointer',
    border: active ? '1px solid var(--ui-theme-color)' : '1px solid transparent',
    background: active ? 'var(--ui-primary-bg)' : 'transparent',
    color: 'var(--ui-text-color)',
  });

  const fieldRow = (label: string, children: React.ReactNode) => (
    <div className="mb-2 flex items-center gap-2">
      <span className="w-[64px] shrink-0 text-right text-xs leading-6 font-medium">{label}</span>
      {children}
    </div>
  );

  // Viewport/Light tab 分段折叠（对齐云端 settings-collapsible-section：标题 + Reset + 折叠）
  const [openSections, setOpenSections] = useState({ viewport: true, camera: true, grid: true, light: true });
  const patchViewport = (patch: Partial<SceneConfigV4['viewport']>) => {
    if (!config) return;
    setConfig({ ...config, viewport: { ...config.viewport, ...patch } });
    setDirty(true);
  };
  const settingsSection = (
    key: 'viewport' | 'camera' | 'grid' | 'light',
    title: string,
    onReset: () => void,
    children: React.ReactNode
  ) => (
    <div className="border-b pb-4" style={{ borderColor: 'var(--ui-header-splitter-color, #e0e0e0)' }}>
      <div className="flex items-center justify-between py-2">
        <button
          type="button"
          className="flex cursor-pointer items-center gap-1 border-0 bg-transparent p-0 text-xs font-semibold"
          style={{ color: 'var(--ui-text-color)' }}
          onClick={() => setOpenSections((prev) => ({ ...prev, [key]: !prev[key] }))}
        >
          {openSections[key] ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
          {title}
        </button>
        <Button type="link" size="small" style={{ color: 'var(--ui-description-card-color)' }} onClick={onReset}>
          {t('scene.editor.sectionReset')}
        </Button>
      </div>
      {openSections[key] ? <div className="flex flex-col gap-2 pl-2">{children}</div> : null}
    </div>
  );
  const viewportNumField = (
    label: string,
    value: number,
    onCommit: (next: number) => void,
    suffix?: string,
    disabled = false
  ) =>
    fieldRow(
      label,
      <InputNumber
        size="small"
        className="flex-1"
        value={value}
        disabled={disabled}
        suffix={
          suffix ? <span style={{ color: 'var(--ui-description-card-color)', fontSize: 11 }}>{suffix}</span> : undefined
        }
        onChange={(next) => {
          if (next === null || !Number.isFinite(next)) return;
          onCommit(next);
        }}
      />
    );
  // Light tab 分组 + 强度控件（Slider + 数字输入联动，limits 对齐云端 render-effects）
  const lightGroup = (title: string, children: React.ReactNode) => (
    <div className="border-t pt-2" style={{ borderColor: 'var(--ui-header-splitter-color, #e0e0e0)' }}>
      <div className="mb-2 text-xs font-medium">{title}</div>
      {children}
    </div>
  );
  const lightIntensityField = (
    key:
      'environmentIntensity' | 'ambientIntensity' | 'keyLightIntensity' | 'fillLightIntensity' | 'backLightIntensity',
    max: number
  ) => {
    if (!config) return null;
    const commit = (next: number | null) => {
      if (next === null || !Number.isFinite(next)) return;
      setConfig({ ...config, light: { ...config.light, [key]: Math.min(Math.max(next, 0), max) } });
      setDirty(true);
    };
    return fieldRow(
      t('scene.editor.intensity'),
      <div className="flex min-w-0 flex-1 items-center gap-2">
        <Slider className="min-w-0 flex-1" min={0} max={max} step={0.05} value={config.light[key]} onChange={commit} />
        <InputNumber
          size="small"
          style={{ width: 64 }}
          min={0}
          max={max}
          step={0.05}
          value={config.light[key]}
          onChange={commit}
        />
      </div>
    );
  };

  // 背景预设色卡（预览色值与源端 viewport-panel 一致，属内容色板而非主题色）
  const bgSwatches: Array<{ key: SceneConfigV4['viewport']['backgroundPreset']; preview: string }> = [
    { key: 'white', preview: '#fcfcfc' },
    { key: 'lightGray', preview: '#b2c1d7' },
    { key: 'black', preview: '#1a1a1a' },
    { key: 'gradient', preview: 'linear-gradient(to top, #1a1a1a 0%, #404040 50%, #808080 100%)' },
  ];

  const axisGroup = (field: 'position' | 'rotation' | 'scale') => {
    if (!selectedPlacement) return null;
    return (
      <div className="mb-3">
        <div className="mb-1 text-xs font-semibold">{t(`scene.editor.${field}`)}</div>
        {(['X', 'Y', 'Z'] as const).map((axis, axisIndex) =>
          fieldRow(
            axis,
            <InputNumber
              size="small"
              className="flex-1"
              style={{ width: '100%' }}
              step={field === 'scale' ? 0.1 : field === 'rotation' ? 5 : 0.1}
              value={Number(
                (field === 'rotation'
                  ? (selectedPlacement[field][axisIndex] * 180) / Math.PI
                  : selectedPlacement[field][axisIndex]
                ).toFixed(2)
              )}
              onChange={(value) => patchPlacement(axisIndex, field, value)}
            />
          )
        )}
      </div>
    );
  };

  return (
    <div className="flex h-full flex-col">
      {/* 顶部 header：返回 + 场景名（左），实时状态 + 保存（右）——对齐 Flow 详情页 */}
      <div className="page-title-bar">
        <div className="flex min-w-0 flex-1 items-center gap-3">
          <Button
            variant="outlined"
            color="default"
            style={{ paddingLeft: '5.5px', gap: '3px' }}
            onClick={() => navigate('/scene')}
          >
            <span className="inline-flex items-center gap-2">
              <ChevronLeft size={16} />
              {t('common.back')}
            </span>
          </Button>
          <div className="truncate">{scene?.name ?? t('common.loading')}</div>
        </div>
        {liveTopics.length > 0 ? (
          <Badge
            status={liveConnected ? 'success' : 'default'}
            text={
              <span className="text-xs font-normal" style={mutedColor}>
                {liveConnected ? t('instance.status.running') : t('instance.mqtt.disconnected')}
              </span>
            }
          />
        ) : null}
        <Button type="primary" icon={<Save size={16} />} loading={saving} disabled={!dirty} onClick={handleSave}>
          {t('scene.editor.save')}
        </Button>
      </div>

      <div className="flex min-h-0 flex-1">
        {/* 画布区：全屏铺底 + 悬浮面板 */}
        <div ref={editorAreaRef} className="relative min-w-0 flex-1 overflow-hidden">
          <div ref={editor.containerRef} className="absolute inset-0" />
          {!editor.ready ? (
            <div
              className="absolute inset-0 flex items-center justify-center"
              style={{ background: 'var(--ui-bg-color)' }}
            >
              <Spin />
            </div>
          ) : null}

          {/* 顶部居中：变换工具条（模式组 + Local/Global 组 + Center/Bottom pivot 组，对齐云端 scene-transform-toolbar） */}
          <div className="pointer-events-none absolute top-4 left-1/2 z-20 flex -translate-x-1/2 items-center gap-2">
            <div className="pointer-events-auto flex items-center gap-0.5 p-1" style={panelStyle}>
              {(
                [
                  ['translate', <Move key="t" size={16} />, t('scene.editor.mode.translate')],
                  ['rotate', <RotateClockwise key="r" size={16} />, t('scene.editor.mode.rotate')],
                  ['scale', <Scale key="s" size={16} />, t('scene.editor.mode.scale')],
                ] as const
              ).map(([value, icon, title]) => (
                <Tooltip key={value} title={title}>
                  <div style={toolButton(mode === value)} onClick={() => applyMode(value)}>
                    {icon}
                  </div>
                </Tooltip>
              ))}
            </div>
            <div className="pointer-events-auto flex items-center gap-0.5 p-1" style={panelStyle}>
              <div style={spaceButton(space === 'local')} onClick={() => applySpace('local')}>
                {t('scene.editor.space.local')}
              </div>
              <div style={spaceButton(space === 'world')} onClick={() => applySpace('world')}>
                {t('scene.editor.space.world')}
              </div>
            </div>
            <div className="pointer-events-auto flex items-center gap-0.5 p-1" style={panelStyle}>
              {(
                [
                  ['center', t('scene.editor.pivot.center')],
                  ['bottomCenter', t('scene.editor.pivot.bottom')],
                ] as const
              ).map(([value, label]) => (
                <div
                  key={value}
                  style={{
                    ...spaceButton(selectedPivot === value),
                    ...(editor.selectedKey ? null : { opacity: 0.45, cursor: 'not-allowed' }),
                  }}
                  onClick={() => {
                    if (!editor.selectedKey) return;
                    editor.setPivot(editor.selectedKey, value);
                    setDirty(true);
                  }}
                >
                  {label}
                </div>
              ))}
            </div>
          </div>

          {/* 右上：视口按钮（重置视角 / 全屏） */}
          <div className="absolute top-4 right-4 z-20 flex items-center gap-0.5 p-1" style={panelStyle}>
            <Tooltip title={t('scene.editor.resetView')}>
              <div style={toolButton(false)} onClick={() => editor.frameAll()}>
                <CenterCircle size={16} />
              </div>
            </Tooltip>
            <Tooltip title={t('scene.editor.fullscreen')}>
              <div
                style={toolButton(false)}
                onClick={() => {
                  const area = editorAreaRef.current;
                  if (!area) return;
                  if (document.fullscreenElement) void document.exitFullscreen();
                  else void area.requestFullscreen();
                }}
              >
                <Maximize size={16} />
              </div>
            </Tooltip>
          </div>

          {/* 左上：实例库悬浮卡片（可折叠 + 搜索） */}
          <div
            className="absolute top-4 left-4 z-20 flex flex-col overflow-hidden transition-all"
            style={{ ...panelStyle, width: libraryCollapsed ? 168 : 260, height: libraryCollapsed ? 48 : 420 }}
          >
            <div
              className="flex h-12 shrink-0 cursor-pointer items-center gap-2 border-b px-3"
              style={{ borderColor: 'var(--ui-header-splitter-color, #e0e0e0)' }}
              onClick={() => setLibraryCollapsed((value) => !value)}
            >
              {libraryCollapsed ? <ChevronRight size={16} /> : <ChevronDown size={16} />}
              <span className="text-sm font-semibold">{t('scene.editor.library')}</span>
            </div>
            {!libraryCollapsed ? (
              <>
                <div className="border-b px-3 py-2" style={{ borderColor: 'var(--ui-header-splitter-color, #e0e0e0)' }}>
                  <Input
                    size="small"
                    allowClear
                    placeholder={t('scene.editor.searchInstance')}
                    value={librarySearch}
                    onChange={(e) => setLibrarySearch(e.target.value)}
                  />
                </div>
                <div className="min-h-0 flex-1 overflow-auto p-2">
                  {filteredInstances.length === 0 ? (
                    <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} />
                  ) : (
                    filteredInstances.map((instance) => (
                      <div
                        key={instance.id}
                        className="mb-2 flex cursor-pointer items-center justify-between gap-2 rounded-[6px] border px-3 py-2 hover:shadow"
                        style={{ borderColor: 'var(--ui-header-splitter-color, #e0e0e0)' }}
                        onClick={() => void addInstance(instance)}
                      >
                        <div className="min-w-0">
                          <div className="truncate text-sm">{instance.name}</div>
                          <div className="mt-0.5 flex items-center gap-1 text-[10px]" style={mutedColor}>
                            {instance.topic ? (
                              <>
                                <LinkIcon size={10} style={{ color: 'var(--ui-status-active-text, #53b483)' }} />
                                {t('scene.editor.connected')}
                              </>
                            ) : (
                              <>
                                <Unlink size={10} />
                                {t('scene.editor.disconnected')}
                              </>
                            )}
                          </div>
                        </div>
                        <div
                          className="flex h-7 w-7 shrink-0 items-center justify-center rounded-[4px] border text-base"
                          style={{ borderColor: 'var(--ui-header-splitter-color, #e0e0e0)' }}
                        >
                          +
                        </div>
                      </div>
                    ))
                  )}
                </div>
              </>
            ) : null}
          </div>
        </div>

        {/* 右侧 295px 固定栏：Scene Collection + 分割条 + Settings */}
        <div
          ref={rightPanelRef}
          className="flex h-full w-[295px] shrink-0 flex-col border-l"
          style={{ borderColor: 'var(--ui-header-splitter-color, #e0e0e0)', background: 'var(--ui-bg-color)' }}
        >
          {/* Collection 区 */}
          <div className="flex min-h-[120px] flex-col" style={{ flexBasis: `${collectionRatio * 100}%` }}>
            <div
              className="flex h-10 shrink-0 items-center border-b px-4 text-sm font-semibold"
              style={{ borderColor: 'var(--ui-header-splitter-color, #e0e0e0)' }}
            >
              {t('scene.editor.collection')}
            </div>
            <div className="min-h-0 flex-1 overflow-auto px-3 py-2">
              {items.length === 0 ? (
                <div className="pt-4 text-center text-xs" style={mutedColor}>
                  {t('scene.editor.emptyCollection')}
                </div>
              ) : (
                items.map((item) => (
                  <div
                    key={item.key}
                    className="group flex cursor-pointer items-center gap-1 rounded-[2px] px-1.5 py-1 text-sm"
                    style={
                      editor.selectedKey === item.key
                        ? { background: 'var(--ui-primary-bg)', color: 'var(--ui-text-color)' }
                        : undefined
                    }
                    onClick={() => editor.select(item.key)}
                  >
                    <span className="min-w-0 flex-1 truncate">{item.instanceName}</span>
                    <Button
                      type="text"
                      size="small"
                      danger
                      className="opacity-0 group-hover:opacity-100"
                      icon={<TrashCan size={14} />}
                      onClick={(e) => {
                        e.stopPropagation();
                        removeItem(item.key);
                      }}
                    />
                  </div>
                ))
              )}
            </div>
          </div>

          {/* 分割拖拽条 */}
          <div className="relative h-2 shrink-0 cursor-row-resize" onPointerDown={handleResizeStart}>
            <div
              className="absolute inset-x-0 top-1/2 h-px -translate-y-1/2"
              style={{ background: 'var(--ui-header-splitter-color, #e0e0e0)' }}
            />
          </div>

          {/* Settings 区 */}
          <div className="flex min-h-[160px] flex-1 flex-col">
            <div
              className="flex h-10 shrink-0 items-center border-y px-4 text-sm font-semibold"
              style={{ borderColor: 'var(--ui-header-splitter-color, #e0e0e0)' }}
            >
              {t('scene.editor.settings')}
            </div>
            <div
              className="flex h-9 shrink-0 items-center gap-5 border-b px-4"
              style={{ borderColor: 'var(--ui-header-splitter-color, #e0e0e0)' }}
            >
              {(
                [
                  ['viewport', t('scene.editor.tab.viewport'), false],
                  ['light', t('scene.editor.tab.light'), false],
                  ['object', t('scene.editor.tab.object'), !editor.selectedKey],
                ] as const
              ).map(([key, label, disabled]) => (
                <div
                  key={key}
                  className="relative h-full cursor-pointer text-sm leading-9"
                  style={{
                    color: disabled
                      ? 'var(--ui-select-d-color)'
                      : activeTab === key
                        ? 'var(--ui-text-color)'
                        : 'var(--ui-description-card-color)',
                    fontWeight: activeTab === key ? 600 : 400,
                    cursor: disabled ? 'not-allowed' : 'pointer',
                    borderBottom: activeTab === key ? '2px solid var(--ui-theme-color)' : '2px solid transparent',
                  }}
                  onClick={() => !disabled && setActiveTab(key)}
                >
                  {label}
                </div>
              ))}
            </div>
            <div className="min-h-0 flex-1 overflow-auto px-4 py-3">
              {activeTab === 'object' && selectedPlacement ? (
                <div>
                  {axisGroup('position')}
                  {axisGroup('rotation')}
                  {axisGroup('scale')}
                  <Button
                    danger
                    size="small"
                    icon={<TrashCan size={14} />}
                    onClick={() => editor.selectedKey && removeItem(editor.selectedKey)}
                  >
                    {t('scene.editor.remove')}
                  </Button>
                </div>
              ) : null}
              {activeTab === 'viewport' && config ? (
                <div className="flex flex-col gap-3">
                  {settingsSection(
                    'viewport',
                    t('scene.editor.viewportSection'),
                    () => {
                      const defaults = defaultSceneConfig().viewport;
                      patchViewport({ backgroundPreset: defaults.backgroundPreset, reflection: defaults.reflection });
                    },
                    <>
                      <div className="mb-1 text-xs font-medium">{t('scene.editor.background')}</div>
                      <div className="grid grid-cols-2 gap-2">
                        {bgSwatches.map((preset) => {
                          const selected = config.viewport.backgroundPreset === preset.key;
                          return (
                            <button
                              key={preset.key}
                              type="button"
                              className="flex cursor-pointer flex-col items-center gap-1.5 rounded border px-2 py-2 text-xs"
                              style={
                                selected
                                  ? {
                                      border: '1px solid var(--ui-theme-color)',
                                      background: 'var(--ui-primary-bg)',
                                      color: 'var(--ui-text-color)',
                                    }
                                  : {
                                      border: '1px solid var(--ui-header-splitter-color, #e0e0e0)',
                                      background: 'transparent',
                                      color: 'var(--ui-text-color)',
                                    }
                              }
                              onClick={() => patchViewport({ backgroundPreset: preset.key })}
                            >
                              <span
                                className="h-5 w-5 rounded-full border"
                                style={{
                                  background: preset.preview,
                                  borderColor: 'var(--ui-header-splitter-color, #e0e0e0)',
                                }}
                              />
                              {t(`scene.editor.bg.${preset.key}`)}
                            </button>
                          );
                        })}
                      </div>
                      <div
                        className="mt-2 border-t pt-2 text-xs font-medium"
                        style={{ borderColor: 'var(--ui-header-splitter-color, #e0e0e0)' }}
                      >
                        {t('scene.editor.reflectionTitle')}
                      </div>
                      {fieldRow(
                        t('scene.editor.reflectionEnabled'),
                        <div className="flex flex-1 justify-end">
                          <Switch
                            size="small"
                            checked={config.viewport.reflection.enabled}
                            onChange={(enabled) =>
                              patchViewport({ reflection: { ...config.viewport.reflection, enabled } })
                            }
                          />
                        </div>
                      )}
                      {viewportNumField(
                        t('scene.editor.reflectionClarity'),
                        config.viewport.reflection.clarity,
                        (clarity) =>
                          patchViewport({
                            reflection: { ...config.viewport.reflection, clarity: Math.min(Math.max(clarity, 0), 100) },
                          }),
                        '%'
                      )}
                      {viewportNumField(
                        t('scene.editor.reflectionDepth'),
                        config.viewport.reflection.depth,
                        (depth) =>
                          patchViewport({
                            reflection: { ...config.viewport.reflection, depth: Math.min(Math.max(depth, 0), 100) },
                          }),
                        '%'
                      )}
                    </>
                  )}
                  {settingsSection(
                    'camera',
                    t('scene.editor.cameraSection'),
                    () => patchViewport({ camera: defaultSceneConfig().viewport.camera }),
                    <>
                      {viewportNumField(
                        t('scene.editor.cameraFocal'),
                        config.viewport.camera.focalLength,
                        (focalLength) => patchViewport({ camera: { ...config.viewport.camera, focalLength } }),
                        'mm'
                      )}
                      {fieldRow(
                        t('scene.editor.cameraLensUnit'),
                        <Select
                          size="small"
                          className="flex-1"
                          value={config.viewport.camera.lensUnit}
                          options={[{ value: 'millimeters', label: t('scene.editor.lensUnitMm') }]}
                          onChange={(lensUnit) => patchViewport({ camera: { ...config.viewport.camera, lensUnit } })}
                        />
                      )}
                      {viewportNumField(t('scene.editor.cameraShiftX'), config.viewport.camera.shiftX, (shiftX) =>
                        patchViewport({ camera: { ...config.viewport.camera, shiftX } })
                      )}
                      {viewportNumField(t('scene.editor.cameraShiftY'), config.viewport.camera.shiftY, (shiftY) =>
                        patchViewport({ camera: { ...config.viewport.camera, shiftY } })
                      )}
                      {fieldRow(
                        t('scene.editor.cameraClipMode'),
                        <Select
                          size="small"
                          className="flex-1"
                          value={config.viewport.camera.clipMode}
                          options={[
                            { value: 'manual', label: t('scene.editor.clipManual') },
                            { value: 'auto', label: t('scene.editor.clipAuto') },
                          ]}
                          onChange={(clipMode) => patchViewport({ camera: { ...config.viewport.camera, clipMode } })}
                        />
                      )}
                      {viewportNumField(
                        t('scene.editor.cameraClipStart'),
                        config.viewport.camera.clipStart,
                        (clipStart) => patchViewport({ camera: { ...config.viewport.camera, clipStart } }),
                        'm',
                        config.viewport.camera.clipMode === 'auto'
                      )}
                      {viewportNumField(
                        t('scene.editor.cameraClipEnd'),
                        config.viewport.camera.clipEnd,
                        (clipEnd) => patchViewport({ camera: { ...config.viewport.camera, clipEnd } }),
                        'm',
                        config.viewport.camera.clipMode === 'auto'
                      )}
                    </>
                  )}
                  {settingsSection(
                    'grid',
                    t('scene.editor.gridSection'),
                    () => patchViewport({ grid: defaultSceneConfig().viewport.grid }),
                    <>
                      {fieldRow(
                        t('scene.editor.grid'),
                        <div className="flex flex-1 justify-end">
                          <Switch
                            size="small"
                            checked={config.viewport.grid.show}
                            onChange={(show) => patchViewport({ grid: { ...config.viewport.grid, show } })}
                          />
                        </div>
                      )}
                      {fieldRow(
                        t('scene.editor.gridAxis'),
                        <div className="flex flex-1 justify-end">
                          <Switch
                            size="small"
                            checked={config.viewport.grid.showAxis}
                            onChange={(showAxis) => patchViewport({ grid: { ...config.viewport.grid, showAxis } })}
                          />
                        </div>
                      )}
                      {viewportNumField(
                        t('scene.editor.gridSize'),
                        config.viewport.grid.areaSize,
                        (areaSize) => patchViewport({ grid: { ...config.viewport.grid, areaSize } }),
                        'm'
                      )}
                      {viewportNumField(
                        t('scene.editor.gridMinCell'),
                        config.viewport.grid.minCellSize,
                        (minCellSize) => patchViewport({ grid: { ...config.viewport.grid, minCellSize } }),
                        'm'
                      )}
                      {fieldRow(
                        t('scene.editor.gridInfinite'),
                        <div className="flex flex-1 justify-end">
                          <Switch
                            size="small"
                            checked={config.viewport.grid.infinite}
                            onChange={(infinite) => patchViewport({ grid: { ...config.viewport.grid, infinite } })}
                          />
                        </div>
                      )}
                    </>
                  )}
                </div>
              ) : null}
              {activeTab === 'light' && config ? (
                <div className="flex flex-col gap-3">
                  {settingsSection(
                    'light',
                    t('scene.editor.lightSection'),
                    () => {
                      setConfig({ ...config, light: defaultSceneConfig().light });
                      setDirty(true);
                    },
                    <>
                      <div className="text-xs font-medium">{t('scene.editor.lightPresetGroup')}</div>
                      {fieldRow(
                        t('scene.editor.lightPreset'),
                        <Select
                          size="small"
                          className="flex-1"
                          value={config.light.preset}
                          options={(['balanced', 'soft', 'contrast'] as const).map((preset) => ({
                            value: preset,
                            label: t(`scene.editor.light.${preset}`),
                          }))}
                          onChange={(preset) => {
                            // 切换预设 → 用预设值重建整份灯光配置（对齐云端 createSceneLightConfigFromPreset）
                            const values = LIGHT_PRESETS[preset] ?? LIGHT_PRESETS.balanced;
                            setConfig({
                              ...config,
                              light: {
                                preset,
                                environmentEnabled: true,
                                environmentIntensity: 1,
                                ambientIntensity: values.ambient,
                                keyLightIntensity: values.key,
                                fillLightIntensity: values.fill,
                                backLightIntensity: values.back,
                              },
                            });
                            setDirty(true);
                          }}
                        />
                      )}
                      {lightGroup(
                        t('scene.editor.iblGroup'),
                        <>
                          {fieldRow(
                            t('scene.editor.envEnabled'),
                            <div className="flex flex-1 justify-end">
                              <Switch
                                size="small"
                                checked={config.light.environmentEnabled}
                                onChange={(environmentEnabled) => {
                                  setConfig({ ...config, light: { ...config.light, environmentEnabled } });
                                  setDirty(true);
                                }}
                              />
                            </div>
                          )}
                          {lightIntensityField('environmentIntensity', 2)}
                        </>
                      )}
                      {lightGroup(t('scene.editor.ambientGroup'), lightIntensityField('ambientIntensity', 2))}
                      {lightGroup(t('scene.editor.keyGroup'), lightIntensityField('keyLightIntensity', 14))}
                      {lightGroup(t('scene.editor.fillGroup'), lightIntensityField('fillLightIntensity', 2))}
                      {lightGroup(t('scene.editor.backGroup'), lightIntensityField('backLightIntensity', 2))}
                    </>
                  )}
                </div>
              ) : null}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
