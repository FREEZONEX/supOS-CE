import type { TableColumnsType } from 'antd';
import type { ProTableColumn, ProTableColumns } from './types';

export const DEFAULT_COL_MIN = 80;
export const DEFAULT_COL_MAX = 420;
export const SELECTION_COL_WIDTH = 35;

type ColumnItem = ProTableColumn;

export const getColumnWidth = (col: ColumnItem) => (typeof col.width === 'number' ? col.width : DEFAULT_COL_MIN);

export const getColumnMin = (col: ColumnItem) => {
  if (typeof col.minWidth === 'number') return col.minWidth;
  if (col.fixed) return getColumnWidth(col);
  return DEFAULT_COL_MIN;
};

export const getColumnMax = (col: ColumnItem) => {
  if (typeof col.maxWidth === 'number') return col.maxWidth;
  if (col.fixed) return Math.max(getColumnWidth(col), 160);
  return DEFAULT_COL_MAX;
};

const isVisibleColumn = (col: ColumnItem) => col.hidden !== true;

const sumWidths = (cols: ProTableColumns, indices: number[]) =>
  indices.reduce((sum, index) => sum + getColumnWidth(cols[index]), 0);

const distributeDelta = (cols: ProTableColumns, indices: number[], delta: number): ProTableColumns => {
  const next = cols.map((col) => ({ ...col }));
  let remaining = Math.abs(delta);
  const grow = delta > 0;
  let pool = [...indices];

  while (remaining > 0.5 && pool.length > 0) {
    const capacities = pool.map((index) => {
      const current = getColumnWidth(next[index]);
      return grow ? getColumnMax(next[index]) - current : current - getColumnMin(next[index]);
    });
    const totalCapacity = capacities.reduce((sum, value) => sum + value, 0);
    if (totalCapacity <= 0) break;

    let moved = 0;
    const nextPool: number[] = [];
    pool.forEach((index, idx) => {
      const capacity = capacities[idx];
      if (capacity <= 0) return;
      const change = Math.min(capacity, remaining * (capacity / totalCapacity));
      const current = getColumnWidth(next[index]);
      const nextWidth = Math.round(current + (grow ? change : -change));
      next[index] = { ...next[index], width: nextWidth };
      moved += change;
      if (capacity - change > 0.5) {
        nextPool.push(index);
      }
    });

    if (moved < 0.5) break;
    remaining -= moved;
    pool = nextPool;
  }

  return next;
};

export const fitColumnsToContainer = (
  cols: TableColumnsType,
  containerWidth: number,
  selectionWidth = 0,
  resize?: { index: number; width: number }
): TableColumnsType => {
  if (!containerWidth) return cols;

  const budget = Math.max(Math.floor(containerWidth - selectionWidth), 0);
  let next: ProTableColumns = cols.map((col) => ({ ...col }));
  const indices = next.map((_, index) => index).filter((index) => isVisibleColumn(next[index]));

  if (resize) {
    const { index, width } = resize;
    const clamped = Math.min(getColumnMax(next[index]), Math.max(getColumnMin(next[index]), Math.round(width)));
    next[index] = { ...next[index], width: clamped };
  }

  let diff = budget - sumWidths(next, indices);
  if (Math.abs(diff) < 1) return next;

  if (diff < 0) {
    let shrinkTargets = indices.filter((index) => !next[index].fixed);
    if (resize) {
      const others = shrinkTargets.filter((index) => index !== resize.index);
      if (others.length > 0) shrinkTargets = others;
    }
    if (shrinkTargets.length > 0) {
      next = distributeDelta(next, shrinkTargets, diff);
      diff = budget - sumWidths(next, indices);
    }
  }

  if (diff > 0) {
    let growTargets = indices.filter((index) => !next[index].fixed);
    if (resize) {
      growTargets = growTargets.filter((index) => index !== resize.index);
    }
    if (growTargets.length === 0) {
      growTargets = indices.filter((index) => index !== resize?.index);
    }
    if (growTargets.length > 0) {
      next = distributeDelta(next, growTargets, diff);
    }
  }

  return next;
};
