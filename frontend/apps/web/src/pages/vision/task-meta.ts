import type { AlgorithmDefinition, AlgorithmType } from '@/apis/core-api/algorithm';
import type { VisionTaskStatus } from '@/apis/core-api/task';

// definition 驱动的算法视图(列表/详情/任务视图均可传入)。
type AlgorithmLike = {
  algoType?: AlgorithmType;
  definition?: AlgorithmDefinition | null;
};

// 各算法类型可选的输出项(fallback,definition.outputs 优先;新分类完全由 definition 驱动)。
export const OUTPUT_OPTIONS: Partial<Record<AlgorithmType, string[]>> = {
  object_detection: ['person_count', 'bounding_boxes'],
  helmet: ['person_count', 'helmet_count', 'no_helmet_count'],
  intrusion: ['intrusion_count', 'intrusion_state', 'intrusion_alert'],
  line_counting: ['total_count', 'in_count', 'out_count'],
};

// 空间配置类型 fallback:line_counting 用计数线,其余用检测区域。
export const spatialMode = (algoType?: AlgorithmType): 'zones' | 'line' =>
  algoType === 'line_counting' ? 'line' : 'zones';

// 输出项:definition.outputs 优先,否则按 algoType 查表,再否则空数组。
export const outputOptionsFor = (algo?: AlgorithmLike | null): string[] => {
  if (algo?.definition?.outputs?.length) return algo.definition.outputs;
  return algo?.algoType ? OUTPUT_OPTIONS[algo.algoType] || [] : [];
};

// 空间配置:definition.spatialMode 优先('none' 表示无区域配置),fallback 按 algoType。
export const spatialModeFor = (algo?: AlgorithmLike | null): 'zones' | 'line' | 'none' =>
  algo?.definition?.spatialMode ?? spatialMode(algo?.algoType);

// 检测频率预设 → 采样间隔(ms)。
export const FREQUENCY_TO_MS: Record<string, number> = { low: 500, standard: 200, high: 100 };

export const DETECTION_FREQUENCY_OPTIONS = [
  { value: 'low', labelKey: 'Vision.task.freqLow' },
  { value: 'standard', labelKey: 'Vision.task.freqStandard' },
  { value: 'high', labelKey: 'Vision.task.freqHigh' },
  { value: 'custom', labelKey: 'Vision.task.freqCustom' },
];

export const RETENTION_OPTIONS = [
  { value: 'permanent', labelKey: 'Vision.task.retentionPermanent' },
  { value: '30', labelKey: 'Vision.task.retention30' },
  { value: '60', labelKey: 'Vision.task.retention60' },
  { value: '180', labelKey: 'Vision.task.retention180' },
  { value: '365', labelKey: 'Vision.task.retention365' },
];

// 每算法结果推送间隔是否展示 fallback(事件类算法展示,object_detection 连续框不展示)。
export const showsPushInterval = (algoType?: AlgorithmType): boolean => algoType !== 'object_detection';

// 推送间隔显隐:definition.showPushInterval 优先,fallback 按 algoType。
export const showsPushIntervalFor = (algo?: AlgorithmLike | null): boolean =>
  algo?.definition?.showPushInterval ?? showsPushInterval(algo?.algoType);

// 相机在线状态徽章。
export const CAMERA_LIVE_META: Record<string, { color: string; labelKey: string }> = {
  online: { color: 'success', labelKey: 'Vision.task.camOnline' },
  offline: { color: 'error', labelKey: 'Vision.task.camOffline' },
  checking: { color: 'processing', labelKey: 'Vision.task.camChecking' },
  unknown: { color: 'default', labelKey: 'Vision.task.camUnknown' },
};

// 算法模型状态徽章。
export const MODEL_STATUS_META: Record<string, { color: string; labelKey: string }> = {
  available: { color: 'success', labelKey: 'Vision.task.modelAvailable' },
  error: { color: 'error', labelKey: 'Vision.task.modelError' },
  missing: { color: 'default', labelKey: 'Vision.task.modelMissing' },
};

// 采样间隔(ms)→ 约等 FPS(列表展示 Recognition Params)。
export const samplingToFps = (ms?: number): number => (ms && ms > 0 ? Math.round(1000 / ms) : 5);

// tone 对应 index.module.scss 里的胶囊配色(设计稿一状态一色)。
export type TaskStatusTone = 'running' | 'starting' | 'warning' | 'error' | 'idle' | 'scheduled';

export const TASK_STATUS_META: Record<VisionTaskStatus, { tone: TaskStatusTone; labelKey: string }> = {
  stopped: { tone: 'idle', labelKey: 'Vision.task.statusStopped' },
  starting: { tone: 'starting', labelKey: 'Vision.task.statusStarting' },
  running: { tone: 'running', labelKey: 'Vision.task.statusRunning' },
  stopping: { tone: 'warning', labelKey: 'Vision.task.statusStopping' },
  error: { tone: 'error', labelKey: 'Vision.task.statusError' },
  offline: { tone: 'idle', labelKey: 'Vision.task.statusOffline' },
  recovering: { tone: 'warning', labelKey: 'Vision.task.statusRecovering' },
  idle: { tone: 'scheduled', labelKey: 'Vision.task.statusIdle' },
};

// 运行态(有实时统计、可停止)。
export const ACTIVE_STATUSES: VisionTaskStatus[] = ['running', 'starting', 'stopping', 'recovering', 'idle'];

// desired=running 语义下用于展示实时统计的状态。
export const RUNNING_STATUSES: VisionTaskStatus[] = ['running', 'recovering'];
