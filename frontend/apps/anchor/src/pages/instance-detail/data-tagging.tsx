// Data Tagging tab — 布局对齐云端 data-tagging-editor：
//   "静态标签" 标题 + 添加标签按钮（placement 模式点 3D 直接创建，无弹窗）；
//   160px 标签卡片流（超 2 行折叠 + 展开全部）；选中标签编辑面板（名称/XYZ 坐标可编辑/字段/单位）；
//   Motion Tags 可见性控制（全部显示/全部隐藏/自定义 + 每 key 是/否），持久化到 binding.motionTagVisibility。
import { Add, ChevronDown, ChevronUp, TrashCan } from '@carbon/icons-react';
import { Button, Empty, Input, InputNumber, Segmented, Select } from 'antd';
import dayjs from 'dayjs';
import { useEffect, useLayoutEffect, useMemo, useRef, useState } from 'react';
import type { InstanceBindingPayload, InstanceDataTag, MotionVisibilitySettings } from '../../api/instances';
import { t } from '../../i18n';

const uniqueId = () => `tag_${Date.now().toString(36)}${Math.random().toString(36).slice(2, 8)}`;
const borderStyle = { borderColor: 'var(--ui-header-splitter-color, #e0e0e0)' } as const;
const mutedColor = { color: 'var(--ui-description-card-color)' } as const;
const COLLAPSED_ROWS = 2;

function formatCreatedAt(value?: string) {
  if (!value) return '';
  const date = dayjs(value);
  if (!date.isValid()) return value;
  return date.format('HH:mm DD/MM/YYYY');
}

export function DataTaggingPanel({
  binding,
  availableKeys,
  placementMode,
  onPlacementModeChange,
  onChange,
  selectedTagId,
  onSelectTag,
}: {
  binding: InstanceBindingPayload;
  availableKeys: string[];
  placementMode: boolean;
  onPlacementModeChange: (next: boolean) => void;
  onChange: (updater: (prev: InstanceBindingPayload) => InstanceBindingPayload, toast?: string) => void;
  selectedTagId: string | null;
  onSelectTag: (id: string | null) => void;
}) {
  const selectedTag = binding.dataTags.find((tag) => tag.id === selectedTagId) ?? null;
  const [tagsExpanded, setTagsExpanded] = useState(false);
  const [needsShowAll, setNeedsShowAll] = useState(false);
  const [collapsedHeight, setCollapsedHeight] = useState(0);
  const cardsRef = useRef<HTMLDivElement>(null);

  const motionKeys = useMemo(() => {
    const seen = new Set<string>();
    const keys: string[] = [];
    for (const mapping of binding.motionMappings) {
      if (mapping.value && !seen.has(mapping.value)) {
        seen.add(mapping.value);
        keys.push(mapping.value);
      }
    }
    return keys;
  }, [binding.motionMappings]);

  const visibility: MotionVisibilitySettings = binding.motionTagVisibility ?? { mode: 'showAll', custom: {} };

  // viewer 点击（placement 模式）→ 直接创建标签并选中（云端行为：无弹窗，默认绑定第一个字段）
  useEffect(() => {
    const listener = (event: Event) => {
      const detail = (event as CustomEvent).detail as { path: string; local: { x: number; y: number; z: number } };
      const tag: InstanceDataTag = {
        id: uniqueId(),
        nodeID: detail.path,
        name: `tag ${String(binding.dataTags.length + 1).padStart(2, '0')}`,
        path: detail.path,
        x: detail.local.x,
        y: detail.local.y,
        z: detail.local.z,
        payload: availableKeys[0] || '',
        unit: '',
        createdAt: new Date().toISOString(),
      };
      onChange((prev) => ({ ...prev, dataTags: [...prev.dataTags, tag] }), t('instance.tags.saved'));
      onSelectTag(tag.id);
      onPlacementModeChange(false);
    };
    window.addEventListener('anchor:create-tag', listener);
    return () => window.removeEventListener('anchor:create-tag', listener);
  }, [availableKeys, binding.dataTags.length, onChange, onPlacementModeChange, onSelectTag]);

  // 测量卡片实际行数，决定是否需要"展开全部"（对齐云端折叠逻辑）
  useLayoutEffect(() => {
    const el = cardsRef.current;
    if (!el || binding.dataTags.length === 0) {
      setNeedsShowAll(false);
      return;
    }
    const items = Array.from(el.children) as HTMLElement[];
    if (items.length === 0) {
      setNeedsShowAll(false);
      return;
    }
    const rowTops = [...new Set(items.map((item) => item.offsetTop))].sort((a, b) => a - b);
    const moreThanCollapsedRows = rowTops.length > COLLAPSED_ROWS;
    setNeedsShowAll(moreThanCollapsedRows);
    if (moreThanCollapsedRows && rowTops.length >= COLLAPSED_ROWS) {
      const lastCollapsedRowItems = items.filter((item) => item.offsetTop === rowTops[COLLAPSED_ROWS - 1]);
      const bottomOfLastRow = Math.max(...lastCollapsedRowItems.map((item) => item.offsetTop + item.offsetHeight));
      setCollapsedHeight(bottomOfLastRow);
    }
  }, [binding.dataTags]);

  const updateTag = (id: string, patch: Partial<InstanceDataTag>, toast?: string) =>
    onChange(
      (prev) => ({ ...prev, dataTags: prev.dataTags.map((tag) => (tag.id === id ? { ...tag, ...patch } : tag)) }),
      toast
    );

  const deleteTag = (id: string) => {
    onChange((prev) => ({ ...prev, dataTags: prev.dataTags.filter((tag) => tag.id !== id) }), t('instance.tags.saved'));
    if (selectedTagId === id) onSelectTag(null);
  };

  const setVisibility = (next: MotionVisibilitySettings) =>
    onChange((prev) => ({ ...prev, motionTagVisibility: next }));

  const fieldBlock = (label: string, children: React.ReactNode) => (
    <div>
      <div className="mb-1.5 text-sm font-medium">{label}</div>
      {children}
    </div>
  );

  return (
    <div className="pt-1">
      {/* 标题 + 添加标签 */}
      <div className="flex items-center justify-between py-2">
        <div className="text-base font-medium">{t('instance.tags.title')}</div>
        <Button type="primary" icon={<Add size={16} />} onClick={() => onPlacementModeChange(!placementMode)}>
          {t('instance.tags.add')}
        </Button>
      </div>
      {placementMode ? (
        <div
          className="mb-2 rounded border border-dashed px-3 py-1.5 text-xs"
          style={{ ...borderStyle, color: 'var(--ui-theme-color)' }}
        >
          {t('instance.tags.placing')}
        </div>
      ) : null}

      {/* 标签卡片流（超 2 行折叠） */}
      {binding.dataTags.length === 0 ? (
        <div className="py-8">
          <Empty
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            description={
              <div>
                <div className="text-sm">{t('instance.tags.emptyTitle')}</div>
                <div className="mt-1 text-xs" style={mutedColor}>
                  {t('instance.tags.emptyDesc')}
                </div>
              </div>
            }
          >
            <Button type="primary" icon={<Add size={16} />} onClick={() => onPlacementModeChange(true)}>
              {t('instance.tags.clickStart')}
            </Button>
          </Empty>
        </div>
      ) : (
        <div className="py-2">
          <div
            ref={cardsRef}
            className="relative flex flex-wrap gap-3 overflow-hidden"
            style={needsShowAll && !tagsExpanded ? { maxHeight: collapsedHeight } : undefined}
          >
            {binding.dataTags.map((tag) => {
              const isSelected = tag.id === selectedTagId;
              return (
                <div
                  key={tag.id}
                  role="button"
                  tabIndex={0}
                  className="group min-w-[160px] max-w-[160px] cursor-pointer rounded-lg border px-4 py-3 text-left"
                  style={
                    isSelected
                      ? { border: '1px solid var(--ui-theme-color)', background: 'var(--ui-primary-bg)' }
                      : borderStyle
                  }
                  onClick={() => {
                    onSelectTag(tag.id);
                    onPlacementModeChange(false);
                  }}
                  onKeyDown={(event) => {
                    if (event.key === 'Enter' || event.key === ' ') {
                      event.preventDefault();
                      onSelectTag(tag.id);
                      onPlacementModeChange(false);
                    }
                  }}
                >
                  <div className="flex items-start justify-between gap-3">
                    <div className="min-w-0">
                      <div className="truncate text-sm font-semibold" title={tag.name}>
                        {tag.name}
                      </div>
                      <div className="mt-1 text-xs" style={mutedColor}>
                        {formatCreatedAt(tag.createdAt)}
                      </div>
                    </div>
                    <Button
                      type="text"
                      size="small"
                      danger
                      className="shrink-0 opacity-0 group-hover:opacity-100"
                      style={isSelected ? { opacity: 1 } : undefined}
                      icon={<TrashCan size={14} />}
                      onClick={(e) => {
                        e.stopPropagation();
                        deleteTag(tag.id);
                      }}
                    />
                  </div>
                </div>
              );
            })}
          </div>
          {needsShowAll ? (
            <Button
              type="text"
              block
              className="mt-2"
              icon={tagsExpanded ? <ChevronUp size={16} /> : <ChevronDown size={16} />}
              onClick={() => setTagsExpanded((value) => !value)}
            >
              {tagsExpanded ? t('instance.tags.showLess') : t('instance.tags.showAll')}
            </Button>
          ) : null}
        </div>
      )}

      {/* 选中标签的属性编辑面板（对齐云端：Name / XYZ 可编辑 / Payload / Unit） */}
      {selectedTag ? (
        <div className="mt-2 space-y-4 rounded border p-4" style={borderStyle}>
          {fieldBlock(
            t('instance.tagging.name'),
            <Input
              value={selectedTag.name}
              onChange={(e) => updateTag(selectedTag.id, { name: e.target.value })}
              onBlur={() => onChange((prev) => prev, t('instance.tags.saved'))}
            />
          )}
          <div className="grid grid-cols-3 gap-3">
            {(['x', 'y', 'z'] as const).map((axis) => (
              <div key={axis}>
                <div className="mb-1.5 text-sm font-medium uppercase">{axis}</div>
                <InputNumber
                  step={0.001}
                  style={{ width: '100%' }}
                  value={Number(selectedTag[axis] ?? 0)}
                  onChange={(value) => {
                    if (value === null || !Number.isFinite(value)) return;
                    updateTag(selectedTag.id, { [axis]: value });
                  }}
                  onBlur={() => onChange((prev) => prev, t('instance.tags.saved'))}
                />
              </div>
            ))}
          </div>
          {fieldBlock(
            t('instance.tagging.payload'),
            <Select
              style={{ width: '100%' }}
              showSearch
              allowClear
              value={selectedTag.payload || undefined}
              options={availableKeys.map((key) => ({ value: key, label: key }))}
              disabled={availableKeys.length === 0}
              onChange={(value) => updateTag(selectedTag.id, { payload: value || '' }, t('instance.tags.saved'))}
            />
          )}
          {fieldBlock(
            t('instance.tagging.unit'),
            <Input
              value={selectedTag.unit}
              placeholder="°C"
              onChange={(e) => updateTag(selectedTag.id, { unit: e.target.value })}
              onBlur={() => onChange((prev) => prev, t('instance.tags.saved'))}
            />
          )}
        </div>
      ) : binding.dataTags.length > 0 ? (
        <div
          className="mt-2 flex min-h-[96px] items-center justify-center rounded-lg border border-dashed text-sm"
          style={{ ...borderStyle, ...mutedColor }}
        >
          {t('instance.tags.selectHint')}
        </div>
      ) : null}

      {/* Motion Tags 可见性（对齐云端：模式切换 + 每 key 是/否，随 binding 持久化） */}
      {motionKeys.length > 0 ? (
        <div className="mt-6 pb-4">
          <div className="mb-3 text-base font-medium">{t('instance.tags.motionTitle')}</div>
          <div className="rounded border px-4" style={borderStyle}>
            <div className="grid grid-cols-2 items-center border-b border-dashed py-4" style={borderStyle}>
              <span className="text-sm font-medium">{t('instance.tags.visMode')}</span>
              <Segmented
                // grid 子项默认被拉伸到整列，控件只需包住文字
                className="justify-self-start"
                value={visibility.mode}
                options={[
                  { value: 'showAll', label: t('instance.tags.visShowAll') },
                  { value: 'hideAll', label: t('instance.tags.visHideAll') },
                  { value: 'custom', label: t('instance.tags.visCustom') },
                ]}
                onChange={(mode) => setVisibility({ ...visibility, mode: mode as MotionVisibilitySettings['mode'] })}
              />
            </div>
            <div className="grid grid-cols-2 py-2">
              <span className="text-sm font-medium" style={mutedColor}>
                {t('instance.tags.keyCol')}
              </span>
              <span className="text-sm font-medium" style={mutedColor}>
                {t('instance.tags.visCol')}
              </span>
            </div>
            {motionKeys.map((key) => {
              const isCustomMode = visibility.mode === 'custom';
              const visible =
                visibility.mode === 'showAll'
                  ? true
                  : visibility.mode === 'hideAll'
                    ? false
                    : (visibility.custom[key] ?? true);
              return (
                <div key={key} className="grid grid-cols-2 items-center py-1.5">
                  <span className="truncate pr-4 font-mono text-sm font-medium" title={key}>
                    {key}
                  </span>
                  <Segmented
                    className="justify-self-start"
                    value={visible ? 'yes' : 'no'}
                    disabled={!isCustomMode}
                    options={[
                      { value: 'yes', label: t('instance.tags.visYes') },
                      { value: 'no', label: t('instance.tags.visNo') },
                    ]}
                    onChange={(value) => {
                      if (!isCustomMode) return;
                      setVisibility({ ...visibility, custom: { ...visibility.custom, [key]: value === 'yes' } });
                    }}
                  />
                </div>
              );
            })}
          </div>
        </div>
      ) : null}
    </div>
  );
}
