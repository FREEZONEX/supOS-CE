import { useEffect, useRef, useState } from 'react';
import type { InstanceDataTag } from '../../api/instances';
import { formatDataTagPreviewValue, layoutDataTagPreviewOverlays } from './tag-overlay-layout';

// 投影只需要 ViewerApi 的子集，GLB 查看器与 SPLAT 查看器都满足
export interface TagOverlayViewer {
  status: string;
  projectModelPoint: (local: { x: number; y: number; z: number }) => { x: number; y: number } | null;
  projectNodePoint: (path: string, local: { x: number; y: number; z: number }) => { x: number; y: number } | null;
}

// 数据标签 3D 标点（对齐源端 data-tag-visuals）：
//   静态标签 = 绿色圆环 marker（可点选）+ 虚线引线 + 数值面板（名称 + 实时值）；
//   运动标签 = 锚点小圆点 + 同款面板（key + 实时值），按可见性设置过滤；
//   面板按源端布局算法分列画布两侧避让，超时数据置灰（stale）。
const STALE_AFTER_MS = 5_000;
const PANEL_HEIGHT = 50;

export interface MotionOverlayLabel {
  key: string;
  path: string;
}

interface OverlayItem {
  id: string;
  kind: 'tag' | 'motion';
  title: string;
  value: string;
  screen: { x: number; y: number };
}

function estimatePanelWidth(title: string, value: string) {
  const chars = Math.max(title.length, value.length);
  return Math.max(92, Math.min(180, 24 + chars * 6.5));
}

export function TagOverlay({
  tags,
  motionLabels = [],
  viewer,
  payload,
  messageTs,
  selectedTagId,
  onSelectTag,
  suspended = false,
}: {
  tags: InstanceDataTag[];
  motionLabels?: MotionOverlayLabel[];
  viewer: TagOverlayViewer;
  payload: Record<string, unknown> | undefined;
  messageTs?: number;
  selectedTagId?: string | null;
  onSelectTag?: (id: string) => void;
  /** 页签隐藏/窗口后台时暂停投影 RAF 循环 */
  suspended?: boolean;
}) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [positions, setPositions] = useState<Record<string, { x: number; y: number } | null>>({});
  const [size, setSize] = useState({ width: 0, height: 0 });
  // stale 依赖 Date.now()，不能在渲染期计算 —— 在 rAF 循环里评估
  const [timedOut, setTimedOut] = useState(false);
  const messageTsRef = useRef<number | undefined>(undefined);
  useEffect(() => {
    messageTsRef.current = messageTs;
  }, [messageTs]);

  useEffect(() => {
    if (suspended || viewer.status !== 'ready' || (tags.length === 0 && motionLabels.length === 0)) {
      setPositions({});
      return;
    }
    let frameId = 0;
    const update = () => {
      frameId = requestAnimationFrame(update);
      const next: Record<string, { x: number; y: number } | null> = {};
      for (const tag of tags) {
        // 标签坐标是模型根局部空间（对齐源端 data-tag-space），path 仅元数据，不参与定位
        next[`tag:${tag.id}`] = viewer.projectModelPoint({ x: tag.x ?? 0, y: tag.y ?? 0, z: tag.z ?? 0 });
      }
      for (const label of motionLabels) {
        next[`motion:${label.key}`] = viewer.projectNodePoint(label.path, { x: 0, y: 0, z: 0 });
      }
      setPositions(next);
      const ts = messageTsRef.current;
      setTimedOut(ts !== undefined && Date.now() - ts >= STALE_AFTER_MS);
      const el = containerRef.current;
      if (el && (el.clientWidth !== size.width || el.clientHeight !== size.height)) {
        setSize({ width: el.clientWidth, height: el.clientHeight });
      }
    };
    update();
    return () => cancelAnimationFrame(frameId);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [viewer, viewer.status, tags, motionLabels, suspended]);

  const stale = !payload || timedOut;

  const items: OverlayItem[] = [];
  for (const tag of tags) {
    const screen = positions[`tag:${tag.id}`];
    if (!screen) continue;
    const raw = tag.payload && payload ? payload[tag.payload] : undefined;
    const display = formatDataTagPreviewValue(raw);
    items.push({
      id: `tag:${tag.id}`,
      kind: 'tag',
      title: tag.name,
      value: display === '--' ? display : `${display}${tag.unit ? ` ${tag.unit}` : ''}`,
      screen,
    });
  }
  for (const label of motionLabels) {
    const screen = positions[`motion:${label.key}`];
    if (!screen) continue;
    items.push({
      id: `motion:${label.key}`,
      kind: 'motion',
      title: label.key,
      value: formatDataTagPreviewValue(payload ? payload[label.key] : undefined),
      screen,
    });
  }

  const layout =
    size.width > 0
      ? layoutDataTagPreviewOverlays({
          overlayWidth: size.width,
          overlayHeight: size.height,
          anchors: items.map((item) => ({
            id: item.id,
            screenX: item.screen.x,
            screenY: item.screen.y,
            width: estimatePanelWidth(item.title, item.value),
            height: PANEL_HEIGHT,
          })),
        })
      : [];
  const layoutById = new Map(layout.map((entry) => [entry.id, entry]));

  const themeColor = 'var(--ui-theme-color)';
  const panelBorder = stale ? 'var(--ui-header-splitter-color, #e0e0e0)' : themeColor;
  const panelBg = stale ? 'var(--ui-bg-color)' : 'var(--ui-primary-bg)';

  return (
    <div ref={containerRef} className="pointer-events-none absolute inset-0 overflow-hidden">
      {items.map((item) => {
        const entry = layoutById.get(item.id);
        const isTagMarker = item.kind === 'tag';
        const tagId = isTagMarker ? item.id.slice(4) : null;
        const selected = tagId !== null && tagId === selectedTagId;
        return (
          <div key={item.id} className="absolute" style={{ left: item.screen.x, top: item.screen.y }}>
            {/* 引线（虚线，从锚点指向面板） */}
            {entry ? (
              <div
                className="absolute"
                style={{
                  width: entry.lineLength,
                  borderTop: `2px dashed ${panelBorder}`,
                  transform: `rotate(${entry.angle}rad)`,
                  transformOrigin: '0 50%',
                }}
              />
            ) : null}
            {/* 锚点：静态标签是可点选的圆环 marker，运动标签是小圆点 */}
            {isTagMarker ? (
              <button
                type="button"
                className="pointer-events-auto absolute -translate-x-1/2 -translate-y-1/2 cursor-pointer rounded-full"
                style={{
                  width: selected ? 22 : 18,
                  height: selected ? 22 : 18,
                  background: 'var(--ui-bg-color, #fff)',
                  border: `${selected ? 7 : 5}px solid ${themeColor}`,
                  boxShadow: '0 1px 4px rgba(0,0,0,0.25)',
                }}
                onClick={() => tagId && onSelectTag?.(tagId)}
              />
            ) : (
              <div
                className="absolute h-1.5 w-1.5 -translate-x-1/2 -translate-y-1/2 rounded-full"
                style={{ background: themeColor }}
              />
            )}
            {/* 数值面板：名称/key + 实时值 */}
            {entry ? (
              <div
                className="absolute whitespace-nowrap"
                style={{
                  left: entry.panelOffsetX,
                  top: entry.panelOffsetY,
                  minWidth: 92,
                  maxWidth: 180,
                  padding: '8px 10px',
                  borderRadius: 8,
                  border: `2px solid ${panelBorder}`,
                  background: panelBg,
                  boxShadow: '0 4px 12px rgba(15, 23, 42, 0.08)',
                }}
              >
                <div
                  className="overflow-hidden text-ellipsis text-[10px] leading-3"
                  style={{ color: 'var(--ui-description-card-color)' }}
                >
                  {item.title}
                </div>
                <div
                  className="mt-0.5 overflow-hidden text-ellipsis text-xs leading-3"
                  style={{ color: 'var(--ui-text-color)' }}
                >
                  {item.value}
                </div>
              </div>
            ) : null}
          </div>
        );
      })}
    </div>
  );
}
