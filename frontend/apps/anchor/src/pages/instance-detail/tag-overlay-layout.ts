// 数据标签 overlay 面板避让布局（移植自源端 viewer-primitives/data-tag-overlay-layout.ts）：
// 锚点左右分道，面板贴画布两侧竖排避让，中央留安全区，引线按角度连回锚点。
export interface DataTagOverlayAnchor {
  id: string;
  screenX: number;
  screenY: number;
  width: number;
  height: number;
}

export interface DataTagOverlayLayoutEntry extends DataTagOverlayAnchor {
  panelOffsetX: number;
  panelOffsetY: number;
  lineLength: number;
  angle: number;
}

interface MutableLayoutEntry extends DataTagOverlayAnchor {
  desiredTop: number;
  placeRight: boolean;
  targetPanelLeft: number;
}

export function layoutDataTagPreviewOverlays(input: {
  overlayWidth: number;
  overlayHeight: number;
  anchors: DataTagOverlayAnchor[];
}): DataTagOverlayLayoutEntry[] {
  const { overlayWidth, overlayHeight, anchors } = input;
  const sideMargin = 12;
  const panelGap = 50;
  const centerSafeWidth = Math.min(120, overlayWidth * 0.18);
  const safeZoneLeft = overlayWidth / 2 - centerSafeWidth / 2;
  const safeZoneRight = overlayWidth / 2 + centerSafeWidth / 2;
  const upperTiltRadians = Math.PI / 6;
  const rightEntries: MutableLayoutEntry[] = [];
  const leftEntries: MutableLayoutEntry[] = [];

  anchors.forEach((anchor) => {
    const placeRight = anchor.screenX > overlayWidth * 0.5;
    const rightLaneMinLeft = Math.min(overlayWidth - anchor.width - sideMargin, Math.max(sideMargin, safeZoneRight));
    const rightLaneMaxLeft = overlayWidth - anchor.width - sideMargin;
    const leftLaneMinLeft = sideMargin;
    const leftLaneMaxLeft = Math.max(sideMargin, Math.min(safeZoneLeft - anchor.width, rightLaneMaxLeft));
    const preferredPanelLeft = placeRight ? anchor.screenX + panelGap : anchor.screenX - anchor.width - panelGap;
    const targetPanelLeft = placeRight
      ? Math.min(rightLaneMaxLeft, Math.max(rightLaneMinLeft, preferredPanelLeft))
      : Math.max(leftLaneMinLeft, Math.min(leftLaneMaxLeft, preferredPanelLeft));
    const anchorX = placeRight ? targetPanelLeft - anchor.screenX : targetPanelLeft - anchor.screenX + anchor.width;
    const anchorInsetY = Math.min(28, anchor.height / 2);
    const tiltOffsetY = Math.max(18, Math.abs(anchorX) * Math.tan(upperTiltRadians));
    const desiredTop =
      anchor.screenY <= overlayHeight * 0.5
        ? anchor.screenY - tiltOffsetY - anchorInsetY
        : anchor.screenY + tiltOffsetY - anchorInsetY;
    const entry: MutableLayoutEntry = {
      ...anchor,
      desiredTop,
      placeRight,
      targetPanelLeft,
    };

    if (placeRight) {
      rightEntries.push(entry);
    } else {
      leftEntries.push(entry);
    }
  });

  const layoutSide = (entries: MutableLayoutEntry[]) => {
    const gap = 12;
    const minTop = 12;
    const maxBottom = overlayHeight - 12;
    const sorted = [...entries].sort((left, right) => left.desiredTop - right.desiredTop);
    let cursorTop = minTop;

    sorted.forEach((entry) => {
      entry.desiredTop = Math.max(entry.desiredTop, cursorTop);
      cursorTop = entry.desiredTop + entry.height + gap;
    });

    for (let index = sorted.length - 1; index >= 0; index -= 1) {
      const entry = sorted[index];
      const bottom = entry.desiredTop + entry.height;
      if (bottom > maxBottom) {
        entry.desiredTop = maxBottom - entry.height;
      }
      if (index < sorted.length - 1) {
        const nextEntry = sorted[index + 1];
        entry.desiredTop = Math.min(entry.desiredTop, nextEntry.desiredTop - entry.height - gap);
      }
      entry.desiredTop = Math.max(minTop, entry.desiredTop);
    }

    return sorted;
  };

  return [...layoutSide(rightEntries), ...layoutSide(leftEntries)].map((entry) => {
    const panelOffsetX = entry.targetPanelLeft - entry.screenX;
    const panelOffsetY = entry.desiredTop - entry.screenY;
    const anchorX = entry.placeRight ? panelOffsetX : panelOffsetX + entry.width;
    const anchorY = panelOffsetY + Math.min(28, entry.height / 2);
    const lineLength = Math.max(12, Math.hypot(anchorX, anchorY));
    const angle = Math.atan2(anchorY, anchorX);

    return {
      id: entry.id,
      screenX: entry.screenX,
      screenY: entry.screenY,
      width: entry.width,
      height: entry.height,
      panelOffsetX,
      panelOffsetY,
      lineLength,
      angle,
    };
  });
}

// 数值展示格式（移植自源端 formatDataTagPreviewValue）
export function formatDataTagPreviewValue(value: unknown) {
  if (value === null || value === undefined || value === '') return '--';
  if (typeof value === 'number' && Number.isFinite(value)) {
    return Number.isInteger(value) ? String(value) : Number(value.toFixed(3)).toString();
  }
  if (typeof value === 'string') {
    const parsed = Number(value);
    if (Number.isFinite(parsed)) {
      return Number.isInteger(parsed) ? String(parsed) : Number(parsed.toFixed(3)).toString();
    }
    return value;
  }
  if (typeof value === 'boolean') {
    return value ? 'true' : 'false';
  }
  try {
    return JSON.stringify(value);
  } catch {
    return String(value);
  }
}
