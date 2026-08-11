import type { AlgorithmModelStatus, AlgorithmSource } from '@/apis/core-api/algorithm';

export const SOURCE_META: Record<AlgorithmSource, { color: string; labelKey: string }> = {
  builtin: { color: 'blue', labelKey: 'Vision.algorithm.sourceBuiltin' },
  custom: { color: 'default', labelKey: 'Vision.algorithm.sourceCustom' },
};

export const MODEL_STATUS_META: Record<AlgorithmModelStatus, { color: string; labelKey: string }> = {
  available: { color: 'success', labelKey: 'Vision.algorithm.modelAvailable' },
  missing: { color: 'default', labelKey: 'Vision.algorithm.modelMissing' },
  error: { color: 'error', labelKey: 'Vision.algorithm.modelError' },
};

/** 算法 label 展示:person_down → Person Down(对齐 Figma Recognizable Labels) */
export const formatAlgorithmLabel = (raw: string) =>
  String(raw || '')
    .split(/[_\s-]+/)
    .filter(Boolean)
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join(' ');

const CHIP_GAP = 8;
const CHIP_HORIZONTAL_PADDING = 16;
const CHIP_BORDER = 2;
const CHIP_FONT = '400 12px "IBM Plex Sans", "PingFang SC", sans-serif';

const measureTextWidth = (() => {
  let canvas: HTMLCanvasElement | null = null;
  return (text: string) => {
    if (typeof document === 'undefined') return text.length * 7;
    if (!canvas) canvas = document.createElement('canvas');
    const context = canvas.getContext('2d');
    if (!context) return text.length * 7;
    context.font = CHIP_FONT;
    return context.measureText(text).width;
  };
})();

const chipWidth = (text: string) => Math.ceil(measureTextWidth(text)) + CHIP_HORIZONTAL_PADDING + CHIP_BORDER;

export const fitAlgorithmLabels = (labels: string[], maxWidth: number) => {
  const formatted = (labels || []).map(formatAlgorithmLabel);
  if (formatted.length === 0) return { shown: [] as string[], rest: 0, all: [] as string[] };
  if (formatted.length === 1) return { shown: formatted, rest: 0, all: formatted };

  const shown: string[] = [];
  let used = 0;
  for (let index = 0; index < formatted.length; index += 1) {
    const current = formatted[index];
    const remaining = formatted.length - (index + 1);
    const nextUsed = used + (shown.length > 0 ? CHIP_GAP : 0) + chipWidth(current);
    const moreWidth = remaining > 0 ? CHIP_GAP + chipWidth(`+${remaining}`) : 0;
    if (nextUsed + moreWidth <= maxWidth || shown.length === 0) {
      shown.push(current);
      used = nextUsed;
      continue;
    }
    break;
  }
  return { shown, rest: formatted.length - shown.length, all: formatted };
};
