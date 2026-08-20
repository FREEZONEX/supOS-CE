import { useEffect, useRef, useState, type ReactNode } from 'react';
import { Input, Switch } from 'antd';
import type { VisionCameraRegion, VisionPoint, VisionZone } from '@/apis/core-api/task';
import { Add, Checkmark, ChevronLeft, ChevronRight, Edit, TrashCan, Undo } from '@/components/lucide-icon/carbon';
import { useTranslate } from '@/hooks';
import styles from './index.module.scss';

type RegionEditorProps = {
  snapshotUrl: string;
  mode: 'zones' | 'line';
  value: VisionCameraRegion;
  onChange: (region: VisionCameraRegion) => void;
  // 画布左上角的摄像头切换器与右下角的额外操作,由上层注入(多相机时才有)。
  cameraPicker?: ReactNode;
  extraActions?: ReactNode;
};

type DragTarget = { kind: 'zonePoint'; zone: number; point: number } | { kind: 'lineStart' } | { kind: 'lineEnd' };

const clamp01 = (v: number) => Math.min(1, Math.max(0, v));

const RegionEditor = ({ snapshotUrl, mode, value, onChange, cameraPicker, extraActions }: RegionEditorProps) => {
  const formatMessage = useTranslate();
  const wrapRef = useRef<HTMLDivElement>(null);
  const [drag, setDrag] = useState<DragTarget | null>(null);
  const [draft, setDraft] = useState<VisionPoint[]>([]); // 正在绘制的多边形
  const [drawing, setDrawing] = useState(false); // 是否处于「新增区域」的绘制态
  const [toolbarOpen, setToolbarOpen] = useState(true);
  const [imgError, setImgError] = useState(false);

  useEffect(() => {
    setImgError(false);
  }, [snapshotUrl]);

  const zones = value.zones || [];
  const line = value.countingLine || null;

  const toNorm = (clientX: number, clientY: number): VisionPoint => {
    const rect = wrapRef.current!.getBoundingClientRect();
    return { x: clamp01((clientX - rect.left) / rect.width), y: clamp01((clientY - rect.top) / rect.height) };
  };

  const handleClick = (e: React.MouseEvent) => {
    if (drag) return;
    const p = toNorm(e.clientX, e.clientY);
    if (mode === 'zones') {
      // 画布只在「新增区域」模式下响应打点,平时点击不误加顶点。
      if (drawing) setDraft((d) => [...d, p]);
    } else if (drawing) {
      // 计数线:第一次点定起点、第二次定终点;已有两点时再点表示从这里重画一条。
      setDraft((d) => (d.length >= 2 ? [p] : [...d, p]));
    }
  };

  const handlePointerMove = (e: React.MouseEvent) => {
    if (!drag) return;
    const p = toNorm(e.clientX, e.clientY);
    if (drag.kind === 'zonePoint') {
      const next = zones.map((z, zi) =>
        zi === drag.zone ? { ...z, points: z.points.map((pt, pi) => (pi === drag.point ? p : pt)) } : z
      );
      onChange({ ...value, zones: next });
    } else if (line) {
      onChange({ ...value, countingLine: drag.kind === 'lineStart' ? { ...line, start: p } : { ...line, end: p } });
    }
  };

  // 计数线用两点提交,覆盖原有的线。
  const commitLine = () => {
    if (draft.length < 2) return;
    onChange({ ...value, countingLine: { start: draft[0], end: draft[1] } });
    setDraft([]);
    setDrawing(false);
  };

  const commitDraft = () => {
    if (draft.length < 3) return;
    const name = `${formatMessage('Vision.region.zone')} ${zones.length + 1}`;
    onChange({ ...value, zones: [...zones, { name, enabled: true, points: draft }] });
    setDraft([]);
    setDrawing(false);
  };

  const undoPoint = () => setDraft((d) => d.slice(0, -1));

  const cancelDraw = () => {
    setDraft([]);
    setDrawing(false);
  };

  const removeZone = (zi: number) => onChange({ ...value, zones: zones.filter((_, i) => i !== zi) });
  const patchZone = (zi: number, patch: Partial<VisionZone>) =>
    onChange({ ...value, zones: zones.map((z, i) => (i === zi ? { ...z, ...patch } : z)) });
  const pct = (p: VisionPoint) => ({ cx: `${p.x * 100}%`, cy: `${p.y * 100}%` });
  const polyPoints = (pts: VisionPoint[]) => pts.map((p) => `${p.x * 100},${p.y * 100}`).join(' ');

  return (
    <div className={styles.regionEditor}>
      <div
        ref={wrapRef}
        className={`${styles.regionStage} ${drawing ? styles.regionStageDrawing : ''}`}
        onClick={handleClick}
        onMouseMove={handlePointerMove}
        onMouseUp={() => setDrag(null)}
        onMouseLeave={() => setDrag(null)}
      >
        {imgError ? (
          <div className={styles.regionNoImg}>{formatMessage('Vision.region.noSnapshot')}</div>
        ) : (
          <img
            src={snapshotUrl}
            alt="snapshot"
            className={styles.regionImg}
            draggable={false}
            onError={() => setImgError(true)}
          />
        )}
        <svg className={styles.regionSvg} viewBox="0 0 100 100" preserveAspectRatio="none">
          {/* 已完成多边形 */}
          {mode === 'zones' &&
            zones.map((z: VisionZone, zi) => (
              <polygon
                key={zi}
                points={polyPoints(z.points)}
                className={z.enabled === false ? styles.regionPolyOff : styles.regionPoly}
              />
            ))}
          {/* 绘制中的多边形 */}
          {mode === 'zones' && draft.length > 0 && (
            <polyline points={polyPoints(draft)} className={styles.regionDraft} />
          )}
          {/* 计数线:未设置时画一条默认中线;绘制态画黄色虚线草稿;否则画已提交的绿实线 */}
          {mode === 'line' && !drawing && !line && (
            <line x1={50} y1={0} x2={50} y2={100} className={styles.regionLineDefault} />
          )}
          {mode === 'line' && drawing && draft.length === 2 && (
            <line
              x1={draft[0].x * 100}
              y1={draft[0].y * 100}
              x2={draft[1].x * 100}
              y2={draft[1].y * 100}
              className={styles.regionDraft}
            />
          )}
          {mode === 'line' && !drawing && line && (
            <line
              x1={line.start.x * 100}
              y1={line.start.y * 100}
              x2={line.end.x * 100}
              y2={line.end.y * 100}
              className={styles.regionLine}
            />
          )}
        </svg>
        {/* 端点/顶点(用 HTML 层放置,便于拖拽) */}
        <svg className={styles.regionSvg}>
          {mode === 'zones' &&
            zones.flatMap((z, zi) =>
              z.points.map((p, pi) => (
                <circle
                  key={`${zi}-${pi}`}
                  {...pct(p)}
                  r={5}
                  className={styles.regionHandle}
                  onMouseDown={(e) => {
                    e.stopPropagation();
                    setDrag({ kind: 'zonePoint', zone: zi, point: pi });
                  }}
                />
              ))
            )}
          {mode === 'zones' &&
            draft.map((p, pi) => <circle key={`d-${pi}`} {...pct(p)} r={4} className={styles.regionDraftHandle} />)}
          {mode === 'line' &&
            drawing &&
            draft.map((p, pi) => <circle key={`ld-${pi}`} {...pct(p)} r={5} className={styles.regionDraftHandle} />)}
          {mode === 'line' && !drawing && line && (
            <>
              <circle
                {...pct(line.start)}
                r={6}
                className={styles.regionHandle}
                onMouseDown={(e) => {
                  e.stopPropagation();
                  setDrag({ kind: 'lineStart' });
                }}
              />
              <circle
                {...pct(line.end)}
                r={6}
                className={styles.regionHandle}
                onMouseDown={(e) => {
                  e.stopPropagation();
                  setDrag({ kind: 'lineEnd' });
                }}
              />
            </>
          )}
        </svg>
        {/* 每个区域在画面上标出自己的名字,锚在多边形最上方的顶点 */}
        {mode === 'zones' &&
          zones.map((z: VisionZone, zi) => {
            const top = z.points.reduce((a, b) => (b.y < a.y ? b : a), z.points[0]);
            if (!top) return null;
            return (
              <span
                key={`label-${zi}`}
                className={styles.regionZoneLabel}
                style={{ left: `${top.x * 100}%`, top: `${top.y * 100}%` }}
              >
                {z.name ?? `${formatMessage('Vision.region.zone')} ${zi + 1}`}
              </span>
            );
          })}
        {/* 左上角摄像头切换、右下角操作条,都浮在画面上(对齐设计稿) */}
        {cameraPicker && (
          <div className={styles.regionStagePicker} onClick={(e) => e.stopPropagation()}>
            {cameraPicker}
          </div>
        )}
        {/* 右下角操作条按设计稿分三态:收起 / 空闲(新增+复制) / 绘制中(回退+完成+取消) */}
        <div className={styles.regionStageActions} onClick={(e) => e.stopPropagation()}>
          <button
            type="button"
            className={styles.regionStageBtn}
            title={formatMessage(toolbarOpen ? 'Vision.region.collapseTools' : 'Vision.region.expandTools')}
            aria-label={formatMessage(toolbarOpen ? 'Vision.region.collapseTools' : 'Vision.region.expandTools')}
            onClick={() => setToolbarOpen((v) => !v)}
          >
            {toolbarOpen ? <ChevronRight size={14} /> : <ChevronLeft size={14} />}
          </button>
          {toolbarOpen && mode === 'zones' && !drawing && (
            <>
              <button
                type="button"
                className={styles.regionStageBtn}
                title={formatMessage('Vision.region.addZone')}
                aria-label={formatMessage('Vision.region.addZone')}
                onClick={() => setDrawing(true)}
              >
                <Add size={14} />
              </button>
              {extraActions}
            </>
          )}
          {toolbarOpen && mode === 'zones' && drawing && (
            <>
              <button
                type="button"
                className={styles.regionStageBtn}
                disabled={draft.length === 0}
                title={formatMessage('Vision.region.undoPoint')}
                aria-label={formatMessage('Vision.region.undoPoint')}
                onClick={undoPoint}
              >
                <Undo size={14} />
              </button>
              <button
                type="button"
                className={`${styles.regionStageBtn} ${styles.regionStageBtnLabeled}`}
                disabled={draft.length < 3}
                title={`${formatMessage('Vision.region.finishZone')} (${draft.length})`}
                aria-label={formatMessage('Vision.region.finishZone')}
                onClick={commitDraft}
              >
                <Checkmark size={14} />
                {formatMessage('Vision.region.finishZone')}
              </button>
              <button
                type="button"
                className={styles.regionStageBtn}
                title={formatMessage('common.cancel')}
                aria-label={formatMessage('common.cancel')}
                onClick={cancelDraw}
              >
                <TrashCan size={14} />
              </button>
            </>
          )}
          {toolbarOpen && mode === 'line' && !drawing && (
            <>
              <button
                type="button"
                className={styles.regionStageBtn}
                title={formatMessage('Vision.region.addLine')}
                aria-label={formatMessage('Vision.region.addLine')}
                onClick={() => {
                  setDraft([]);
                  setDrawing(true);
                }}
              >
                <Edit size={14} />
              </button>
              {extraActions}
            </>
          )}
          {toolbarOpen && mode === 'line' && drawing && (
            <>
              <button
                type="button"
                className={styles.regionStageBtn}
                disabled={draft.length === 0}
                title={formatMessage('Vision.region.undoPoint')}
                aria-label={formatMessage('Vision.region.undoPoint')}
                onClick={undoPoint}
              >
                <Undo size={14} />
              </button>
              <button
                type="button"
                className={`${styles.regionStageBtn} ${styles.regionStageBtnLabeled}`}
                disabled={draft.length < 2}
                title={formatMessage('Vision.region.finishZone')}
                aria-label={formatMessage('Vision.region.finishZone')}
                onClick={commitLine}
              >
                <Checkmark size={14} />
                {formatMessage('Vision.region.finishZone')}
              </button>
              <button
                type="button"
                className={styles.regionStageBtn}
                title={formatMessage('common.cancel')}
                aria-label={formatMessage('common.cancel')}
                onClick={cancelDraw}
              >
                <TrashCan size={14} />
              </button>
            </>
          )}
        </div>
      </div>

      {mode === 'zones' && zones.length > 0 && (
        <div className={styles.regionList}>
          {zones.map((z, zi) => (
            <div key={zi} className={styles.regionListItem}>
              <Input
                size="small"
                className={styles.regionName}
                value={z.name ?? `${formatMessage('Vision.region.zone')} ${zi + 1}`}
                onChange={(e) => patchZone(zi, { name: e.target.value })}
              />
              <span className={styles.regionPointCount}>
                {z.points.length} {formatMessage('Vision.region.points')}
              </span>
              <Switch
                size="small"
                checked={z.enabled !== false}
                onChange={(checked) => patchZone(zi, { enabled: checked })}
              />
              <button
                type="button"
                className={styles.regionRowBtn}
                aria-label={formatMessage('common.delete')}
                onClick={() => removeZone(zi)}
              >
                <TrashCan size={14} />
              </button>
            </div>
          ))}
        </div>
      )}

      {mode === 'zones' && zones.length === 0 && draft.length === 0 && (
        <div className={styles.regionEmptyHint}>{formatMessage('Vision.region.zoneEmpty')}</div>
      )}
      {mode === 'line' && !line && !drawing && (
        <div className={styles.regionEmptyHint}>{formatMessage('Vision.region.lineEmpty')}</div>
      )}
    </div>
  );
};

export default RegionEditor;
