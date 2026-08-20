// Model Mapping 表格 — 布局对齐云端 motion-mapping/index.tsx：
//   Object(n) 头部（搜索 + Edit Selected）；Name | ✥ Position Binding | ⟳ Rotation Binding 三列；
//   单元格内每个轴一行：轴 chip（点击翻转正负号）+ 字段下拉（搜索/解绑）+ 实时值/系数 + 删轴/加轴。
import {
  Add,
  ChevronDown,
  CloseFilled,
  Edit,
  Link as LinkIcon,
  Move,
  RotateClockwise,
  Search,
  Unlink,
} from '@carbon/icons-react';
import { Button, Dropdown, Input, Table } from 'antd';
import { useMemo, useState } from 'react';
import type { InstanceBindingPayload, MotionMappingNode, SelectedObjectNode, SignedAxis } from '../../api/instances';
import { t } from '../../i18n';
import { resolveBindingValue } from '../../viewer/live-binding';

type TransformType = 'position' | 'rotation';
type AxisLetter = 'x' | 'y' | 'z';

// 云端轴配色：X 品红 / Y 绿 / Z 蓝
const AXIS_COLORS: Record<AxisLetter, string> = {
  x: 'var(--ui-t-magenta-color-60)',
  y: 'var(--ui-status-active-text, #53b483)',
  z: 'var(--ui-t-blue-color-70)',
};

const letterOf = (axis: SignedAxis): AxisLetter => axis[1] as AxisLetter;
const flipSign = (axis: SignedAxis): SignedAxis =>
  (axis.startsWith('-') ? `+${letterOf(axis)}` : `-${letterOf(axis)}`) as SignedAxis;

// 云端 Name 列为 RTL 后缀截断
const suffixPath = (path: string, max = 34) => (path.length > max ? `…${path.slice(-(max - 1))}` : path);

interface MappingTableProps {
  binding: InstanceBindingPayload;
  availableKeys: string[];
  payload: Record<string, unknown> | undefined;
  onChange: (updater: (prev: InstanceBindingPayload) => InstanceBindingPayload, toast?: string) => void;
  onHover: (path: string | null) => void;
  onEditSelected: () => void;
}

function AxisChip({ axis, onFlip }: { axis: SignedAxis; onFlip: () => void }) {
  const color = AXIS_COLORS[letterOf(axis)];
  return (
    <button
      type="button"
      className="h-7 shrink-0 cursor-pointer rounded border px-2 font-mono text-xs font-semibold"
      style={{ color, borderColor: color, background: 'transparent' }}
      title={t('scene.editor.mode.rotate')}
      onClick={onFlip}
    >
      {axis.toUpperCase()}
    </button>
  );
}

function PayloadDropdown({
  value,
  availableKeys,
  onSelect,
  onUnbind,
}: {
  value: string;
  availableKeys: string[];
  onSelect: (key: string) => void;
  onUnbind: () => void;
}) {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState('');
  const filtered = availableKeys.filter((key) => (search ? key.toLowerCase().includes(search.toLowerCase()) : true));
  return (
    <Dropdown
      open={open}
      onOpenChange={(next) => {
        setOpen(next);
        if (!next) setSearch('');
      }}
      trigger={['click']}
      popupRender={() => (
        <div
          className="w-56 rounded border py-1 shadow-lg"
          style={{ background: 'var(--ui-bg-color)', borderColor: 'var(--ui-header-splitter-color, #e0e0e0)' }}
        >
          <div className="flex items-center gap-1 px-2 py-1.5">
            <Input
              size="small"
              prefix={<Search size={13} />}
              placeholder={t('instance.searchObjects')}
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              onClick={(e) => e.stopPropagation()}
            />
            {value ? (
              <Button
                type="text"
                size="small"
                danger
                icon={<Unlink size={14} />}
                title={t('instance.unbind')}
                onClick={() => {
                  onUnbind();
                  setOpen(false);
                }}
              />
            ) : null}
          </div>
          <div className="max-h-48 overflow-auto">
            {filtered.length === 0 ? (
              <div className="px-3 py-2 text-xs" style={{ color: 'var(--ui-description-card-color)' }}>
                {t('instance.noResult')}
              </div>
            ) : (
              filtered.map((key) => (
                <div
                  key={key}
                  className="cursor-pointer px-3 py-1.5 text-sm hover:opacity-70"
                  style={value === key ? { color: 'var(--ui-theme-color)', fontWeight: 600 } : undefined}
                  onClick={() => {
                    onSelect(key);
                    setOpen(false);
                  }}
                >
                  {key}
                </div>
              ))
            )}
          </div>
        </div>
      )}
    >
      <button
        type="button"
        className="flex h-7 min-w-0 cursor-pointer items-center gap-1 rounded border px-2 text-xs"
        style={{ borderColor: 'var(--ui-header-splitter-color, #e0e0e0)', color: 'var(--ui-text-color)' }}
      >
        {value ? (
          <LinkIcon size={13} style={{ color: 'var(--ui-status-active-text, #53b483)' }} />
        ) : (
          <Unlink size={13} />
        )}
        <span className="max-w-24 truncate">{value || t('instance.noMapping')}</span>
        <ChevronDown size={13} />
      </button>
    </Dropdown>
  );
}

export function MappingTable({
  binding,
  availableKeys,
  payload,
  onChange,
  onHover,
  onEditSelected,
}: MappingTableProps) {
  const [search, setSearch] = useState('');

  const rows = useMemo(
    () =>
      binding.selectedObjects.filter((node) =>
        search ? node.path.toLowerCase().includes(search.toLowerCase()) : true
      ),
    [binding.selectedObjects, search]
  );

  const mappingsFor = (node: SelectedObjectNode, type: TransformType) =>
    binding.motionMappings.filter((item) => item.nodeID === node.nodeID && item.type === type);

  const patchMappings = (mutate: (list: MotionMappingNode[]) => MotionMappingNode[], toast?: string) =>
    onChange((prev) => ({ ...prev, motionMappings: mutate([...prev.motionMappings]) }), toast);

  const addAxis = (node: SelectedObjectNode, type: TransformType, letter: AxisLetter) => {
    const mapping: MotionMappingNode = {
      nodeID: node.nodeID,
      name: node.name,
      path: node.path,
      type,
      value: '',
      axis: `+${letter}` as SignedAxis,
      unit: '',
      factor: 1,
    };
    patchMappings((list) => [...list, mapping], t('instance.axisAdded', { axis: letter.toUpperCase() }));
  };

  const updateMapping = (target: MotionMappingNode, patch: Partial<MotionMappingNode>, toast?: string) =>
    patchMappings(
      (list) =>
        list.map((item) =>
          item === target ||
          (item.nodeID === target.nodeID &&
            item.type === target.type &&
            item.axis === target.axis &&
            item.value === target.value)
            ? { ...item, ...patch }
            : item
        ),
      toast
    );

  const removeMapping = (target: MotionMappingNode, toast?: string) =>
    patchMappings(
      (list) =>
        list.filter(
          (item) =>
            !(
              item.nodeID === target.nodeID &&
              item.type === target.type &&
              item.axis === target.axis &&
              item.value === target.value
            )
        ),
      toast
    );

  const bindingCell = (node: SelectedObjectNode, type: TransformType) => {
    const mappings = mappingsFor(node, type);
    const usedLetters = new Set(mappings.map((item) => letterOf(item.axis)));
    const freeLetters = (['x', 'y', 'z'] as AxisLetter[]).filter((letter) => !usedLetters.has(letter));

    const addMenu = freeLetters.length ? (
      <Dropdown
        trigger={['click']}
        menu={{
          items: freeLetters.map((letter) => ({
            key: letter,
            label: <span style={{ color: AXIS_COLORS[letter] }}>{letter.toUpperCase()}</span>,
          })),
          onClick: ({ key }) => addAxis(node, type, key as AxisLetter),
        }}
      >
        <Button type="text" size="small" icon={<Add size={15} />} title={t('instance.addAxis')} />
      </Dropdown>
    ) : null;

    if (mappings.length === 0) {
      return (
        <div className="flex items-center gap-2" style={{ minHeight: 32 }}>
          <span
            className="flex flex-1 items-center gap-1.5 text-xs"
            style={{ color: 'var(--ui-description-card-color)' }}
          >
            <Unlink size={13} />
            {t('instance.noMapping')}
          </span>
          {addMenu}
        </div>
      );
    }

    return (
      <div className="flex flex-col">
        {mappings.map((mapping, index) => {
          const live = mapping.value ? resolveBindingValue(payload, mapping.value) : undefined;
          return (
            <div
              key={`${mapping.axis}-${index}`}
              className="flex items-center gap-2"
              style={{ minHeight: 32, marginTop: index > 0 ? 8 : 0 }}
            >
              <AxisChip
                axis={mapping.axis}
                onFlip={() => updateMapping(mapping, { axis: flipSign(mapping.axis) }, t('instance.signOk'))}
              />
              <PayloadDropdown
                value={mapping.value}
                availableKeys={availableKeys}
                onSelect={(key) => updateMapping(mapping, { value: key }, t('instance.mapSaved'))}
                onUnbind={() => updateMapping(mapping, { value: '' }, t('instance.mapRemoved'))}
              />
              {mapping.value ? (
                <div
                  className="flex shrink-0 items-center gap-1.5 rounded px-2 py-0.5"
                  style={{ background: 'var(--ui-card-bg)' }}
                >
                  <span
                    className="font-mono text-xs"
                    style={{
                      color: live !== undefined ? 'var(--ui-text-color)' : 'var(--ui-description-card-color)',
                    }}
                  >
                    {live !== undefined ? live.toFixed(2) : '--'}
                  </span>
                  {type === 'position' ? (
                    <Input
                      size="small"
                      className="h-6"
                      style={{ width: 64 }}
                      prefix={<span className="text-xs">x</span>}
                      defaultValue={String(mapping.factor ?? 1)}
                      onBlur={(e) => {
                        const parsed = Number(e.target.value);
                        updateMapping(
                          mapping,
                          { factor: Number.isFinite(parsed) && parsed !== 0 ? parsed : 1 },
                          t('instance.factorOk')
                        );
                      }}
                    />
                  ) : null}
                </div>
              ) : null}
              <div className="ml-auto flex shrink-0 items-center">
                <Button
                  type="text"
                  size="small"
                  icon={<CloseFilled size={14} style={{ color: 'var(--ui-select-d-color)' }} />}
                  onClick={() => removeMapping(mapping, t('instance.axisRemoved'))}
                />
                {addMenu}
              </div>
            </div>
          );
        })}
      </div>
    );
  };

  return (
    <div>
      {/* Object(n) 头部：计数 + 搜索 + Edit Selected */}
      <div
        className="flex items-center justify-between gap-4 border-b px-1 py-2"
        style={{ borderColor: 'var(--ui-header-splitter-color, #e0e0e0)' }}
      >
        <span className="shrink-0 text-sm font-medium">
          {t('instance.objectCount', { count: binding.selectedObjects.length })}
        </span>
        <div className="flex items-center gap-2">
          <Input
            size="small"
            allowClear
            style={{ width: 180 }}
            prefix={<Search size={13} />}
            placeholder={t('instance.searchObjects')}
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
          <Button size="small" icon={<Edit size={13} />} onClick={onEditSelected}>
            {t('instance.editSelected')}
          </Button>
        </div>
      </div>

      <Table
        rowKey={(record) => record.nodeID}
        size="small"
        pagination={false}
        dataSource={rows}
        onRow={(record) => ({
          onMouseEnter: () => onHover(record.path),
          onMouseLeave: () => onHover(null),
        })}
        columns={[
          {
            title: t('instance.mapping.name'),
            dataIndex: 'path',
            width: '24%',
            render: (path: string) => (
              <span className="font-mono text-xs" title={path}>
                {suffixPath(path)}
              </span>
            ),
          },
          {
            title: (
              <span className="flex items-center gap-1">
                <Move size={13} /> {t('instance.mapping.position')}
              </span>
            ),
            width: '38%',
            render: (_: unknown, record) => bindingCell(record, 'position'),
          },
          {
            title: (
              <span className="flex items-center gap-1">
                <RotateClockwise size={13} /> {t('instance.mapping.rotation')}
              </span>
            ),
            width: '38%',
            render: (_: unknown, record) => bindingCell(record, 'rotation'),
          },
        ]}
      />
    </div>
  );
}
